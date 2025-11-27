package client

import (
	"context"
	"github.com/minkezhang/bene-api/client/errors"

	dpb "github.com/minkezhang/bene-api/proto/go/data"
)

// C is an atomic client.
type C interface {
	Query(ctx context.Context, q *dpb.QueryAtom) ([]*dpb.Atom, errors.E)
}
