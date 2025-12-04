package metadata

import (
	"github.com/minkezhang/bene-api/db/enums"
	"google.golang.org/protobuf/proto"
)

type M interface {
	AtomType() enums.AtomType
	Copy() M
	Merge(o M) M
	Equal(o M) bool
	Unmarshal() (proto.Message, error)
	Marshal() ([]byte, error)
}
