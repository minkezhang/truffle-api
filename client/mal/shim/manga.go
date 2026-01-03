package shim

import (
	"context"
	"strconv"

	"github.com/minkezhang/truffle-api/data/source"
	"github.com/minkezhang/truffle-api/util/slice"
	"github.com/nstratos/go-myanimelist/mal"

	cpb "github.com/minkezhang/truffle-api/proto/go/config"
)

type MangaClient struct {
	Config *cpb.MAL
	MAL    mal.Client
}

func (c MangaClient) Get(ctx context.Context, header source.H) (source.S, error) {
	id, err := strconv.ParseInt(header.ID(), 10, 0)
	if err != nil {
		return source.S{}, err
	}
	result, _, err := c.MAL.Manga.Details(ctx, int(id), fields[ModeManga])
	if err != nil {
		return source.S{}, err
	}

	return source.Make(Manga{*result}.PB()), nil
}

func (c MangaClient) Search(ctx context.Context, query string, nsfw bool) ([]source.S, error) {
	f := func(resp *mal.Response) ([]source.S, *mal.Response, error) {
		if resp != nil && resp.NextOffset == 0 {
			return nil, nil, nil
		}
		var offset int
		if resp != nil {
			offset = resp.NextOffset
		}
		results, resp, err := c.MAL.Manga.List(
			ctx,
			query,
			fields[ModeManga],
			mal.Limit(100),
			mal.Offset(offset),
			mal.NSFW(nsfw),
		)
		return slice.Apply(results, func(v mal.Manga) source.S {
			return source.Make(Manga{v}.PB())
		}), resp, err
	}

	var results []source.S
	var page []source.S
	var resp *mal.Response
	var err error

	// Aggregate all results
	for page, resp, err = f(nil); err == nil && page != nil; page, resp, err = f(resp) {
		results = append(results, page...)
		if len(results) >= int(c.Config.GetMaxResults()) {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	end := len(results)
	if int(c.Config.GetMaxResults()) < len(results) {
		end = int(c.Config.GetMaxResults())
	}
	return results[:end], nil
}
