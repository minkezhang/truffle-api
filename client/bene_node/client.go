package bene_node

import (
	"context"

	dpb "github.com/minkezhang/bene-api/proto/go/data"
)

type C struct{}

func (c *C) Query(ctx context.Context, q *dpb.Query) (*dpb.Node, error) {
	return nil, nil
}
