// Package node is a discrete, logical object of interest of a clear media type.
//
// This node encapsulates multiple atoms from potentially multiple data sources.
package node

import (
	"fmt"

	"github.com/minkezhang/bene-api/db/atom"
)

type N[A atom.A] struct {
	Type     atom.AtomType
	ID       string
	IsQueued bool
	Atoms    map[atom.ClientAPI]map[string]A
	Related  map[string]bool // Related nodes
}

func (n *N[A]) CreateAtom(a A) error {
	if a.GetType() == atom.AtomTypeNone {
		return fmt.Errorf("atom type must be specified")
	}
	if a.GetAPI() == atom.ClientAPINone {
		return fmt.Errorf("atom client API must be specified")
	}
	if a.GetType() != n.Type {
		return fmt.Errorf("invalid atom type: %v != %v", a.GetType(), n.Type)
	}
	if a.GetID() == "" {
		return fmt.Errorf("atom ID must be non-empty")
	}

	if _, ok := n.Atoms[a.GetAPI()]; !ok {
		n.Atoms[a.GetAPI()] = map[string]A{}
	}
	n.Atoms[a.GetAPI()][a.GetID()] = a
	return nil
}

func (n *N[A]) GetData() (A, error) {
	var a A
	for _, atoms := range n.Atoms {
		for _, b := range atoms {
			c, err := a.Merge(b)
			if err != nil { return a, err }
			a = c.(A)
		}
	}
	return a, nil
}

func (n *N[A]) DeleteAtom(api atom.ClientAPI, id string) error {
	if _, ok := n.Atoms[api]; ok {
		delete(n.Atoms[api], id)
	}
	return nil
}

func CreateLink[A atom.A, B atom.A](n *N[A], m *N[B]) {
	n.Related[m.ID] = true
	m.Related[n.ID] = true
}

func DeleteLink[A atom.A, B atom.A](n *N[A], m *N[B]) {
	delete(n.Related, m.ID)
	delete(m.Related, n.ID)
}
