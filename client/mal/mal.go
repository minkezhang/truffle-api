// Package mal defines the shim for MyAnimeList APIv2.
//
// See https://myanimelist.net/apiconfig/references/api/v2 for more information
// on the API. See
// https://help.myanimelist.net/hc/en-us/articles/900003108823-API for more
// information on how to get an API key.
package mal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/minkezhang/truffle-api/client/mal/shim"
	"github.com/minkezhang/truffle-api/client/option"
	"github.com/minkezhang/truffle-api/data/source"
	"github.com/nstratos/go-myanimelist/mal"

	cpb "github.com/minkezhang/truffle-api/proto/go/config"
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

func Make(pb *cpb.MAL) C {
	client := *mal.NewClient(
		&http.Client{
			Transport: transport{id: pb.GetClientId()},
		},
	)

	return C{
		anime: shim.AnimeClient{
			Config: pb,
			MAL:    client,
		},
		manga: shim.MangaClient{
			Config: pb,
			MAL:    client,
		},
	}
}

type C struct {
	anime shim.AnimeClient
	manga shim.MangaClient
}

func (c C) Get(ctx context.Context, header source.H) (source.S, error) {
	switch api := header.API(); api {
	case epb.SourceAPI_SOURCE_API_MAL_MANGA_PARTIAL:
		return c.manga.Get(ctx, header)
	case epb.SourceAPI_SOURCE_API_MAL_ANIME_PARTIAL:
		return c.anime.Get(ctx, header)
	case epb.SourceAPI_SOURCE_API_MAL:
		switch t := header.Type(); t {
		case epb.SourceType_SOURCE_TYPE_SERIES_ANIME:
			fallthrough
		case epb.SourceType_SOURCE_TYPE_MOVIE_ANIME:
			return c.anime.Get(ctx, header)
		case epb.SourceType_SOURCE_TYPE_BOOK_MANGA:
			fallthrough
		case epb.SourceType_SOURCE_TYPE_BOOK_LIGHT_NOVEL:
			return c.manga.Get(ctx, header)
		default:
			return source.S{}, fmt.Errorf("unsupported source type: %v", t.String())
		}
	default:
		return source.S{}, fmt.Errorf("unsupported API: %v", api.String())
	}
}

func (c C) Search(ctx context.Context, query string, opts ...option.O) ([]source.S, error) {
	var results []source.S
	var sources []source.S
	var err error

	nsfw := false
	for _, o := range opts {
		switch o := o.(type) {
		case option.NSFW:
			nsfw = bool(o)
		}
	}

	sources, err = c.anime.Search(ctx, query, nsfw)
	if err != nil {
		return nil, err
	}
	results = append(results, sources...)

	sources, err = c.manga.Search(ctx, query, nsfw)
	if err != nil {
		return nil, err
	}
	results = append(results, sources...)

	return results, nil
}

type transport struct {
	id string
}

func (t transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-MAL-CLIENT-ID", t.id)
	return http.DefaultTransport.RoundTrip(req)
}
