package tv

import (
	"fmt"
	"reflect"

	"github.com/minkezhang/bene-api/db/atom/internal/utils/merge"
	"github.com/minkezhang/bene-api/db/atom/metadata"
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
	pb := msg.(*mpb.TV)
	return New(O{
		Genres:      pb.GetGenres(),
		Showrunners: pb.GetShowrunners(),
		IsAnimated:  pb.GetIsAnimated(),
		IsAnime:     pb.GetIsAnime(),
		Studios:     pb.GetStudios(),
		Networks:    pb.GetNetworks(),
	})
}

func (g G) Save(m metadata.M) proto.Message {
	md := m.(*M)
	return &mpb.TV{
		Genres:      md.Genres(),
		Showrunners: md.Showrunners(),
		IsAnimated:  md.IsAnimated(),
		IsAnime:     md.IsAnime(),
		Studios:     md.Studios(),
		Networks:    md.Networks(),
	}
}

func (g G) Merge(t metadata.T, u metadata.T) metadata.M {
	if t.M.AtomType() != u.M.AtomType() {
		panic(fmt.Errorf("cannot merge mismatching metadata types: %v != %v", t.M.AtomType(), u.M.AtomType()))
	}

	return New(O{
		Genres:      merge_utils.Distinct(t.M.(*M).Genres(), u.M.(*M).Genres()),
		Showrunners: merge_utils.Distinct(t.M.(*M).Showrunners(), u.M.(*M).Showrunners()),
		IsAnimated: merge_utils.Prioritize(
			merge_utils.V[bool]{API: t.API, V: t.M.(*M).IsAnimated()},
			merge_utils.V[bool]{API: u.API, V: u.M.(*M).IsAnimated()},
		),
		IsAnime: merge_utils.Prioritize(
			merge_utils.V[bool]{API: t.API, V: t.M.(*M).IsAnime()},
			merge_utils.V[bool]{API: u.API, V: u.M.(*M).IsAnime()},
		),
		Studios:  merge_utils.Distinct(t.M.(*M).Studios(), u.M.(*M).Studios()),
		Networks: merge_utils.Distinct(t.M.(*M).Networks(), u.M.(*M).Networks()),
	})
}

type O struct {
	Genres      []string
	Showrunners []string
	IsAnimated  bool
	IsAnime     bool
	Studios     []string
	Networks    []string
}

type M struct {
	genres      []string
	showrunners []string
	isAnimated  bool
	isAnime     bool
	studios     []string
	networks    []string
}

func New(o O) *M {
	return &M{
		genres:      append([]string{}, o.Genres...),
		showrunners: append([]string{}, o.Showrunners...),
		isAnimated:  o.IsAnimated,
		isAnime:     o.IsAnime,
		studios:     append([]string{}, o.Studios...),
		networks:    append([]string{}, o.Networks...),
	}
}

func (m *M) AtomType() epb.Type { return epb.Type_TYPE_TV }

func (m *M) Genres() []string      { return append([]string{}, m.genres...) }
func (m *M) Showrunners() []string { return append([]string{}, m.showrunners...) }
func (m *M) IsAnimated() bool      { return m.isAnimated }
func (m *M) IsAnime() bool         { return m.isAnime }
func (m *M) Studios() []string     { return append([]string{}, m.studios...) }
func (m *M) Networks() []string    { return append([]string{}, m.networks...) }

func (m *M) SetGenres(v []string)      { m.genres = append([]string{}, v...) }
func (m *M) SetShowrunners(v []string) { m.showrunners = append([]string{}, v...) }
func (m *M) SetIsAnimated(v bool)      { m.isAnimated = v }
func (m *M) SetIsAnime(v bool)         { m.isAnime = v }
func (m *M) SetStudios(v []string)     { m.studios = append([]string{}, v...) }
func (m *M) SetNetworks(v []string)    { m.studios = append([]string{}, v...) }

func (m *M) Equal(o metadata.M) bool { return reflect.DeepEqual(m, o) }

func (m *M) Copy() metadata.M {
	return New(O{
		Genres:      m.Genres(),
		Showrunners: m.Showrunners(),
		IsAnimated:  m.IsAnimated(),
		IsAnime:     m.IsAnime(),
		Studios:     m.Studios(),
		Networks:    m.Networks(),
	})
}
