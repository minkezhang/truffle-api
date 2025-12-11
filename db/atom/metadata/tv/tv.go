package tv

import (
	"fmt"
	"time"

	"github.com/minkezhang/truffle-api/db/atom/internal/utils/merge"
	"github.com/minkezhang/truffle-api/db/atom/metadata"
	"github.com/minkezhang/truffle-api/db/atom/metadata/movie"
	"github.com/minkezhang/truffle-api/db/atom/metadata/shared/video"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

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
		M:           video.G{}.Load(msg.(*mpb.TV).GetVideo()).(*M).M,
		lastUpdated: msg.(*mpb.TV).GetLastUpdated().AsTime(),
	}
}

func (g G) Save(m metadata.M) proto.Message {
	return &mpb.TV{
		Video:       video.G{}.Save(m.(*M).M).(*mpb.Video),
		LastUpdated: timestamppb.New(m.(*M).LastUpdated()),
	}
}

func (g G) Merge(t metadata.T, u metadata.T) metadata.M {
	if t.M.AtomType() != u.M.AtomType() {
		panic(fmt.Errorf("cannot merge mismatching metadata types: %v != %v", t.M.AtomType(), u.M.AtomType()))
	}

	return &M{
		M: movie.G{}.Merge(
			metadata.T{
				API: t.API,
				M:   t.M.(*M).M,
			},
			metadata.T{
				API: u.API,
				M:   u.M.(*M).M,
			},
		).(*movie.M),
		lastUpdated: merge_utils.Prioritize(
			merge_utils.V[time.Time]{API: t.API, V: t.M.(*M).LastUpdated()},
			merge_utils.V[time.Time]{API: u.API, V: u.M.(*M).LastUpdated()},
		),
	}
}

type O struct {
	movie.O

	LastUpdated time.Time
}

type M struct {
	*movie.M

	lastUpdated time.Time
}

func New(o O) *M {
	return &M{
		M:           movie.New(o.O),
		lastUpdated: o.LastUpdated,
	}
}

func (m *M) LastUpdated() time.Time { return m.lastUpdated }

func (m *M) SetLastUpdated(v time.Time) { m.lastUpdated = v }

func (m *M) AtomType() epb.Type { return epb.Type_TYPE_TV }

func (m *M) Copy() metadata.M {
	return &M{
		M:           m.M.Copy().(*movie.M),
		lastUpdated: m.LastUpdated(),
	}
}
