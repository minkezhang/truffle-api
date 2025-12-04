package metadata

import (
	"google.golang.org/protobuf/proto"

	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

type M interface {
	AtomType() epb.Type
	Copy() M
	Merge(o M) M
	Equal(o M) bool
	Unmarshal() (proto.Message, error)
	Marshal() ([]byte, error)
}
