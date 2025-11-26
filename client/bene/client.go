package bene

import (
	"context"

	dpb "github.com/minkezhang/bene-api/proto/go/data"
)

type C struct{}

func (c *C) Query(ctx context.Context, q *dpb.QueryAtom) ([]*dpb.Atom, error) {
	return nil, nil
}
