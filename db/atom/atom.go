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

	"github.com/minkezhang/bene-api/db/atom/internal/metadata/mock"
	"github.com/minkezhang/bene-api/db/atom/internal/utils/merge"
	"github.com/minkezhang/bene-api/db/atom/metadata"
	"github.com/minkezhang/bene-api/db/atom/metadata/empty"
	"github.com/minkezhang/bene-api/db/atom/metadata/tv"
	"google.golang.org/protobuf/proto"

	apb "github.com/minkezhang/bene-api/proto/go/atom"
	mpb "github.com/minkezhang/bene-api/proto/go/atom/metadata"
	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

func Load(msg proto.Message) *A {
	pb := msg.(*apb.Atom)
	a := New(O{
		AtomType:   pb.GetType(),
		APIType:    pb.GetApi(),
		APIID:      pb.GetId(),
		PreviewURL: pb.GetPreviewUrl(),
		Score:      pb.GetScore(),
		Synopsis:   pb.GetSynopsis(),
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
		a.SetMetadata(tv.G{}.Load(pb.GetMetadataTv()))
	default:
		a.SetMetadata(empty.G{}.Load(pb.GetMetadataEmpty()))
	}

	return a
}

func Save(a *A) proto.Message {
	pb := &apb.Atom{
		Type:       a.AtomType(),
		Api:        a.APIType(),
		Id:         a.APIID(),
		PreviewUrl: a.PreviewURL(),
		Score:      int64(a.Score()),
		Synopsis:   a.Synopsis(),
	}

	for _, t := range a.Titles() {
		pb.Titles = append(pb.GetTitles(), &apb.Title{
			Title:        t.Title,
			Localization: t.Localization,
		})
	}

	switch t := a.AtomType(); t {
	case epb.Type_TYPE_TV:
		pb.Metadata = &apb.Atom_MetadataTv{MetadataTv: tv.G{}.Save(a.Metadata().(*tv.M)).(*mpb.TV)}
	default:
		pb.Metadata = &apb.Atom_MetadataEmpty{MetadataEmpty: empty.G{}.Save(a.Metadata().(empty.M)).(*mpb.Empty)}
	}

	return pb
}

// Merge will combine two atoms, with the following heuristic --
//
//  1. the primitives of two merged atoms will be overwritten by the higher
//     priority API
//  2. slices and maps of the atoms are a union of the two inputs
//  3. structs (i.e. a.M()) will be recursively merged with the same
//     heuristic
func Merge(a *A, o *A) *A {
	if a.AtomType() != o.AtomType() {
		panic(fmt.Errorf("cannot merge mismatching atom types: %v != %v", a.atomType, o.atomType))
	}
	return New(O{
		APIType: merge_utils.Prioritize(
			merge_utils.V[epb.API]{API: a.APIType(), V: a.APIType()},
			merge_utils.V[epb.API]{API: o.APIType(), V: o.APIType()},
		),
		APIID: merge_utils.Prioritize(
			merge_utils.V[string]{API: a.APIType(), V: a.APIID()},
			merge_utils.V[string]{API: o.APIType(), V: o.APIID()},
		),
		Titles: merge_utils.Distinct(a.Titles(), o.Titles()),
		PreviewURL: merge_utils.Prioritize(
			merge_utils.V[string]{API: a.APIType(), V: a.PreviewURL()},
			merge_utils.V[string]{API: o.APIType(), V: o.PreviewURL()},
		),
		Score: merge_utils.Prioritize(
			merge_utils.V[int64]{API: a.APIType(), V: a.Score()},
			merge_utils.V[int64]{API: o.APIType(), V: o.Score()},
		),
		AtomType: a.AtomType(),
		Synopsis: merge_utils.Prioritize(
			merge_utils.V[string]{API: a.APIType(), V: a.Synopsis()},
			merge_utils.V[string]{API: o.APIType(), V: o.Synopsis()},
		),
		Metadata: MergeMetadata(
			metadata.T{
				API: a.APIType(),
				M:   a.metadata,
			},
			metadata.T{
				API: o.APIType(),
				M:   o.metadata,
			},
		),
	})
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
	Synopsis   string
	AtomType   epb.Type
	Metadata   metadata.M
}

type A struct {
	apiType    epb.API                           // Read-only
	apiID      string                            // Read-only
	titles     map[string]map[string]interface{} // e.g. a.titles["us-en"]["Firefly"]
	previewURL string
	score      int64
	synopsis   string

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
		synopsis:   o.Synopsis,
	}
	a.SetTitles(o.Titles)
	return a
}

func (a *A) APIType() epb.API     { return a.apiType }
func (a *A) APIID() string        { return a.apiID }
func (a *A) PreviewURL() string   { return a.previewURL }
func (a *A) Score() int64         { return a.score }
func (a *A) Synopsis() string     { return a.synopsis }
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
		Synopsis:   a.Synopsis(),
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
func (a *A) SetSynopsis(v string)   { a.synopsis = v }
func (a *A) SetMetadata(v metadata.M) {
	if a.atomType != v.AtomType() {
		panic(fmt.Errorf("cannot set mismatching atom types: %v != %v", a.atomType, v.AtomType()))
	}
	a.metadata = v.Copy()
}

func MergeMetadata(t metadata.T, u metadata.T) metadata.M {
	switch mt := t.M.(type) {
	case *mock.M:
		return mock.G{}.Merge(t, u)
	case empty.M:
		return empty.G{}.Merge(t, u)
	case *tv.M:
		return tv.G{}.Merge(t, u)
	default:
		panic(fmt.Errorf("cannot merge unsupported metadata type: %v", mt))
	}
}
