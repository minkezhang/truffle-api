// Package node represents a logically distinct collection of atoms.
//
// That is, a node may be backed by multiple data sources, e.g. both IMDB and
// MAL have entries for Our War Game! --
//
//	https://www.imdb.com/title/tt0313441/
//	https://myanimelist.net/anime/2397/
package node

import (
	"fmt"

	"github.com/minkezhang/truffle-api/db/atom"
	"github.com/minkezhang/truffle-api/db/atom/metadata/empty"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

type O struct {
	ID              string
	IsAuthoritative bool
	AtomType        epb.Type
	IsQueued        bool
	Notes           string
	Atoms           []*atom.A
}

type N struct {
	id       string   // Read-only
	atomType epb.Type // Read-only

	isAuthoritative bool // Is saved locally
	isQueued        bool // Is starred by the user
	notes           string
	atoms           map[epb.API]map[string]*atom.A
}

func New(o O) *N {
	n := &N{
		id:              o.ID,
		atomType:        o.AtomType,
		isAuthoritative: o.IsAuthoritative,
		notes:           o.Notes,
		isQueued:        o.IsQueued,
	}
	n.SetAtoms(o.Atoms)
	return n
}

func (n *N) ID() string            { return n.id }
func (n *N) AtomType() epb.Type    { return n.atomType }
func (n *N) IsQueued() bool        { return n.isQueued }
func (n *N) Notes() bool           { return n.notes }
func (n *N) IsAuthoritative() bool { return n.isAuthoritative }

func (n *N) Atoms() []*atom.A {
	res := []*atom.A{}
	for _, as := range n.atoms {
		for _, a := range as {
			res = append(res, a)
		}
	}
	return res
}

func (n *N) SetAtoms(as []*atom.A) {
	n.atoms = map[epb.API]map[string]*atom.A{}
	for _, a := range as {
		n.AddAtom(a)
	}
}

func (n *N) SetIsQueued(v bool)        { n.isQueued = v }
func (n *N) SetIsAuthoritative(v bool) { n.isAuthoritative = v }
func (n *N) SetNotes(v string)         { n.notes = v }

func (n *N) AddAtom(a *atom.A) {
	if n.atomType != a.AtomType() {
		panic(fmt.Errorf("cannot link atom with mismatching type: %v != %v", n.atomType, a.AtomType()))
	}
	if _, ok := n.atoms[a.APIType()]; !ok {
		n.atoms[a.APIType()] = map[string]*atom.A{}
	}
	n.atoms[a.APIType()][a.APIID()] = a.Copy()
}

func (n *N) RemoveAtom(api epb.API, id string) {
	if vs, ok := n.atoms[api]; ok {
		delete(vs, id)
	}
}

func (n *N) Copy() *N {
	return New(O{
		ID:              n.ID(),
		AtomType:        n.AtomType(),
		IsAuthoritative: n.IsAuthoritative(),
		IsQueued:        n.IsQueued(),
		Notes:           n.Notes(),
		Atoms:           n.Atoms(),
	})
}

// Virtual returns the merged data of all atoms encapsulated in this node.
func (n *N) Virtual() *atom.A {
	res := atom.New(atom.O{
		APIType:  epb.API_API_VIRTUAL,
		AtomType: n.atomType,
		Metadata: empty.M{},
	})
	for _, a := range n.Atoms() {
		res = atom.Merge(res, a)
	}
	return atom.New(atom.O{
		APIType:    epb.API_API_VIRTUAL,
		APIID:      "",
		Titles:     res.Titles(),
		PreviewURL: res.PreviewURL(),
		Synopsis:   res.Synopsis(),
		Score:      res.Score(),
		AtomType:   res.AtomType(),
		Metadata:   res.Metadata(),
	})
}
