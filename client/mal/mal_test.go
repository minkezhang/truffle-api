package mal

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/minkezhang/bene-api/client/query"
	"github.com/minkezhang/bene-api/db/atom"
	"github.com/minkezhang/bene-api/db/atom/metadata/shared/video"
	"github.com/minkezhang/bene-api/db/atom/metadata/tv"

	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

const (
	// This is a publically-known API key for the MAL Android app.
	MALClientID = "6114d00ca681b7701d1e15fe11a4987e"
)

func TestGet(t *testing.T) {
	c := New(O{
		ClientID: MALClientID,
	})

	got, err := c.Get(context.Background(), query.G{
		AtomType: epb.Type_TYPE_TV,
		ID:       "235", // Detective Conan
	})
	if err != nil {
		t.Errorf("Get() returned unexpected error: %v", err)
	}

	want := atom.New(atom.O{
		APIType: epb.API_API_MAL,
		APIID:   "235",
		Titles: []atom.T{
			atom.T{Title: "Meitantei Conan"},
			atom.T{Title: "Case Closed", Localization: "en"},
		},
		PreviewURL: "https://cdn.myanimelist.net/images/anime/7/75199l.jpg",
		Score:      81,
		AtomType:   epb.Type_TYPE_TV,
		Metadata: tv.New(tv.O{
			IsAnimated: true,
			IsAnime:    true,
			Genres:     []string{"Adventure", "Comedy", "Detective", "Mystery", "Shounen"},
			Studios:    []string{"TMS Entertainment"},
		}),
	})
	if diff := cmp.Diff(
		want,
		got,
		cmp.AllowUnexported(
			atom.A{},
			tv.M{},
			video.M{},
		),
	); diff != "" {
		t.Errorf("Get() mismatch (-want +got):\n%s", diff)
	}
}
