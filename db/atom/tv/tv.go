// Package tv defines an atom of the TV type.
package tv

import (
	"github.com/minkezhang/bene-api/db/atom"
)

var (
	_ atom.A[*T] = &T{}
)

type O struct {
	atom.O

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

func DebugNewOrDie(o O) *T {
	a, err := New(o)
	if err != nil {
		panic(err)
	}
	return a
}

func New(o O) (*T, error) {
	a, err := atom.New(o.O)
	if err != nil {
		return nil, err
	}
	return &T{
		Base:             a,
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
	if t == nil {
		t = &T{}
	}
	return &T{
		Base:       t.Base.Merge(other.Base),
		Season:     other.Season,
		IsAnimated: other.IsAnimated,
		Genres: append(
			append([]string{}, t.Genres...),
			other.Genres...),
		Showrunners: append(
			append([]string{}, t.Showrunners...),
			other.Showrunners...),
		Directors: append(
			append([]string{}, t.Directors...),
			other.Directors...),
		Writers: append(
			append([]string{}, t.Writers...),
			other.Writers...),
		Cinematography: append(
			append([]string{}, t.Cinematography...),
			other.Cinematography...),
		Composers: append(
			append([]string{}, t.Composers...),
			other.Composers...),
		Starring: append(
			append([]string{}, t.Starring...),
			other.Starring...),
		AnimationStudios: append(
			append([]string{}, t.AnimationStudios...),
			other.AnimationStudios...),
	}
}
