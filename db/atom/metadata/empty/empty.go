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
	"github.com/minkezhang/bene-api/db/atom/metadata"
	"google.golang.org/protobuf/proto"

	mpb "github.com/minkezhang/bene-api/proto/go/atom/metadata"
	epb "github.com/minkezhang/bene-api/proto/go/enums"
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

type M struct{}

func (m M) AtomType() epb.Type            { return epb.Type_TYPE_NONE }
func (m M) Equal(o metadata.M) bool       { return true }
func (m M) Copy() metadata.M              { return M{} }
func (m M) Merge(v metadata.M) metadata.M { return v.Copy() }
