package tv

import (
	"fmt"
	"reflect"

	"github.com/minkezhang/bene-api/db/atom/metadata"
	"github.com/minkezhang/bene-api/db/atom/metadata/shared/video"
	"google.golang.org/protobuf/proto"

	mpb "github.com/minkezhang/bene-api/proto/go/atom/metadata"
	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

var (
	_ metadata.M = &M{}
	_ metadata.G = G{}
)

type G struct{}

func (g G) Load(msg proto.Message) metadata.M {
	return &M{
		M: video.G{}.Load(msg.(*mpb.TV).GetVideo()).(*video.M),
	}
}

func (g G) Save(m metadata.M) proto.Message {
	return &mpb.TV{
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

type O video.O

type M struct {
	*video.M
}

func New(o O) *M {
	return &M{
		M: video.New(video.O{
			Genres:      o.Genres,
			Showrunners: o.Showrunners,
			IsAnimated:  o.IsAnimated,
			IsAnime:     o.IsAnime,
			Studios:     o.Studios,
			Networks:    o.Networks,
		}),
	}
}

func (m *M) AtomType() epb.Type      { return epb.Type_TYPE_TV }
func (m *M) Equal(o metadata.M) bool { return reflect.DeepEqual(m, o) }
func (m *M) Copy() metadata.M        { return &M{M: m.M.Copy().(*video.M)} }
