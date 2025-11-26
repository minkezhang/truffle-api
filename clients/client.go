package client

import (
	"context"

	dpb "github.com/minkezhang/bene-api/proto/go/data"
)

type E interface {
	error

	Status() int
}

type C interface {
	Query(ctx context.Context, q *dpb.Query) (*dpb.Node, E)
}
