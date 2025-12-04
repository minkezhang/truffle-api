package db

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/minkezhang/bene-api/db/atom"
	"github.com/minkezhang/bene-api/db/atom/empty"
	"github.com/minkezhang/bene-api/db/enums"
	"github.com/minkezhang/bene-api/db/node"
)

func TestGet(t *testing.T) {
	db, err := New(context.Background(), O{})
	if err != nil {
		t.Errorf("New() returned unexpected error: %v", err)
	}

	n := node.New(node.O{
		ID:              "foo",
		IsAuthoritative: false,
		AtomType:        enums.AtomTypeTV,
		Atoms: []*atom.A{
			atom.New(atom.O{
				APIType: enums.ClientAPIBene,
				APIID:   "bar",
				Titles: []atom.T{
					{Title: "Firefly"},
				},
				PreviewURL: "",
				Score:      91,
				AtomType:   enums.AtomTypeTV,
				Aux:        empty.A{},
			}),
			atom.New(atom.O{
				APIType: enums.ClientAPIBene,
				APIID:   "baz",
				Titles: []atom.T{
					{Title: "Firefly: The Complete Series"},
				},
				PreviewURL: "",
				Score:      92,
				AtomType:   enums.AtomTypeTV,
				Aux:        empty.A{},
			}),
		},
	})
	want := n.Copy()
	want.SetIsAuthoritative(true)

	db.Add(context.Background(), n)
	got, err := db.Get(context.Background(), want.ID())
	if err != nil {
		t.Errorf("Get() raised unexpected error: %v", err)
	}
	if diff := cmp.Diff(
		want,
		got,
		cmp.AllowUnexported(node.N{}, atom.A{}),
	); diff != "" {
		t.Errorf("Get() mismatch (-want +got):\n%s", diff)
	}
}

func TestRemove(t *testing.T) {
	db, err := New(context.Background(), O{})
	if err != nil {
		t.Errorf("New() returned unexpected error: %v", err)
	}

	n := node.New(node.O{
		ID:              "foo",
		IsAuthoritative: false,
		AtomType:        enums.AtomTypeTV,
		Atoms: []*atom.A{
			atom.New(atom.O{
				APIType: enums.ClientAPIBene,
				APIID:   "bar",
				Titles: []atom.T{
					{Title: "Firefly"},
				},
				PreviewURL: "",
				Score:      91,
				AtomType:   enums.AtomTypeTV,
				Aux:        empty.A{},
			}),
			atom.New(atom.O{
				APIType: enums.ClientAPIBene,
				APIID:   "baz",
				Titles: []atom.T{
					{Title: "Firefly: The Complete Series"},
				},
				PreviewURL: "",
				Score:      92,
				AtomType:   enums.AtomTypeTV,
				Aux:        empty.A{},
			}),
		},
	})

	db.Add(context.Background(), n)
	if err := db.Remove(context.Background(), n.ID()); err != nil {
		t.Errorf("Remove() returned unexpected error: %v", err)
	}
	got, err := db.Get(context.Background(), n.ID())
	if err != nil {
		t.Errorf("Get() raised unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("Get() = %v, want = nil", got)
	}
}
