package client

import (
	"context"

	"github.com/minkezhang/bene-api/client/query"
	"github.com/minkezhang/bene-api/db/atom"

	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

type C interface {
	APIType() epb.API

	// Get returns a single atom given an ID associated with the API.
	Get(ctx context.Context, g query.G) (*atom.A, error)

	// Query returns a (potentially empty) list of atoms with the given
	// input.
	Query(ctx context.Context, q *query.Q) ([]*atom.A, error)
}
