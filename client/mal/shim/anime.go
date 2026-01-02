package shim

import (
	"context"
	"strconv"

	"github.com/minkezhang/truffle-api/data/source"
	"github.com/minkezhang/truffle-api/data/source/util"
	"github.com/nstratos/go-myanimelist/mal"

	cpb "github.com/minkezhang/truffle-api/proto/go/config"
)

type AnimeClient struct {
	Config *cpb.MAL
	MAL    mal.Client
}

func (c AnimeClient) Get(ctx context.Context, header source.H) (source.S, error) {
	id, err := strconv.ParseInt(header.ID(), 10, 0)
	if err != nil {
		return source.S{}, err
	}
	result, _, err := c.MAL.Anime.Details(ctx, int(id), fields[ModeAnime])
	if err != nil {
		return source.S{}, err
	}

	return source.Make(Anime{*result}.PB()), nil
}

func (c AnimeClient) Search(ctx context.Context, query string, nsfw bool) ([]source.S, error) {
	f := func(resp *mal.Response) ([]source.S, *mal.Response, error) {
		if resp != nil && resp.NextOffset == 0 {
			return nil, nil, nil
		}
		var offset int
		if resp != nil {
			offset = resp.NextOffset
		}
		results, resp, err := c.MAL.Anime.List(
			ctx,
			query,
			fields[ModeAnime],
			mal.Limit(100),
			mal.Offset(offset),
			mal.NSFW(nsfw),
		)
		return util.Apply(results, func(v mal.Anime) source.S {
			return source.Make(Anime{v}.PB())
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
