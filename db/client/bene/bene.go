package bene

import (
	"context"
	"regexp"

	"github.com/minkezhang/bene-api/db/atom"
	"github.com/minkezhang/bene-api/db/client"
)

type O[T atom.A[T]] struct {
	Cache []T
}

func DebugNewOrDie[T atom.A[T]](o O[T]) *Bene[T] {
	c, err := New(o)
	if err != nil {
		panic(err)
	}
	return c
}

func New[T atom.A[T]](o O[T]) (*Bene[T], error) {
	a, err := client.New(client.O{
		API: atom.ClientAPIBene,
		SupportedTypes: []atom.AtomType{
			atom.AtomTypeTV,
		},
	})
	if err != nil {
		return nil, err
	}
	c := &Bene[T]{
		Base:  a,
		cache: map[string]T{},
	}
	for _, a := range o.Cache {
		c.cache[a.ID()] = a
	}
	return c, nil
}

type Bene[T atom.A[T]] struct {
	*client.Base
	cache map[string]T
}

func (c *Bene[T]) Get(ctx context.Context, id string) (T, error) {
	if v, ok := c.cache[id]; ok {
		return v, nil
	}
	var a T
	return a, nil
}

func (c *Bene[T]) Query(ctx context.Context, q client.Q[T]) ([]T, error) {
	pattern, err := regexp.Compile(q.Title)
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

func (c *Bene[T]) Create(ctx context.Context, a atom.A[T]) error {
	c.cache[a.ID()] = a
	return nil
}

func (c *Bene[T]) Delete(ctx context.Context, id string) error {
	delete(c.cache, id)
	return nil
}
