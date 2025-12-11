package movie

import (
	"fmt"

	"github.com/minkezhang/truffle-api/db/atom/metadata"
	"github.com/minkezhang/truffle-api/db/atom/metadata/shared/video"
	"github.com/minkezhang/truffle-api/db/atom/metadata/tv"
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
	return &M{
		M: video.G{}.Load(msg.(*mpb.Movie).GetVideo()).(*video.M),
	}
}

func (g G) Save(m metadata.M) proto.Message {
	return &mpb.Movie{
		Video: video.G{}.Save(m.(*M).M).(*mpb.Video),
	}
}

func (g G) Merge(t metadata.T, u metadata.T) metadata.M {
	if t.M.AtomType() != u.M.AtomType() {
		panic(fmt.Errorf("cannot merge mismatching metadata types: %v != %v", t.M.AtomType(), u.M.AtomType()))
	}

	return &M{
		M: video.G{}.Merge(
			metadata.T{
				API: t.API,
				M:   t.M.(*M).M,
			},
			metadata.T{
				API: u.API,
				M:   u.M.(*M).M,
			},
		).(*video.M),
	}
}

type M tv.M
type O tv.O

func New(o O) *M { return (*M)(tv.New(tv.O(o))) }

func (m *M) AtomType() epb.Type { return epb.Type_TYPE_MOVIE }
