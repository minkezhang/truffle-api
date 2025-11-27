package node

import (
	dpb "github.com/minkezhang/bene-api/proto/go/data"
)

type N struct {
	parent *N
}

func (n *N) Type() dpb.Type { return n.pb.GetType() }
func (n *N) ID() string     { return n.pb.GetID() }

func (n *N) Parent() *N     { return n.parent }
func (n *N) IsQueued() bool { return n.pb.GetIsQueued() }
