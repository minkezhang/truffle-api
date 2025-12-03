package query

import (
	"github.com/minkezhang/bene-api/db/enums"
)

type O struct {
	APIs      []enums.ClientAPI
	AtomTypes []enums.AtomType
	Title     string
}

type Q struct {
	apis  map[enums.ClientAPI]bool
	types map[enums.AtomType]bool
	title string
}

func New(o O) *Q {
	q := &Q{
		apis:  map[enums.ClientAPI]bool{},
		types: map[enums.AtomType]bool{},
	}
	for _, api := range o.APIs {
		q.apis[api] = true
	}
	for _, t := range o.AtomTypes {
		q.types[t] = true
	}
	return q
}

func (q *Q) IsSupportedAPI(v enums.ClientAPI) bool     { return q.apis[v] }
func (q *Q) IsSupportedAtomType(v enums.AtomType) bool { return q.types[v] }
func (q *Q) Title() string                             { return q.title }
