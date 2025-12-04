package mock

import (
	"context"

	"github.com/minkezhang/bene-api/client"
	"github.com/minkezhang/bene-api/client/query"
	"github.com/minkezhang/bene-api/db/atom"
	"github.com/minkezhang/bene-api/db/enums"
)

type O struct {
	Data []*atom.A
}

func New(o O) *C {
	c := &C{
		data: map[string]*atom.A{},
	}
	for _, a := range o.Data {
		c.data[a.APIID()] = a
	}
	return c
}

type C struct {
	data map[string]*atom.A
}

func (c *C) APIType() enums.ClientAPI                            { return enums.ClientAPIMAL }
func (c *C) Get(ctx context.Context, id string) (*atom.A, error) { return c.data[id], nil }

func (c *C) Query(ctx context.Context, q query.Q) ([]*atom.A, error) {
	res := []*atom.A{}
	for _, a := range c.data {
		match, err := client.Match(q, a)
		if err != nil {
			return nil, err
		}
		if match {
			res = append(res, a)
		}
	}
	return res, nil
}
