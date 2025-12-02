package node

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/minkezhang/bene-api/db/atom"
	"github.com/minkezhang/bene-api/db/atom/tv"
)

func TestUnion(t *testing.T) {
	want := tv.DebugNewOrDie(tv.O{
		O: atom.O{
			API: atom.ClientAPIVirtual,
			ID:  "",
			Titles: []atom.T{
				{
					Title: "Firefly",
				},
				{
					Title: "Firefly: The Complete Series",
				},
			},
			Score: 91,
		},
		Genres: []string{"space", "western"},
	})
	n := DebugNewOrDie(O[*tv.T]{
		Type: atom.AtomTypeTV,
		ID:   "foo",
		Atoms: []*tv.T{
			tv.DebugNewOrDie(tv.O{
				O: atom.O{
					API: atom.ClientAPIBene,
					ID:  "bar",
					Titles: []atom.T{
						{
							Title: "Firefly",
						},
					},
					Score: 92,
				},
				Genres: []string{"space"},
			}),
			tv.DebugNewOrDie(tv.O{
				O: atom.O{
					API: atom.ClientAPIBene,
					ID:  "baz",
					Titles: []atom.T{
						{
							Title: "Firefly: The Complete Series",
						},
					},
					Score: 91,
				},
				Genres: []string{"western"},
			}),
		},
	})
	got, err := n.Union()
	if err != nil {
		t.Errorf("Union() raised unexpected error: %v", err)
	}
	if diff := cmp.Diff(
		want,
		got,
		cmp.AllowUnexported(atom.Base{}),
		cmpopts.EquateEmpty(),
		cmpopts.SortSlices(func(a, b string) bool { return a < b }),
	); diff != "" {
		t.Errorf("Union() mismatch (-want +got):\n%s", diff)
	}
}
