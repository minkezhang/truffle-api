package mal

import (
	"context"
	"fmt"
	"testing"

	"github.com/minkezhang/bene-api/client/query"
	"github.com/minkezhang/bene-api/db/atom"
	"google.golang.org/protobuf/encoding/prototext"

	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

func TestGet(t *testing.T) {
	c := New(O{
		ClientID: "6114d00ca681b7701d1e15fe11a4987e",
	})
	r, _ := c.Get(context.Background(), query.G{
		AtomType: epb.Type_TYPE_TV,
		ID:       "52807",
	})

	s, err := prototext.MarshalOptions{
		Multiline: true,
		Indent:    " ",
	}.Marshal(atom.Save(r))
	fmt.Println(string(s), err)
}
