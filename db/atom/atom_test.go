package atom

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/minkezhang/bene-api/db/atom/internal/metadata/mock"

	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

func TestMerge(t *testing.T) {
	got := Merge(
		New(O{
			APIType: epb.API_API_BENE,
			APIID:   "foo",
			Titles: []T{
				{Title: "Firefly"},
			},
			PreviewURL: "",
			Score:      91,
			AtomType:   epb.Type_TYPE_TV,
			Metadata: mock.New(mock.O{
				Producers: []string{"Joss Whedon"},
			}),
		}),
		New(O{
			APIType: epb.API_API_BENE,
			APIID:   "foo",
			Titles: []T{
				{Title: "Firefly"},
			},
			PreviewURL: "overwrite",
			Score:      92,
			AtomType:   epb.Type_TYPE_TV,
			Metadata: mock.New(mock.O{
				Producers: []string{"Tim Minear"},
			}),
		}))

	want := New(O{
		APIType: epb.API_API_BENE,
		APIID:   "foo",
		Titles: []T{
			{Title: "Firefly"}, // Remove duplicates
		},
		PreviewURL: "overwrite",
		Score:      92,
		AtomType:   epb.Type_TYPE_TV,
		Metadata: mock.New(mock.O{
			Producers: []string{"Joss Whedon", "Tim Minear"},
		}),
	})

	if diff := cmp.Diff(
		want,
		got,
		cmp.AllowUnexported(A{}, mock.M{}),
		cmpopts.SortSlices(func(a, b string) bool { return a > b }),
	); diff != "" {
		t.Errorf("Merge() mismatch (-want +got):\n%s", diff)
	}
}
