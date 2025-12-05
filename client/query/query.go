package query

import (
	"regexp"

	"github.com/minkezhang/bene-api/db/atom"

	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

type G struct {
	AtomType epb.Type
	ID       string
}

type O struct {
	AtomTypes []epb.Type
	Title     string
}

type Q struct {
	types map[epb.Type]bool
	title string
}

func New(o O) *Q {
	q := &Q{
		types: map[epb.Type]bool{},
		title: o.Title,
	}
	for _, t := range o.AtomTypes {
		q.types[t] = true
	}
	return q
}

func (q *Q) IsSupportedType(v epb.Type) bool { return q.types[v] }
func (q *Q) Title() string                   { return q.title }

func (q *Q) Match(a *atom.A) (bool, error) {
	pattern, err := regexp.Compile(q.Title())
	if err != nil {
		return false, err
	}

	if !q.IsSupportedType(a.AtomType()) {
		return false, nil
	}

	for _, t := range a.Titles() {
		if pattern.MatchString(t.Title) {
			return true, nil
		}
	}
	return false, nil

}
