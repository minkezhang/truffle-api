package query

import (
	"github.com/minkezhang/bene-api/client/query"

	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

type O struct {
	APIs      []epb.API
	AtomTypes []epb.Type
	Title     string
}

type Q struct {
	*query.Q
	apis map[epb.API]bool
}

func New(o O) *Q {
	q := &Q{
		Q: query.New(query.O{
			AtomTypes: append([]epb.Type{}, o.AtomTypes...),
			Title:     o.Title,
		}),
		apis: map[epb.API]bool{},
	}
	for _, api := range o.APIs {
		q.apis[api] = true
	}
	return q
}

func (q *Q) IsSupportedAPI(v epb.API) bool { return q.apis[v] }
