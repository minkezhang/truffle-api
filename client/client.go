package client

import (
	"context"

	dpb "github.com/minkezhang/bene-api/proto/go/data"
)

type E interface {
	error

	Status() int
}

// C is an atomic client.
type C interface {
	Query(ctx context.Context, q *dpb.QueryAtom) ([]*dpb.Atom, error)
}
