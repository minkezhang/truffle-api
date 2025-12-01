// Package tv defines an atom of the TV type.
package tv

import (
	"github.com/minkezhang/bene-api/db/atom"
)

var (
	_ atom.A[*T] = &T{}
)

type O struct {
	Base   atom.O
	Header atom.H

	Season           int
	IsAnimated       bool
	Genres           []string
	Showrunners      []string
	Directors        []string
	Writers          []string
	Cinematography   []string
	Composers        []string
	Starring         []string
	AnimationStudios []string
}

func New(o O) (*T, error) {
	return &T{
		Base:             atom.New(o.Base).WithHeader(o.Header),
		Season:           o.Season,
		IsAnimated:       o.IsAnimated,
		Genres:           o.Genres,
		Showrunners:      o.Showrunners,
		Directors:        o.Directors,
		Writers:          o.Writers,
		Cinematography:   o.Cinematography,
		Composers:        o.Composers,
		Starring:         o.Starring,
		AnimationStudios: o.AnimationStudios,
	}, nil
}

type T struct {
	*atom.Base

	Season           int
	IsAnimated       bool
	Genres           []string
	Showrunners      []string
	Directors        []string
	Writers          []string
	Cinematography   []string
	Composers        []string
	Starring         []string
	AnimationStudios []string
}

func (t *T) Type() atom.AtomType { return atom.AtomTypeTV }
func (t *T) Merge(other *T) *T {
	return nil // unimplmented
}
