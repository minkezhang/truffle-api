package db

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/minkezhang/truffle-api/client"
	"github.com/minkezhang/truffle-api/db/atom"
	"github.com/minkezhang/truffle-api/db/atom/metadata/empty"
	"github.com/minkezhang/truffle-api/db/internal/client/mock"
	"github.com/minkezhang/truffle-api/db/node"
	"github.com/minkezhang/truffle-api/db/query"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

func TestGet(t *testing.T) {
	db, err := New(context.Background(), O{})
	if err != nil {
		t.Errorf("New() returned unexpected error: %v", err)
	}

	n := node.New(node.O{
		ID:              "foo",
		IsAuthoritative: false,
		AtomType:        epb.Type_TYPE_TV,
		Atoms: []*atom.A{
			atom.New(atom.O{
				APIType: epb.API_API_TRUFFLE,
				APIID:   "bar",
				Titles: []atom.T{
					{Title: "Firefly"},
				},
				PreviewURL: "",
				Score:      91,
				AtomType:   epb.Type_TYPE_TV,
				Metadata:   empty.M{},
			}),
			atom.New(atom.O{
				APIType: epb.API_API_TRUFFLE,
				APIID:   "baz",
				Titles: []atom.T{
					{Title: "Firefly: The Complete Series"},
				},
				PreviewURL: "",
				Score:      92,
				AtomType:   epb.Type_TYPE_TV,
				Metadata:   empty.M{},
			}),
		},
	})
	want := n.Copy()
	want.SetIsAuthoritative(true)

	if _, err := db.Add(context.Background(), n); err != nil {
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
		AtomType:        epb.Type_TYPE_TV,
		Atoms: []*atom.A{
			atom.New(atom.O{
				APIType: epb.API_API_TRUFFLE,
				APIID:   "bar",
				Titles: []atom.T{
					{Title: "Firefly"},
				},
				PreviewURL: "",
				Score:      91,
				AtomType:   epb.Type_TYPE_TV,
				Metadata:   empty.M{},
			}),
			atom.New(atom.O{
				APIType: epb.API_API_TRUFFLE,
				APIID:   "baz",
				Titles: []atom.T{
					{Title: "Firefly: The Complete Series"},
				},
				PreviewURL: "",
				Score:      92,
				AtomType:   epb.Type_TYPE_TV,
				Metadata:   empty.M{},
			}),
		},
	})

	if _, err := db.Add(context.Background(), n); err != nil {
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
			AtomType:        epb.Type_TYPE_TV,
			Atoms: []*atom.A{
				atom.New(atom.O{
					APIType: epb.API_API_TRUFFLE,
					APIID:   "bar",
					Titles: []atom.T{
						{Title: "Firefly"},
					},
					PreviewURL: "",
					Score:      91,
					AtomType:   epb.Type_TYPE_TV,
					Metadata:   empty.M{},
				}),
				atom.New(atom.O{
					APIType: epb.API_API_TRUFFLE,
					APIID:   "baz",
					Titles: []atom.T{
						{Title: "Firefly: The Complete Series"},
					},
					PreviewURL: "",
					Score:      92,
					AtomType:   epb.Type_TYPE_TV,
					Metadata:   empty.M{},
				}),
			},
		})
		m := n.Copy()
		m.SetIsAuthoritative(true)
		want := []*node.N{m.Copy()}

		db.Add(context.Background(), n)

		got, err := db.Query(context.Background(), query.New(query.O{
			APIs:      []epb.API{epb.API_API_TRUFFLE},
			AtomTypes: []epb.Type{m.AtomType()},
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
			t.Errorf("Query() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("External", func(t *testing.T) {
		n := node.New(node.O{
			AtomType: epb.Type_TYPE_TV,
			Atoms: []*atom.A{
				atom.New(atom.O{
					APIType: epb.API_API_TRUFFLE,
					APIID:   "bar",
					Titles: []atom.T{
						{Title: "Firefly"},
					},
					PreviewURL: "",
					Score:      91,
					AtomType:   epb.Type_TYPE_TV,
					Metadata:   empty.M{},
				}),
			},
		})
		c := mock.New(mock.O{
			Data: n.Atoms(),
		})
		db, err := New(context.Background(), O{Clients: []client.C{c}})
		if err != nil {
			t.Errorf("New() returned unexpected error: %v", err)
		}

		want := []*node.N{n.Copy()}
		got, err := db.Query(context.Background(), query.New(query.O{
			APIs:      []epb.API{epb.API_API_MAL},
			AtomTypes: []epb.Type{epb.Type_TYPE_TV},
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
			t.Errorf("Query() mismatch (-want +got):\n%s", diff)
		}
	})
}
