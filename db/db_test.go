package db

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/minkezhang/bene-api/db/atom"
	"github.com/minkezhang/bene-api/db/atom/empty"
	"github.com/minkezhang/bene-api/db/enums"
	"github.com/minkezhang/bene-api/db/node"
	"github.com/minkezhang/bene-api/db/query"
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

	if err := db.Add(context.Background(), n); err != nil {
		t.Errorf("Add() raised unexpected error: %v", err)
	}
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

	if err := db.Add(context.Background(), n); err != nil {
		t.Errorf("Add() raised unexpected error: %v", err)
	}
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

func TestQuery(t *testing.T) {
	t.Run("NoExternal", func(t *testing.T) {
		db, _ := New(context.Background(), O{})
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
		m := n.Copy()
		m.SetIsAuthoritative(true)
		want := []*node.N{m.Copy()}

		db.Add(context.Background(), n)

		got, err := db.Query(context.Background(), query.New(query.O{
			APIs:      []enums.ClientAPI{enums.ClientAPIBene},
			AtomTypes: []enums.AtomType{m.AtomType()},
			Title:     "fly",
		}))
		if err != nil {
			t.Errorf("Query() returned unexpected error: %v", err)
		}

		if diff := cmp.Diff(
			want,
			got,
			cmp.AllowUnexported(node.N{}, atom.A{}),
		); diff != "" {
			t.Errorf("Get() mismatch (-want +got):\n%s", diff)
		}
	})
}
