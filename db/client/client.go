package client

import (
	"context"

	"github.com/minkezhang/bene-api/db/atom"
)

type Q[T atom.A[T]] struct {
	Body T
}

type C[T atom.A[T]] interface {
	API() atom.ClientAPI
	IsSupported(t atom.AtomType) bool
	Query(ctx context.Context, q Q[T]) ([]T, error)
}

type Base[T atom.A[T]] struct {
	cache     map[string]T
	api       atom.ClientAPI
	supported map[atom.AtomType]interface{}
}

type O[T atom.A[T]] struct {
	Cache          map[string]T
	API            atom.ClientAPI
	SupportedTypes []atom.AtomType
}

func New[T atom.A[T]](o O[T]) (*Base[T], error) {
	c := &Base[T]{
		cache:     o.Cache,
		api:       o.API,
		supported: map[atom.AtomType]interface{}{},
	}
	for _, t := range o.SupportedTypes {
		c.supported[t] = struct{}{}
	}
	return c, nil
}

func (c *Base[T]) API() atom.ClientAPI { return c.api }

func (c *Base[T]) IsSupported(t atom.AtomType) bool {
	_, ok := c.supported[t]
	return ok
}

func (c *Base[T]) Query(ctx context.Context, q Q[T]) ([]T, error) {
	if c.IsSupported(q.Body.Type()) && q.Body.API() == c.API() && q.Body.ID() != "" {
		if v, ok := c.cache[q.Body.ID()]; ok {
			return []T{v}, nil
		}
	}
	return nil, nil
}

type Bene[T atom.A[T]] struct {
	*Base[T]
}

func (c *Bene[T]) Query(ctx context.Context, q Q[T]) ([]T, error) {
	return c.Base.Query(ctx, q)
}
