package db

import (
	"fmt"
	"testing"

	apb "github.com/minkezhang/bene-api/proto/go/api"
	dpb "github.com/minkezhang/bene-api/proto/go/data"
)

var (
	mock []*dpb.Node = []*dpb.Node{
		{
			Id:   "foo",
			Type: dpb.NodeType_NODE_TV,
			Atoms: []*dpb.Atom{
				{
					Api: apb.API_API_BENE,
					Id:  "foo-atom",
					Titles: []*dpb.Title{
						&dpb.Title{
							Title:        "Firefly",
							Localization: "us-en",
						},
					},
				},
			},
		},
	}
)

func TestNewDB(t *testing.T) {
	db, err := NewDB(O{Data: mock})
	if err != nil {
		t.Errorf("New() returned unexpected error: %v", err)
	}
	fmt.Println(db.Get("foo").Type())
}
