package client

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/minkezhang/truffle-api/client/mock"
	"github.com/minkezhang/truffle-api/client/option"
	"github.com/minkezhang/truffle-api/data/source"
	"google.golang.org/protobuf/testing/protocmp"

	dpb "github.com/minkezhang/truffle-api/proto/go/data"
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

var (
	frieren = &dpb.Source{
		Header: &dpb.SourceHeader{
			Api:  epb.SourceAPI_SOURCE_API_MOCK,
			Type: epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
			Id:   "foo",
		},
		Titles: []*dpb.Title{
			&dpb.Title{
				Title:        "Frieren",
				Localization: "en",
			},
		},
	}
	sources = []source.S{
		source.Make(frieren),
	}
)

func TestGet(t *testing.T) {
	t.Run("DNE", func(t *testing.T) {
		m := mock.New(sources)
		c := New(context.Background(), m, nil)

		got, err := c.Get(context.Background(), source.H{}, option.Remote(true))
		if err != nil {
			t.Errorf("Get() returned non-nil error: %v", err)
		}

		if got, want := len(m.GetHistory), 1; got != want {
			t.Errorf("len(m.GetHistory) = %d, want = %d", got, want)
		}

		want := source.S{}
		if diff := cmp.Diff(want, got, protocmp.Transform(), cmp.AllowUnexported(source.S{})); diff != "" {
			t.Errorf("Get() mismatch (-want +got):\n%v", diff)
		}
	})

	t.Run("Remote", func(t *testing.T) {
		m := mock.New(sources)
		c := New(context.Background(), m, nil)

		got, err := c.Get(context.Background(), source.Make(frieren).Header(), option.Remote(true))
		if err != nil {
			t.Errorf("Get() returned non-nil error: %v", err)
		}

		if got, want := len(m.GetHistory), 1; got != want {
			t.Errorf("len(m.GetHistory) = %d, want = %d", got, want)
		}

		want := source.Make(frieren)
		if diff := cmp.Diff(want, got, protocmp.Transform(), cmp.AllowUnexported(source.S{})); diff != "" {
			t.Errorf("Get() mismatch (-want +got):\n%v", diff)
		}
	})

	t.Run("Cache", func(t *testing.T) {
		m := mock.New(sources)
		c := New(context.Background(), m, []*dpb.Source{frieren})

		got, err := c.Get(context.Background(), source.Make(frieren).Header(), option.Remote(false))
		if err != nil {
			t.Errorf("Get() returned non-nil error: %v", err)
		}

		if got, want := len(m.GetHistory), 0; got != want {
			t.Errorf("len(m.GetHistory) = %d, want = %d", got, want)
		}

		want := source.Make(frieren)
		if diff := cmp.Diff(want, got, protocmp.Transform(), cmp.AllowUnexported(source.S{})); diff != "" {
			t.Errorf("Get() mismatch (-want +got):\n%v", diff)
		}
	})

	t.Run("Remote/WithNodeID", func(t *testing.T) {
		m := mock.New(sources)
		c := New(context.Background(), m, nil)

		got, err := c.Get(context.Background(), source.Make(frieren).Header(), true)
		if err != nil {
			t.Errorf("Get() returned non-nil error: %v", err)
		}

		if _, err := c.Put(context.Background(), got.WithNodeID("bar")); err != nil {
			t.Errorf("Put() returned non-nil error: %v", err)
		}

		got, err = c.Get(context.Background(), source.Make(frieren).Header(), true)

		if got, want := len(m.GetHistory), 2; got != want {
			t.Errorf("len(m.GetHistory) = %d, want = %d", got, want)
		}

		want := source.Make(frieren).WithNodeID("bar")
		if diff := cmp.Diff(want, got, protocmp.Transform(), cmp.AllowUnexported(source.S{})); diff != "" {
			t.Errorf("Get() mismatch (-want +got):\n%v", diff)
		}
	})

}

func TestSearch(t *testing.T) {
	t.Run("Local", func(t *testing.T) {
		m := mock.New(sources)
		c := New(context.Background(), m, []*dpb.Source{frieren})

		got, err := c.Search(context.Background(), "Frieren", option.Remote(false))
		if err != nil {
			t.Errorf("Search() returned non-nil error: %v", err)
		}

		if got, want := len(m.SearchHistory), 0; got != want {
			t.Errorf("len(m.SearchHistory) = %d, want = %d", got, want)
		}

		want := []source.S{source.Make(frieren)}
		if diff := cmp.Diff(want, got, protocmp.Transform(), cmp.AllowUnexported(source.S{})); diff != "" {
			t.Errorf("Search() mismatch (-want +got):\n%v", diff)
		}
	})

	t.Run("Remote", func(t *testing.T) {
		m := mock.New(sources)
		c := New(context.Background(), m, nil)

		got, err := c.Search(context.Background(), "Frieren", option.Remote(true))
		if err != nil {
			t.Errorf("Search() returned non-nil error: %v", err)
		}

		if _, err := c.Put(context.Background(), got[0].WithNodeID("bar")); err != nil {
			t.Errorf("Put() returned non-nil error: %v", err)
		}

		got, err = c.Search(context.Background(), "Frieren", option.Remote(true))

		if got, want := len(m.SearchHistory), 2; got != want {
			t.Errorf("len(m.SearchHistory) = %d, want = %d", got, want)
		}

		want := []source.S{source.Make(frieren).WithNodeID("bar")}
		if diff := cmp.Diff(want, got, protocmp.Transform(), cmp.AllowUnexported(source.S{})); diff != "" {
			t.Errorf("Search() mismatch (-want +got):\n%v", diff)
		}
	})
}
