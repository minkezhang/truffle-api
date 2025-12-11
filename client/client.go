// Package client defines the interface for different data sources to interact
// with the Truffle API.
package client

import (
	"context"

	"github.com/minkezhang/truffle-api/client/query"
	"github.com/minkezhang/truffle-api/db/atom"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

type C interface {
	APIType() epb.API

	// Get returns a single atom given an ID associated with the API.
	Get(ctx context.Context, g query.G) (*atom.A, error)

	// Query returns a (potentially empty) list of atoms with the given
	// input.
	Query(ctx context.Context, q *query.Q) ([]*atom.A, error)
}
