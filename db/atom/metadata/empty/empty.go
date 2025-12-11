// Package empty is used to initialize an `atom.M` without any specific
// metadatailiary data.
//
// Example:
//
//	atom.New(atom.O{
//		...
//		Metadata: empty.M{},
//	})
package empty

import (
	"github.com/minkezhang/truffle-api/db/atom/metadata"
	"google.golang.org/protobuf/proto"

	mpb "github.com/minkezhang/truffle-api/proto/go/atom/metadata"
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

var (
	_ metadata.M = M{}
	_ metadata.G = G{}
)

type G struct{}

func (g G) Load(msg proto.Message) metadata.M {
	_ = msg.(*mpb.Empty)
	return M{}
}

func (g G) Save(m metadata.M) proto.Message { return &mpb.Empty{} }

func (g G) Merge(t metadata.T, u metadata.T) metadata.M { return u.M.Copy() }

type M struct{}

func (m M) AtomType() epb.Type { return epb.Type_TYPE_NONE }
func (m M) Copy() metadata.M   { return M{} }
