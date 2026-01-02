package cache

import (
	"context"

	"github.com/minkezhang/truffle-api/client/option"
	"github.com/minkezhang/truffle-api/data/node"
	"github.com/minkezhang/truffle-api/data/source"
	"github.com/minkezhang/truffle-api/data/source/util"
	"github.com/minkezhang/truffle-api/util/match"

	dpb "github.com/minkezhang/truffle-api/proto/go/data"
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

func New(ctx context.Context, data []*dpb.Source) *C {
	cache := &C{
		cache: map[epb.SourceAPI]map[epb.SourceType]map[string]source.S{},
	}
	for _, s := range data {
		cache.Put(ctx, source.Make(s))
	}
	return cache
}

type C struct {
	cache map[epb.SourceAPI]map[epb.SourceType]map[string]source.S
}

func (c *C) Delete(ctx context.Context, header source.H) error {
	if _, ok := c.cache[header.API()]; !ok {
		return nil
	}
	if _, ok := c.cache[header.API()][header.Type()]; !ok {
		return nil
	}
	delete(c.cache[header.API()][header.Type()], header.ID())
	return nil
}

func (c *C) Get(ctx context.Context, header source.H) (source.S, error) {
	if _, ok := c.cache[header.API()]; !ok {
		return source.S{}, nil
	}
	if _, ok := c.cache[header.API()][header.Type()]; !ok {
		return source.S{}, nil
	}
	return c.cache[header.API()][header.Type()][header.ID()], nil
}

func (c *C) Put(ctx context.Context, s source.S) (source.H, error) {
	if _, ok := c.cache[s.Header().API()]; !ok {
		c.cache[s.Header().API()] = map[epb.SourceType]map[string]source.S{}
	}
	if _, ok := c.cache[s.Header().API()][s.Header().Type()]; !ok {
		c.cache[s.Header().API()][s.Header().Type()] = map[string]source.S{}
	}
	c.cache[s.Header().API()][s.Header().Type()][s.Header().ID()] = s
	return s.Header(), nil
}

func (c *C) Search(ctx context.Context, query string, opts ...option.O) ([]source.S, error) {
	var results []source.S
	for api := range c.cache {
		for t := range c.cache[api] {
			for _, v := range c.cache[api][t] {
				if h, _ := match.RegExp(query, v); h > 0 {
					results = append(results, v)
				}
			}
		}
	}
	return results, nil
}

func (c *C) PB() []*dpb.Source {
	sources := []source.S{}
	for api := range c.cache {
		for t := range c.cache[api] {
			for _, v := range c.cache[api][t] {
				sources = append(sources, v)
			}
		}
	}
	return util.Apply(sources, func(v source.S) *dpb.Source { return v.PB() })
}

func (c *C) SearchByNodeID(ctx context.Context, header node.H) ([]source.S, error) {
	results := []source.S{}
	for api := range c.cache {
		for t := range c.cache[api] {
			for _, v := range c.cache[api][t] {
				if v.NodeID() == header.ID() {
					v, err := c.Get(ctx, v.Header())
					if err != nil {
						return nil, err
					}
					results = append(results, v)
				}
			}
		}
	}
	return results, nil
}
