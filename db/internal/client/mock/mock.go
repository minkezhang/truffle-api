package mock

import (
	"context"

	"github.com/minkezhang/bene-api/client/query"
	"github.com/minkezhang/bene-api/db/atom"

	epb "github.com/minkezhang/bene-api/proto/go/enums"
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

func (c *C) APIType() epb.API                                    { return epb.API_API_MAL }
func (c *C) Get(ctx context.Context, g query.G) (*atom.A, error) { return c.data[g.ID], nil }

func (c *C) Query(ctx context.Context, q *query.Q) ([]*atom.A, error) {
	res := []*atom.A{}
	for _, a := range c.data {
		match, err := q.Match(a)
		if err != nil {
			return nil, err
		}
		if match {
			res = append(res, a)
		}
	}
	return res, nil
}
