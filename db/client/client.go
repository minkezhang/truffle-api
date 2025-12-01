// Package client is a collection of API hooks to distinct metadata aggregation
// sources.
//
// Bene will query these sources for data.
package client

import (
	"context"
	"regexp"

	"github.com/minkezhang/bene-api/db/atom"
)

type Q[T atom.A[T]] struct {
	Title string
	ID    string
}

type C[T atom.A[T]] interface {
	API() atom.ClientAPI
	IsSupported(t atom.AtomType) bool
	Query(ctx context.Context, q Q[T]) ([]T, error)
}

type Base[T atom.A[T]] struct {
	api       atom.ClientAPI
	supported map[atom.AtomType]interface{}
}

type O[T atom.A[T]] struct {
	API            atom.ClientAPI
	SupportedTypes []atom.AtomType
}

func New[T atom.A[T]](o O[T]) (*Base[T], error) {
	c := &Base[T]{
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

func (c *Base[T]) Query(ctx context.Context, q Q[T]) ([]T, error) { return nil, nil }

type Bene[T atom.A[T]] struct {
	*Base[T]
	cache map[string]T
}

func (c *Bene[T]) Query(ctx context.Context, q Q[T]) ([]T, error) {
	if q.ID != "" {
		if v, ok := c.cache[q.ID]; ok {
			return []T{v}, nil
		}
		return nil, nil
	}
	pattern, err := regexp.Compile(q.ID)
	if err != nil {
		return nil, err
	}
	res := []T{}
	for _, v := range c.cache {
		if pattern.MatchString(v.ID()) {
			res = append(res, v)
		}
		/*for _, t := range v.Titles {
			if pattern.MatchString(t.Title) {
				res = append(res, v)
			}
		}
		*/
	}
	return res, nil
}
