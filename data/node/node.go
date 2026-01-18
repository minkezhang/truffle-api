package node

import (
	"github.com/minkezhang/truffle-api/data/source"
	"google.golang.org/protobuf/proto"

	dpb "github.com/minkezhang/truffle-api/proto/go/data"
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

func Make(pb *dpb.Node) N {
	return N{
		pb: pb,
	}
}

type H struct {
	_type epb.SourceType
	_id   string
}

func (h H) ID() string           { return h._id }
func (h H) Type() epb.SourceType { return h._type }

func (h H) PB() *dpb.NodeHeader {
	return &dpb.NodeHeader{
		Type: h.Type(),
		Id:   h.ID(),
	}
}

type N struct {
	pb      *dpb.Node
	sources []source.S // Virtual only
}

func (n N) Header() H {
	return H{
		_type: n.PB().GetHeader().GetType(),
		_id:   n.PB().GetHeader().GetId(),
	}
}

func (n N) Sources() []source.S { return n.sources }
func (n N) PB() *dpb.Node       { return proto.Clone(n.pb).(*dpb.Node) }

func (n N) WithSources(vs []source.S) N {
	pb := n.PB()
	if pb == nil {
		pb = &dpb.Node{}
	}
	m := Make(pb)
	m.sources = vs
	return m
}

func (n N) Virtual() (source.S, error) {
	res := source.Make(&dpb.Source{
		NodeId: n.Header().ID(),
		Header: &dpb.SourceHeader{
			Type: n.Header().Type(),
			Api: epb.SourceAPI_SOURCE_API_TRUFFLE,
		},
	})
	for _, s := range n.sources {
		var err error
		res, err = source.Merge(res, s)
		if err != nil {
			return source.S{}, err
		}
	}
	return res, nil
}
