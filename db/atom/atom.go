// Package atom is a collection of discrete data sources used to represent
// different media types.
//
// Each atom is backed by a single source of truth, e.g. MAL, Spotify, etc.
//
// Different media types have media-specific types of data -- a song for example
// will need the concept of a composer, which is not the case for a book. This
// media-specific data is wrapped in the `A.aux` field.
package atom

import (
	"fmt"

	"github.com/minkezhang/bene-api/db/enums"
)

type T struct {
	Title        string
	Localization string
}

type Aux interface {
	AtomType() enums.AtomType
	Copy() Aux
	Merge(o Aux) Aux
	Equal(o Aux) bool
}

type O struct {
	APIType    enums.ClientAPI
	APIID      string
	Titles     []T
	PreviewURL string
	Score      int
	AtomType   enums.AtomType
	Aux        Aux
}

type A struct {
	// Shared properties
	apiType    enums.ClientAPI                   // Read-only
	apiID      string                            // Read-only
	titles     map[string]map[string]interface{} // e.g. a.titles["us-en"]["Firefly"]
	previewURL string
	score      int

	atomType enums.AtomType // Read-only
	aux      Aux            // Media-specific data
}

func New(o O) *A {
	a := &A{
		apiType:    o.APIType,
		apiID:      o.APIID,
		previewURL: o.PreviewURL,
		score:      o.Score,
		atomType:   o.AtomType,
		aux:        o.Aux,
	}
	a.SetTitles(o.Titles)
	return a
}

func (a *A) APIType() enums.ClientAPI { return a.apiType }
func (a *A) APIID() string            { return a.apiID }
func (a *A) PreviewURL() string       { return a.previewURL }
func (a *A) Score() int               { return a.score }
func (a *A) AtomType() enums.AtomType { return a.atomType }
func (a *A) Aux() Aux                 { return a.aux.Copy() }
func (a *A) Copy() *A {
	return New(O{
		APIType:    a.apiType,
		APIID:      a.apiID,
		PreviewURL: a.previewURL,
		Score:      a.score,
		AtomType:   a.atomType,
		Aux:        a.aux.Copy(),
		Titles:     a.Titles(),
	})
}

func (a *A) Titles() []T {
	res := []T{}
	for l, ts := range a.titles {
		for t, _ := range ts {
			res = append(res, T{Title: t, Localization: l})
		}
	}
	return res
}

func (a *A) SetTitles(v []T) {
	a.titles = map[string]map[string]interface{}{}
	for _, t := range v {
		if _, ok := a.titles[t.Localization]; !ok {
			a.titles[t.Localization] = map[string]interface{}{}
		}
		a.titles[t.Localization][t.Title] = struct{}{}
	}
}

func (a *A) SetPreviewURL(v string) { a.previewURL = v }
func (a *A) SetScore(v int)         { a.score = v }
func (a *A) SetAux(v Aux) {
	if a.atomType != v.AtomType() {
		panic(fmt.Errorf("cannot set mismatching atom types: %v != %v", a.atomType, v.AtomType()))
	}
	a.aux = v.Copy()
}

// Merge will combine two atoms, with the following heuristic --
//
//  1. the primitives of two merged atoms will be overwritten by the second
//     atom
//  2. slices and maps of the atoms are a union of the two inputs
//  3. structs (i.e. a.Aux()) will be recursively merged with the same
//     heuristic
func (a *A) Merge(o *A) *A {
	if a.atomType != o.atomType {
		panic(fmt.Errorf("cannot merge mismatching atom types: %v != %v", a.atomType, o.atomType))
	}
	return New(O{
		APIType: o.apiType,
		APIID:   o.apiID,
		Titles: append(
			append([]T{}, a.Titles()...),
			o.Titles()...,
		),
		PreviewURL: o.previewURL,
		Score:      o.score,
		AtomType:   o.atomType,
		Aux:        a.aux.Merge(o.aux),
	})
}
