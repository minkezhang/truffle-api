package mal

import (
	"context"
	"net/http"

	"github.com/minkezhang/bene-api/client"
	"github.com/minkezhang/bene-api/client/mal/anime"
	"github.com/minkezhang/bene-api/client/mal/manga"
	"github.com/minkezhang/bene-api/client/query"
	"github.com/minkezhang/bene-api/db/atom"
	"github.com/nstratos/go-myanimelist/mal"

	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

var (
	_ client.C = &C{}
)

type O struct {
	ClientID         string
	PopularityCutoff int64
	MaxResults       int64
	NSFW             bool
}

func New(o O) *C {
	return &C{
		anime: &anime.C{
			Cutoff:  int(o.PopularityCutoff),
			Results: int(o.MaxResults),
			NSFW:    o.NSFW,
			MAL: *mal.NewClient(
				&http.Client{
					Transport: transport{id: o.ClientID},
				},
			),
		},
		manga: &manga.C{
			Cutoff:  int(o.PopularityCutoff),
			Results: int(o.MaxResults),
			NSFW:    o.NSFW,
			MAL: *mal.NewClient(
				&http.Client{
					Transport: transport{id: o.ClientID},
				},
			),
		},
	}
}

type C struct {
	anime *anime.C
	manga *manga.C
}

func (c *C) APIType() epb.API { return epb.API_API_MAL }

func (c *C) Get(ctx context.Context, g query.G) (*atom.A, error) {
	switch t := g.AtomType; t {
	case epb.Type_TYPE_TV:
		fallthrough
	case epb.Type_TYPE_MOVIE:
		return c.anime.Get(ctx, g)
	case epb.Type_TYPE_BOOK:
		return c.manga.Get(ctx, g)
	default:
		return nil, nil /* unimplemented */
	}
}

func (c *C) Query(ctx context.Context, q *query.Q) ([]*atom.A, error) {
	var results []*atom.A
	if q.IsSupportedType(epb.Type_TYPE_TV) || q.IsSupportedType(epb.Type_TYPE_MOVIE) {
		atoms, err := c.anime.Query(ctx, q)
		if err != nil {
			return nil, err
		}
		results = append(results, atoms...)
	}
	if q.IsSupportedType(epb.Type_TYPE_BOOK) {
		/* unimplemented */
	}
	return results, nil
}

type transport struct {
	id string
}

func (t transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-MAL-CLIENT-ID", t.id)
	return http.DefaultTransport.RoundTrip(req)
}
