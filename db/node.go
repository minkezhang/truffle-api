package db

import (
	dpb "github.com/minkezhang/bene-api/proto/go/data"
)

type Node struct {
	db *DB
	pb *dpb.Node
}

func (n *Node) Type() dpb.NodeType { return n.pb.GetType() }
func (n *Node) IsQueued() bool     { return n.pb.GetIsQueued() }
func (n *Node) SetIsQueued(v bool) { n.pb.SetIsQueued(v) }

func (n *Node) Atoms() []*Atom {
	return nil
}

func (n *Node) Related() []*Node {
	ns := []*Node{}
	for _, m := range pb.GetRelated() {
		ns = append(ns, n.nodes.Get(m.GetId()))
	}
	return ns
}
