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
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"

	mpb "github.com/minkezhang/bene-api/proto/go/atom/metadata"
	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

var (
	_ metadata.M = M{}
)

type M struct{}

func (a M) AtomType() epb.Type                { return epb.Type_TYPE_NONE }
func (a M) Equal(o metadata.M) bool           { return true }
func (a M) Copy() metadata.M                  { return M{} }
func (a M) Merge(o metadata.M) metadata.M     { return o.Copy() }
func (a M) Unmarshal() (proto.Message, error) { return &mpb.Empty{}, nil }

func (a M) Marshal() ([]byte, error) {
	pb, err := a.Unmarshal()
	if err != nil {
		return nil, err
	}

	return prototext.Marshal(pb.(*mpb.Empty))
}
