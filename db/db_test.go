package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/minkezhang/bene-api/db/atom"
	"github.com/minkezhang/bene-api/db/atom/empty"
	"github.com/minkezhang/bene-api/db/enums"
	"github.com/minkezhang/bene-api/db/node"
)

func TestAdd(t *testing.T) {
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
	db.AddNode(context.Background(), n)
	fmt.Println(db.data[enums.AtomTypeTV]["foo"])
	n, err = db.Get(context.Background(), "foo")
	fmt.Println(n.Virtual())
}
