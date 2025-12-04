package query

import (
	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

type O struct {
	APIs      []epb.API
	AtomTypes []epb.Type
	Title     string
}

type Q struct {
	apis  map[epb.API]bool
	types map[epb.Type]bool
	title string
}

func New(o O) *Q {
	q := &Q{
		apis:  map[epb.API]bool{},
		types: map[epb.Type]bool{},
	}
	for _, api := range o.APIs {
		q.apis[api] = true
	}
	for _, t := range o.AtomTypes {
		q.types[t] = true
	}
	return q
}

func (q *Q) IsSupportedAPI(v epb.API) bool   { return q.apis[v] }
func (q *Q) IsSupportedType(v epb.Type) bool { return q.types[v] }
func (q *Q) Title() string                   { return q.title }
