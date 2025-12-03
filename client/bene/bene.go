// Package bene provides a client interface for bene.
//
// This client caches all atoms (i.e. atoms from all client APIs).
package bene

import (
	"context"
	"regexp"

	"github.com/minkezhang/bene-api/client"
	"github.com/minkezhang/bene-api/db/atom"
	"github.com/minkezhang/bene-api/db/enums"
)

var (
	_ client.C = &C{}
)

type O struct {
	Cache []*atom.A
}

type C struct {
	cache map[string]*atom.A
}

func New(o O) (*C, error) {
	c := &C{
		cache: map[string]*atom.A{},
	}
	for _, a := range o.Cache {
		c.Add(a)
	}
	return c, nil
}

func (c *C) APIType() enums.ClientAPI                            { return enums.ClientAPIBene }
func (c *C) Add(a *atom.A)                                       { c.cache[a.APIID()] = a.Copy() }
func (c *C) Get(ctx context.Context, id string) (*atom.A, error) { return c.cache[id], nil }
func (c *C) Remove(id string)                                    { delete(c.cache, id) }

func (c *C) Query(ctx context.Context, q client.Q) ([]*atom.A, error) {
	pattern, err := regexp.Compile(q.Title)
	if err != nil {
		return nil, err
	}

	res := []*atom.A{}
	for _, a := range c.cache {
		for _, t := range a.Titles() {
			if pattern.MatchString(t.Title) {
				res = append(res, a.Copy())
			}
		}
	}
	return res, nil
}
