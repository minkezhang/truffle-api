package bene

import (
	"context"
	"testing"

	"github.com/minkezhang/bene-api/db/atom"
	"github.com/minkezhang/bene-api/db/atom/tv"
)

func TestGet(t *testing.T) {
	want := tv.DebugNewOrDie(
		tv.O{
			O: atom.O{
				Titles: []atom.T{
					{
						Title: "Firefly",
					},
				},
				ID: "foo",
			},
		},
	)
	c := DebugNewOrDie(O[*tv.T]{
		Cache: []*tv.T{
			want,
			tv.DebugNewOrDie(
				tv.O{
					O: atom.O{
						Titles: []atom.T{
							{
								Title: "Buffy the Vampire Slayer",
							},
						},
						ID: "bar",
					},
				},
			),
		},
	})
	a, err := c.Get(context.Background(), "foo")
	if err != nil {
		t.Errorf("Query() returned unexpected error: %v", err)
	}
	if got := a.GetBase().Titles[0].Title; got != want.GetBase().Titles[0].Title {
		t.Errorf("Query() = %v, want = %v", got, want)
	}
}
