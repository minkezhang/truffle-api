package metadata

import (
	"google.golang.org/protobuf/proto"

	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

// G is a generator interface to load and unload a metadata instance.
type G interface {
	Load(msg proto.Message) M
	Save(m M) proto.Message
}

type M interface {
	AtomType() epb.Type
	Copy() M
	Merge(o M) M
	Equal(o M) bool
}
