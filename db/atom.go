package db

import (
	dpb "github.com/minkezhang/bene-api/proto/go/data"
)

type Atom struct {
	db     *DB
	nodeID string
	atomID string

	pb *dpb.Atom
}

func (a *Atom) Node() *Node            { return a.nodes.Get(a.nodeID) }
func (a *Atom) NodeType() dpb.NodeType { return a.Node().Type() }
