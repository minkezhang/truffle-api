package client

import (
	"context"

	"github.com/minkezhang/truffle-api/client/cache"
	"github.com/minkezhang/truffle-api/client/option"
	"github.com/minkezhang/truffle-api/data/node"
	"github.com/minkezhang/truffle-api/data/source"

	dpb "github.com/minkezhang/truffle-api/proto/go/data"
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

type RW interface {
	Put(ctx context.Context, s source.S) (source.H, error)
	Delete(ctx context.Context, header source.H) error
}

type C interface {
	Get(ctx context.Context, header source.H) (source.S, error)
	Search(ctx context.Context, query string, opts ...option.O) ([]source.S, error)
}

func New(ctx context.Context, c C, data []*dpb.Source) *Cache {
	return &Cache{
		cache: cache.New(ctx, data),
		c:     c,
	}
}

type Cache struct {
	cache *cache.C
	c     C
}

func (c *Cache) Put(ctx context.Context, s source.S) (source.H, error) {
	switch client := c.c.(type) {
	case RW:
		var header source.H
		var err error
		header, err = client.Put(ctx, s)
		if err != nil {
			return source.H{}, err
		}
		s = s.WithHeader(header)
	}

	return c.cache.Put(ctx, s)
}

func (c *Cache) Delete(ctx context.Context, header source.H) error {
	switch client := c.c.(type) {
	case RW:
		if err := client.Delete(ctx, header); err != nil {
			return err
		}
	}
	return c.cache.Delete(ctx, header)
}

func (c *Cache) Get(ctx context.Context, header source.H, remote option.Remote) (source.S, error) {
	remote = option.Remote(
		bool(remote) || header.API() == epb.SourceAPI_SOURCE_API_TRUFFLE,
	)

	var s source.S
	if bool(remote) {
		var err error
		s, err = c.c.Get(ctx, header)
		if err != nil {
			return source.S{}, nil
		}
		if s != (source.S{}) {
			t, _ := c.cache.Get(ctx, header)
			s = s.WithNodeID(t.NodeID())
			if _, err := c.cache.Put(ctx, s); err != nil {
				return source.S{}, err
			}
		}
		return s, err
	}
	return c.cache.Get(ctx, header)
}

func (c *Cache) Search(ctx context.Context, query string, opts ...option.O) ([]source.S, error) {
	remote := false
	for _, o := range opts {
		switch o := o.(type) {
		case option.Remote:
			remote = bool(o)
		}
	}

	if remote {
		results, err := c.c.Search(ctx, query, opts...)
		if err != nil {
			return nil, err
		}
		for i, s := range results {
			t, _ := c.cache.Get(ctx, s.Header())
			s = s.WithNodeID(t.NodeID())
			if _, err := c.cache.Put(ctx, s); err != nil {
				return nil, err
			}
			results[i] = s
		}
		return results, nil
	}
	return c.cache.Search(ctx, query, opts...)
}

func (c *Cache) PB() []*dpb.Source { return c.cache.PB() }

func (c *Cache) SearchByNodeID(ctx context.Context, header node.H, remote option.Remote) ([]source.S, error) {
	results, err := c.cache.SearchByNodeID(ctx, header)
	if err != nil {
		return nil, err
	}
	if bool(remote) {
		for i, s := range results {
			t, _ := c.c.Get(ctx, s.Header())
			s = t.WithNodeID(s.NodeID())
			if _, err := c.cache.Put(ctx, s); err != nil {
				return nil, err
			}
			results[i] = s
		}
	}
	return results, nil
}
