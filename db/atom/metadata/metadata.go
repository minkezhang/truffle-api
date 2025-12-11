// Package metadata defines the auxiliary media type-specific data stored for
// each atom.
package metadata

import (
	"google.golang.org/protobuf/proto"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

// T is a target Metadata instance along with the AtomType of the container.
type T struct {
	API epb.API
	M   M
}

// G is a generator interface to load and unload a metadata instance.
type G interface {
	Load(msg proto.Message) M
	Save(m M) proto.Message

	// Merge will take in two metadata instances and outputs a metadata
	// instance of the same type. All data in the merged instance must be
	// copy-by-value (including the underlying structs).
	Merge(t T, u T) M
}

type M interface {
	AtomType() epb.Type
	Copy() M
}
