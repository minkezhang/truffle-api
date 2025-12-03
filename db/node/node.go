package node

import (
	"fmt"

	"github.com/minkezhang/bene-api/db/atom"
	"github.com/minkezhang/bene-api/db/enums"
)

type O struct {
	ID              string
	IsAuthoritative bool
	AtomType        enums.AtomType
	IsQueued        bool
	Atoms           []*atom.A
}

type N struct {
	id       string         // Read-only
	atomType enums.AtomType // Read-only

	isAuthoritative bool
	isQueued        bool
	atoms           map[enums.ClientAPI]map[string]*atom.A
}

func New(o O) *N {
	n := &N{
		id:              o.ID,
		atomType:        o.AtomType,
		isAuthoritative: o.IsAuthoritative,
		isQueued:        o.IsQueued,
	}
	n.SetAtoms(o.Atoms)
	return n
}

func (n *N) ID() string               { return n.id }
func (n *N) AtomType() enums.AtomType { return n.atomType }
func (n *N) IsQueued() bool           { return n.isQueued }
func (n *N) IsAuthoritative() bool    { return n.isAuthoritative }

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
	n.atoms = map[enums.ClientAPI]map[string]*atom.A{}
	for _, a := range as {
		n.AddAtom(a)
	}
}

func (n *N) SetIsQueued(v bool)        { n.isQueued = v }
func (n *N) SetIsAuthoritative(v bool) { n.isAuthoritative = v }

func (n *N) AddAtom(a *atom.A) {
	if n.atomType != a.AtomType() {
		panic(fmt.Errorf("cannot link atom with mismatching type: %v != %v", n.atomType, a.AtomType()))
	}
	if _, ok := n.atoms[a.APIType()]; !ok {
		n.atoms[a.APIType()] = map[string]*atom.A{}
	}
	n.atoms[a.APIType()][a.APIID()] = a.Copy()
}

func (n *N) RemoveAtom(api enums.ClientAPI, id string) {
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
		Atoms:           n.Atoms(),
	})
}

// Virtual returns the merged data of all atoms encapsulated in this node.
func (n *N) Virtual() *atom.A {
	res := atom.New(atom.O{
		AtomType: n.atomType,
	})
	for _, a := range n.Atoms() {
		res = res.Merge(a)
	}
	return atom.New(atom.O{
		APIType:    res.APIType(),
		APIID:      "",
		Titles:     res.Titles(),
		PreviewURL: res.PreviewURL(),
		Score:      res.Score(),
		AtomType:   res.AtomType(),
		Aux:        res.Aux(),
	})
}
