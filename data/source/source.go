package source

import (
	"cmp"
	"fmt"
	"strings"
	"time"

	"github.com/minkezhang/truffle-api/data/source/util/merge"
	"github.com/minkezhang/truffle-api/util/slice"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	dpb "github.com/minkezhang/truffle-api/proto/go/data"
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

type S struct {
	pb *dpb.Source
}

type T struct {
	_title        string
	_localization string
}

func (t T) Title() string        { return t._title }
func (t T) Localization() string { return t._localization }

func (t T) PB() *dpb.Title {
	return &dpb.Title{
		Title:        t.Title(),
		Localization: t.Localization(),
	}
}

type H struct {
	_api  epb.SourceAPI
	_type epb.SourceType
	_id   string
}

func (h H) API() epb.SourceAPI   { return h._api }
func (h H) Type() epb.SourceType { return h._type }
func (h H) ID() string           { return h._id }

func (h H) PB() *dpb.SourceHeader {
	return &dpb.SourceHeader{
		Api:  h.API(),
		Type: h.Type(),
		Id:   h.ID(),
	}
}

func Make(pb *dpb.Source) S {
	if pb == nil {
		return S{}
	}
	return S{
		pb: clean(pb),
	}
}

func (s S) WithNodeID(node_id string) S {
	pb := s.PB()
	if pb == nil {
		pb = &dpb.Source{}
	}
	pb.NodeId = node_id
	return Make(pb)
}

func (s S) WithHeader(header H) S {
	pb := s.PB()
	if pb == nil {
		pb = &dpb.Source{}
	}
	pb.Header = header.PB()
	return Make(pb)
}

func (s S) Title() T {
	if titles := s.pb.GetTitles(); len(titles) > 0 {
		return T{
			_title:        titles[0].GetTitle(),
			_localization: titles[0].GetLocalization(),
		}
	}
	return T{}
}

func (s S) NodeID() string { return s.PB().GetNodeId() }

func (s S) Header() H {
	return H{
		_api:  s.PB().GetHeader().GetApi(),
		_type: s.PB().GetHeader().GetType(),
		_id:   s.PB().GetHeader().GetId(),
	}
}

func (s S) RelatedHeaders() []H {
	return slice.Apply(s.PB().GetRelatedHeaders(), func(header *dpb.SourceHeader) H {
		return H{
			_api:  header.GetApi(),
			_type: header.GetType(),
			_id:   header.GetId(),
		}
	})
}

func (s S) Titles() []T {
	return slice.Apply(s.PB().GetTitles(), func(title *dpb.Title) T {
		return T{
			_title:        title.GetTitle(),
			_localization: title.GetLocalization(),
		}
	})
}

func (s S) PreviewURL() string       { return s.PB().GetPreviewUrl() }
func (s S) Score() int               { return int(s.PB().GetScore()) }
func (s S) Synopsis() string         { return s.PB().GetSynopsis() }
func (s S) Notes() string            { return s.PB().GetNotes() }
func (s S) Genres() []string         { return s.PB().GetGenres() }
func (s S) Status() epb.SourceStatus { return s.PB().GetStatus() }
func (s S) Studios() []string        { return s.PB().GetStudios() }
func (s S) Seasons() []string        { return s.PB().GetSeasons() }
func (s S) Authors() []string        { return s.PB().GetAuthors() }
func (s S) Illustrators() []string   { return s.PB().GetIllustrators() }
func (s S) LastUpdated() time.Time   { return s.PB().GetLastUpdated().AsTime() }

func (s S) PB() *dpb.Source { return proto.Clone(s.pb).(*dpb.Source) }

func Merge(u, v S) (S, error) {
	if u.Header().Type() != v.Header().Type() {
		return S{}, fmt.Errorf(
			"cannot merge sources of different types: %v != %v",
			u.Header().Type().String(),
			v.Header().Type().String(),
		)
	}

	return Make(&dpb.Source{
		NodeId: merge.Prioritize(u.Header().API(), u.NodeID(), v.Header().API(), v.NodeID()),
		Header: merge.Prioritize(u.Header().API(), u.Header(), v.Header().API(), v.Header()).PB(),
		RelatedHeaders: append(
			slice.Apply(
				u.RelatedHeaders(),
				func(v H) *dpb.SourceHeader { return v.PB() },
			),
			slice.Apply(
				v.RelatedHeaders(),
				func(v H) *dpb.SourceHeader { return v.PB() },
			)...,
		),
		Titles: append(
			slice.Apply(
				u.Titles(),
				func(v T) *dpb.Title { return v.PB() },
			),
			slice.Apply(
				v.Titles(),
				func(v T) *dpb.Title { return v.PB() },
			)...,
		),
		PreviewUrl:   merge.Prioritize(u.Header().API(), u.PreviewURL(), v.Header().API(), v.PreviewURL()),
		Score:        int64(merge.Prioritize(u.Header().API(), u.Score(), v.Header().API(), v.Score())),
		Synopsis:     strings.Join([]string{u.Synopsis(), v.Synopsis()}, "\n\n"),
		Notes:        strings.Join([]string{u.Notes(), v.Notes()}, "\n\n"),
		Genres:       append(u.Genres(), v.Genres()...),
		Status:       merge.Prioritize(u.Header().API(), u.Status(), v.Header().API(), v.Status()),
		Studios:      append(u.Studios(), v.Studios()...),
		Seasons:      append(u.Seasons(), v.Seasons()...),
		Authors:      append(u.Authors(), v.Authors()...),
		Illustrators: append(u.Illustrators(), v.Illustrators()...),
		LastUpdated:  merge.Prioritize(u.Header().API(), timestamppb.New(u.LastUpdated()), v.Header().API(), timestamppb.New(v.LastUpdated())),
	}), nil
}

func clean(src *dpb.Source) *dpb.Source {
	if src == nil {
		return nil
	}
	dst := proto.Clone(src).(*dpb.Source)

	dst.RelatedHeaders = slice.DeduplicateFunc(
		dst.GetRelatedHeaders(),
		func(u, v *dpb.SourceHeader) int { // cmp
			return cmp.Compare(u.GetApi(), v.GetApi())
		},
		func(u, v *dpb.SourceHeader) bool { // eq
			return u.GetApi() == v.GetApi() && u.GetId() == v.GetId()
		},
	)

	dst.Titles = merge.DeduplicateTitles(dst.GetTitles())
	dst.Genres = merge.DeduplicateStrings(dst.GetGenres())
	dst.Studios = merge.DeduplicateStrings(dst.GetStudios())
	dst.Seasons = merge.DeduplicateStrings(dst.GetSeasons())
	dst.Authors = merge.DeduplicateStrings(dst.GetAuthors())
	dst.Illustrators = merge.DeduplicateStrings(dst.GetIllustrators())
	return dst
}
