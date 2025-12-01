package bene

import (
	"context"
	"regexp"

	"github.com/minkezhang/bene-api/db/atom"
	"github.com/minkezhang/bene-api/db/client"
)

type O[T atom.A[T]] struct {
	Cache map[string]T
}

func New[T atom.A[T]](o O[T]) *Bene[T] {
	return &Bene[T]{
		Base: client.New(client.O{
			API: atom.ClientAPIBene,
			SupportedTypes: []atom.AtomType{
				atom.AtomTypeTV,
			},
		}),
		cache: o.Cache,
	}
}

type Bene[T atom.A[T]] struct {
	*client.Base
	cache map[string]T
}

func (c *Bene[T]) Query(ctx context.Context, q client.Q[T]) ([]T, error) {
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
		for _, t := range v.GetBase().Titles {
			if pattern.MatchString(t.Title) {
				res = append(res, v)
			}
		}
	}
	return res, nil
}
