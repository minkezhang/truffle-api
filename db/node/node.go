// Package node is a discrete, logical object of interest of a clear media type.
//
// This node encapsulates multiple atoms from potentially multiple data sources.
package node

import (
	"fmt"

	"github.com/minkezhang/bene-api/db/atom"
)

var (
	_ = &N[*atom.Base]{}
)

// N is a generic container which contains multiple atoms of type atom.A[T].
//
// T here is a concrete atom pointer type, e.g. *atom.TV.
type N[T atom.A[T]] struct {
	Type     atom.AtomType
	ID       string
	IsQueued bool
	Atoms    map[atom.ClientAPI]map[string]T
	Related  map[string]bool // Related nodes
}

func (n *N[T]) LinkAtom(a T) error {
	if a.API() == atom.ClientAPINone {
		return fmt.Errorf("atom client API must be specified")
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
	n.Related[m.ID] = true
	m.Related[n.ID] = true
}

func UnlinkNode[T atom.A[T], U atom.A[U]](n *N[T], m *N[U]) {
	delete(n.Related, m.ID)
	delete(m.Related, n.ID)
}

func (n *N[T]) Data() (T, error) {
	var res T
	for _, atoms := range n.Atoms {
		for _, a := range atoms {
			res = res.Merge(a)
		}
	}
	return res, nil
}

