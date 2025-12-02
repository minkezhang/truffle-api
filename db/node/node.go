// Package node is a discrete, logical object of interest of a clear media type.
//
// This node encapsulates multiple atoms from potentially multiple data sources.
package node

import (
	"fmt"

	"github.com/minkezhang/bene-api/db/atom"
)

type O[T atom.A[T]] struct {
	Type     atom.AtomType
	ID       string
	IsQueued bool
	Atoms    []T
	Related  []string
}

func DebugNewOrDie[T atom.A[T]](o O[T]) *N[T] {
	n, err := New(o)
	if err != nil {
		panic(err)
	}
	return n
}

func New[T atom.A[T]](o O[T]) (*N[T], error) {
	n := &N[T]{
		t:        o.Type,
		id:       o.ID,
		IsQueued: o.IsQueued,
		Atoms:    map[atom.ClientAPI]map[string]T{}, // e.g. n.Atoms[ClientAPIBene]["foo"]
		Related:  map[string]interface{}{},
	}
	for _, a := range o.Atoms {
		if err := n.LinkAtom(a); err != nil {
			return nil, err
		}
	}
	for _, l := range o.Related {
		n.Related[l] = struct{}{}
	}
	return n, nil
}

// N is a generic container which contains multiple atoms of type atom.A[T].
//
// T here is a concrete atom pointer type, e.g. *atom.TV.
type N[T atom.A[T]] struct {
	id string
	t  atom.AtomType

	IsQueued bool
	Atoms    map[atom.ClientAPI]map[string]T
	Related  map[string]interface{} // Related nodes
}

func (n *N[T]) ID() string          { return n.id }
func (n *N[T]) Type() atom.AtomType { return n.t }

func (n *N[T]) LinkAtom(a T) error {
	if a.API() == atom.ClientAPIVirtual {
		return fmt.Errorf("atom client API must be non-virtual")
	}
	if a.ID() == "" {
		return fmt.Errorf("atom ID must be non-empty")
	}

	if _, ok := n.Atoms[a.API()]; !ok {
		n.Atoms[a.API()] = map[string]T{}
	}
	n.Atoms[a.API()][a.ID()] = a
	return nil
}

func (n *N[A]) UnlinkAtom(api atom.ClientAPI, id string) error {
	if _, ok := n.Atoms[api]; ok {
		delete(n.Atoms[api], id)
	}
	return nil
}

func LinkNode[T atom.A[T], U atom.A[U]](n *N[T], m *N[U]) {
	n.Related[m.ID()] = struct{}{}
	m.Related[n.ID()] = struct{}{}
}

func UnlinkNode[T atom.A[T], U atom.A[U]](n *N[T], m *N[U]) {
	delete(n.Related, m.ID())
	delete(m.Related, n.ID())
}

// Union merges all atoms into a representative struct.
func (n *N[T]) Union() (T, error) {
	var res T
	for _, atoms := range n.Atoms {
		for _, a := range atoms {
			res = res.Merge(a)
		}
	}
	return res, nil
}
