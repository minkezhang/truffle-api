package tv

import (
	"fmt"
	"reflect"

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
	_ = msg.(*mpb.TV)
	return &M{}
}

func (g G) Save(m metadata.M) proto.Message { return &mpb.TV{} }

type O struct {
	Seasons        []int64
	Genres         []string
	Showrunners    []string
	IsAnimated     bool
	Directors      []string
	Writers        []string
	Cinematography []string
	Composers      []string
	Studios        []string
}

type M struct {
	seasons        []int64
	genres         []string
	showrunners    []string
	isAnimated     bool
	directors      []string
	writers        []string
	cinematography []string
	composers      []string
	studios        []string
}

func New(o O) *M {
	return &M{
		seasons:        append([]int64{}, o.Seasons...),
		genres:         append([]string{}, o.Genres...),
		showrunners:    append([]string{}, o.Showrunners...),
		isAnimated:     o.IsAnimated,
		directors:      append([]string{}, o.Directors...),
		writers:        append([]string{}, o.Writers...),
		cinematography: append([]string{}, o.Cinematography...),
		composers:      append([]string{}, o.Composers...),
		studios:        append([]string{}, o.Studios...),
	}
}

func (m *M) AtomType() epb.Type { return epb.Type_TYPE_TV }

func (m *M) Seasons() []int64         { return append([]int64{}, m.seasons...) }
func (m *M) Genres() []string         { return append([]string{}, m.genres...) }
func (m *M) Showrunners() []string    { return append([]string{}, m.showrunners...) }
func (m *M) IsAnimated() bool         { return m.isAnimated }
func (m *M) Directors() []string      { return append([]string{}, m.directors...) }
func (m *M) Writers() []string        { return append([]string{}, m.writers...) }
func (m *M) Cinematography() []string { return append([]string{}, m.cinematography...) }
func (m *M) Composers() []string      { return append([]string{}, m.composers...) }
func (m *M) Studios() []string        { return append([]string{}, m.studios...) }

func (m *M) SetSeasons(v []int64)         { m.seasons = append([]int64{}, v...) }
func (m *M) SetGenres(v []string)         { m.genres = append([]string{}, v...) }
func (m *M) SetShowrunners(v []string)    { m.showrunners = append([]string{}, v...) }
func (m *M) SetIsAnimated(v bool)         { m.isAnimated = v }
func (m *M) SetDirectors(v []string)      { m.directors = append([]string{}, v...) }
func (m *M) SetWriters(v []string)        { m.writers = append([]string{}, v...) }
func (m *M) SetCinematography(v []string) { m.cinematography = append([]string{}, v...) }
func (m *M) SetComposers(v []string)      { m.composers = append([]string{}, v...) }
func (m *M) SetStudios(v []string)        { m.studios = append([]string{}, v...) }

func (m *M) Equal(o metadata.M) bool { return reflect.DeepEqual(m, o) }

func (m *M) Copy() metadata.M {
	return New(O{
		Seasons:        m.Seasons(),
		Genres:         m.Genres(),
		Showrunners:    m.Genres(),
		IsAnimated:     m.IsAnimated(),
		Directors:      m.Directors(),
		Writers:        m.Writers(),
		Cinematography: m.Cinematography(),
		Composers:      m.Composers(),
		Studios:        m.Studios(),
	})
}

func (m *M) Merge(v metadata.M) metadata.M {
	if m.AtomType() != v.AtomType() {
		panic(fmt.Errorf("cannot merge mismatching metadata types: %v != %v", m.AtomType(), v.AtomType()))
	}

	o := v.(*M)

	return New(O{
		Seasons:        append(m.Seasons(), o.Seasons()...),
		Genres:         append(m.Genres(), o.Genres()...),
		Showrunners:    append(m.Showrunners(), o.Showrunners()...),
		IsAnimated:     o.IsAnimated(),
		Directors:      append(m.Directors(), o.Directors()...),
		Writers:        append(m.Writers(), o.Writers()...),
		Cinematography: append(m.Cinematography(), o.Cinematography()...),
		Composers:      append(m.Composers(), o.Composers()...),
		Studios:        append(m.Studios(), o.Studios()...),
	})
}
