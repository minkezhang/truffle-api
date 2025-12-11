package book

import (
	"fmt"

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
	pb := msg.(*mpb.Book)
	return New(O{
		Genres:        pb.GetGenres(),
		Authors:       pb.GetAuthors(),
		Illustrators:  pb.GetIllustrators(),
		IsIllustrated: pb.GetIsIllustrated(),
		IsManga:       pb.GetIsManga(),
	})
}

func (g G) Save(m metadata.M) proto.Message {
	md := m.(*M)
	return &mpb.Book{
		Genres:        md.Genres(),
		Authors:       md.Authors(),
		Illustrators:  md.Illustrators(),
		IsIllustrated: md.IsIllustrated(),
		IsManga:       md.IsManga(),
	}
}

func (g G) Merge(t metadata.T, u metadata.T) metadata.M {
	if t.M.AtomType() != u.M.AtomType() {
		panic(fmt.Errorf("cannot merge mismatching metadata types: %v != %v", t.M.AtomType(), u.M.AtomType()))
	}

	return New(O{
		Genres:       merge_utils.Distinct(t.M.(*M).Genres(), u.M.(*M).Genres()),
		Authors:      merge_utils.Distinct(t.M.(*M).Authors(), u.M.(*M).Authors()),
		Illustrators: merge_utils.Distinct(t.M.(*M).Illustrators(), u.M.(*M).Illustrators()),
		IsIllustrated: merge_utils.Prioritize(
			merge_utils.V[bool]{API: t.API, V: t.M.(*M).IsIllustrated()},
			merge_utils.V[bool]{API: u.API, V: u.M.(*M).IsIllustrated()},
		),
		IsManga: merge_utils.Prioritize(
			merge_utils.V[bool]{API: t.API, V: t.M.(*M).IsManga()},
			merge_utils.V[bool]{API: u.API, V: u.M.(*M).IsManga()},
		),
	})
}

type O struct {
	Genres        []string
	Authors       []string
	Illustrators  []string
	IsIllustrated bool
	IsManga       bool
}

func New(o O) *M {
	return &M{
		genres:        o.Genres,
		authors:       o.Authors,
		illustrators:  o.Illustrators,
		isIllustrated: o.IsIllustrated,
		isManga:       o.IsManga,
	}
}

type M struct {
	genres        []string
	authors       []string
	illustrators  []string
	isIllustrated bool
	isManga       bool
}

func (m *M) Genres() []string       { return append([]string{}, m.genres...) }
func (m *M) Authors() []string      { return append([]string{}, m.authors...) }
func (m *M) Illustrators() []string { return append([]string{}, m.illustrators...) }
func (m *M) IsIllustrated() bool    { return m.isIllustrated }
func (m *M) IsManga() bool          { return m.isManga }

func (m *M) SetGenres(v []string)       { m.genres = append([]string{}, v...) }
func (m *M) SetAuthors(v []string)      { m.authors = append([]string{}, v...) }
func (m *M) SetIllustrators(v []string) { m.illustrators = append([]string{}, v...) }
func (m *M) SetIsIllustrated(v bool)    { m.isIllustrated = v }
func (m *M) SetIsManga(v bool)          { m.isManga = v }

func (m *M) AtomType() epb.Type { return epb.Type_TYPE_BOOK }

func (m *M) Copy() metadata.M {
	return New(O{
		Genres:        m.Genres(),
		Authors:       m.Authors(),
		Illustrators:  m.Illustrators(),
		IsIllustrated: m.IsIllustrated(),
		IsManga:       m.IsManga(),
	})
}
