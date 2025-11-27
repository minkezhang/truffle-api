package bene

import (
	"context"

	"github.com/minkezhang/bene-api/client/errors"

	dpb "github.com/minkezhang/bene-api/proto/go/data"
)

type C struct{}

func (c *C) Query(ctx context.Context, q *dpb.QueryAtom) ([]*dpb.Atom, errors.E) {
	return nil, errors.Unimplemented{}
}
