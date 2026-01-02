package truffle

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/minkezhang/truffle-api/data/source"
	"google.golang.org/protobuf/testing/protocmp"

	cpb "github.com/minkezhang/truffle-api/proto/go/config"
	dpb "github.com/minkezhang/truffle-api/proto/go/data"
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

var (
	frieren = &dpb.Source{
		Header: &dpb.SourceHeader{
			Api:  epb.SourceAPI_SOURCE_API_MAL,
			Type: epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
			Id:   "foo",
		},
		Titles: []*dpb.Title{
			&dpb.Title{Title: "Frieren"},
		},
	}
	conan = &dpb.Source{
		Header: &dpb.SourceHeader{
			Api:  epb.SourceAPI_SOURCE_API_MAL,
			Type: epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
			Id:   "bar",
		},
		Titles: []*dpb.Title{
			&dpb.Title{Title: "Meitantei Conan"},
		},
	}

	sources = []*dpb.Source{frieren, conan}
)

func TestSearch(t *testing.T) {
	c := New(&cpb.Truffle{}, sources)
	got, err := c.Search(context.Background(), "conan")
	if err != nil {
		t.Errorf("Search() returned a non-nil error: %v", err)
	}
	want := []source.S{source.Make(conan)}
	if diff := cmp.Diff(
		want,
		got,
		cmp.AllowUnexported(source.S{}),
		protocmp.Transform(),
	); diff != "" {
		t.Errorf("Search() mismatch (-want +got):\n%v", diff)
	}
}

func TestPut(t *testing.T) {
	c := New(&cpb.Truffle{}, nil)

	s := source.Make(
		&dpb.Source{
			Header: &dpb.SourceHeader{
				Type: epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
			},
			Titles: []*dpb.Title{
				&dpb.Title{Title: "Meitantei Conan"},
			},
		},
	)

	header, _ := c.Put(context.Background(), s)
	if header.ID() == "" {
		t.Errorf("Put() header ID is unexpectedly blank")
	}

	pb := s.PB()
	pb.Header = header.PB()
	want := source.Make(pb)

	got, _ := c.Get(context.Background(), header)
	if diff := cmp.Diff(
		want,
		got,
		cmp.AllowUnexported(source.S{}),
		protocmp.Transform(),
	); diff != "" {
		t.Errorf("Get() mismatch (-want +got):\n%v", diff)
	}
}
