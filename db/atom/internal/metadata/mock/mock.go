package mock

import (
	"fmt"
	"reflect"

	"github.com/minkezhang/bene-api/db/atom/metadata"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"

	mpb "github.com/minkezhang/bene-api/proto/go/atom/metadata"
	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

var (
	_ metadata.M = &M{}
)

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

func (a *M) AtomType() epb.Type      { return epb.Type_TYPE_TV }
func (a *M) Producers() []string     { return append([]string{}, a.producers...) }
func (a *M) Equal(o metadata.M) bool { return reflect.DeepEqual(a, o) }

func (a *M) Copy() metadata.M {
	return &M{
		producers: append([]string{}, a.producers...),
	}
}

func (a *M) Merge(o metadata.M) metadata.M {
	if a.AtomType() != o.AtomType() {
		panic(fmt.Errorf("cannot merge mismatching metadata types: %v != %v", a.AtomType(), o.AtomType()))
	}
	return &M{
		producers: append(
			append([]string{}, a.producers...),
			o.(*M).producers...,
		),
	}
}

func (a M) Unmarshal() (proto.Message, error) {
	return &mpb.Mock{
		Producers: a.Producers(),
	}, nil
}

func (a M) Marshal() ([]byte, error) {
	pb, err := a.Unmarshal()
	if err != nil {
		return nil, err
	}

	return prototext.Marshal(pb.(*mpb.Mock))
}
