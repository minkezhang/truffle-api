package mock

import (
	"fmt"
	"reflect"

	"github.com/minkezhang/truffle-api/db/atom/internal/utils/merge"
	"github.com/minkezhang/truffle-api/db/atom/metadata"
	"google.golang.org/protobuf/proto"

	mpb "github.com/minkezhang/truffle-api/proto/go/atom/metadata"
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

var (
	_ metadata.M = &M{}
	_ metadata.G = G{}
)

type G struct{}

func (g G) Load(msg proto.Message) metadata.M {
	return New(O{
		Producers: msg.(*mpb.Mock).GetProducers(),
	})
}

func (g G) Save(m metadata.M) proto.Message {
	return &mpb.Mock{
		Producers: m.(*M).Producers(),
	}
}

func (g G) Merge(t metadata.T, u metadata.T) metadata.M {
	if t.M.AtomType() != u.M.AtomType() {
		panic(fmt.Errorf("cannot merge mismatching metadata types: %v != %v", t.M.AtomType(), u.M.AtomType()))
	}
	return New(O{
		Producers: merge_utils.Distinct(t.M.(*M).Producers(), u.M.(*M).Producers()),
	})
}

type O struct {
	Producers []string
}

func New(o O) *M {
	return &M{
		producers: append([]string{}, o.Producers...),
	}
}

type M struct {
	producers []string
}

func (m *M) AtomType() epb.Type      { return epb.Type_TYPE_TV }
func (m *M) Producers() []string     { return append([]string{}, m.producers...) }
func (m *M) Equal(o metadata.M) bool { return reflect.DeepEqual(m, o) }

func (m *M) Copy() metadata.M {
	return &M{
		producers: append([]string{}, m.producers...),
	}
}
