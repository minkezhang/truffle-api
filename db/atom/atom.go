// Package atom is a collection of discrete data sources used to represent
// different media types.
//
// Each atom is backed by a single source of truth, e.g. MAL, Spotify, etc.
//
// Different media types have media-specific types of data -- a song for example
// will need the concept of a composer, which is not the case for a book. This
// media-specific data is wrapped in the `A.metadata` field.
package atom

import (
	"fmt"

	"github.com/minkezhang/bene-api/db/atom/metadata"
	"github.com/minkezhang/bene-api/db/atom/metadata/empty"
	"google.golang.org/protobuf/proto"

	apb "github.com/minkezhang/bene-api/proto/go/atom"
	mpb "github.com/minkezhang/bene-api/proto/go/atom/metadata"
	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

type G struct{}

func (g G) Load(msg proto.Message) *A {
	pb := msg.(*apb.Atom)
	a := New(O{
		AtomType:   pb.GetType(),
		APIType:    pb.GetApi(),
		APIID:      pb.GetId(),
		PreviewURL: pb.GetPreviewUrl(),
		Score:      pb.GetScore(),
	})

	titles := []T{}
	for _, t := range pb.GetTitles() {
		titles = append(titles, T{
			Title:        t.GetTitle(),
			Localization: t.GetLocalization(),
		})
	}
	a.SetTitles(titles)

	switch t := pb.GetType(); t {
	case epb.Type_TYPE_TV:
		// a.SetMetadata(tv.G{}.Load(pb.GetMetadataTv()))
	default:
		a.SetMetadata(empty.G{}.Load(pb.GetMetadataEmpty()))
	}

	return a
}

func (g G) Save(a *A) proto.Message {
	pb := &apb.Atom{
		Type:       a.AtomType(),
		Api:        a.APIType(),
		Id:         a.APIID(),
		PreviewUrl: a.PreviewURL(),
		Score:      int64(a.Score()),
	}

	for _, t := range a.Titles() {
		pb.Titles = append(pb.GetTitles(), &apb.Title{
			Title:        t.Title,
			Localization: t.Localization,
		})
	}

	switch t := a.AtomType(); t {
	case epb.Type_TYPE_TV:
		// pb.Metadata = &apb.Atom_MetadataTv{tv.G{}.Save(a.Metadata().(tv.M)).(*mpb.TV)}
	default:
		pb.Metadata = &apb.Atom_MetadataEmpty{empty.G{}.Save(a.Metadata().(empty.M)).(*mpb.Empty)}
	}

	return pb
}

type T struct {
	Title        string
	Localization string
}

type O struct {
	APIType    epb.API
	APIID      string
	Titles     []T
	PreviewURL string
	Score      int64
	AtomType   epb.Type
	Metadata   metadata.M
}

type A struct {
	apiType    epb.API                           // Read-only
	apiID      string                            // Read-only
	titles     map[string]map[string]interface{} // e.g. a.titles["us-en"]["Firefly"]
	previewURL string
	score      int64

	atomType epb.Type   // Read-only
	metadata metadata.M // Media-specific data
}

func New(o O) *A {
	a := &A{
		apiType:    o.APIType,
		apiID:      o.APIID,
		previewURL: o.PreviewURL,
		score:      o.Score,
		atomType:   o.AtomType,
		metadata:   o.Metadata,
	}
	a.SetTitles(o.Titles)
	return a
}

func (a *A) APIType() epb.API     { return a.apiType }
func (a *A) APIID() string        { return a.apiID }
func (a *A) PreviewURL() string   { return a.previewURL }
func (a *A) Score() int64         { return a.score }
func (a *A) AtomType() epb.Type   { return a.atomType }
func (a *A) Metadata() metadata.M { return a.metadata.Copy() }
func (a *A) Copy() *A {
	return New(O{
		APIType:    a.apiType,
		APIID:      a.apiID,
		PreviewURL: a.previewURL,
		Score:      a.score,
		AtomType:   a.atomType,
		Metadata:   a.metadata.Copy(),
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
func (a *A) SetScore(v int64)       { a.score = v }
func (a *A) SetMetadata(v metadata.M) {
	if a.atomType != v.AtomType() {
		panic(fmt.Errorf("cannot set mismatching atom types: %v != %v", a.atomType, v.AtomType()))
	}
	a.metadata = v.Copy()
}

// Merge will combine two atoms, with the following heuristic --
//
//  1. the primitives of two merged atoms will be overwritten by the second
//     atom
//  2. slices and maps of the atoms are a union of the two inputs
//  3. structs (i.e. a.M()) will be recursively merged with the same
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
		Metadata:   a.metadata.Merge(o.metadata),
	})
}
