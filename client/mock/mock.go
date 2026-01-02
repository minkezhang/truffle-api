package mock

import (
	"context"

	"github.com/minkezhang/truffle-api/client/option"
	"github.com/minkezhang/truffle-api/data/source"
)

func New(sources []source.S) *C {
	c := &C{
		sources: map[string]source.S{},
	}

	for _, s := range sources {
		c.sources[s.Header().ID()] = s
	}

	return c
}

type C struct {
	sources map[string]source.S

	GetHistory    []source.H
	SearchHistory []string
}

func (c *C) Get(ctx context.Context, header source.H) (source.S, error) {
	c.GetHistory = append(c.GetHistory, header)
	return c.sources[header.ID()], nil
}

func (c *C) Search(ctx context.Context, query string, opts ...option.O) ([]source.S, error) {
	c.SearchHistory = append(c.SearchHistory, query)
	results := []source.S{}
	for _, s := range c.sources {
		for _, t := range s.Titles() {
			if t.Title() == query {
				results = append(results, s)
			}
		}
	}
	return results, nil
}
