package client

import (
	"context"

	"github.com/minkezhang/bene-api/client/query"
	"github.com/minkezhang/bene-api/db/atom"
	"github.com/minkezhang/bene-api/db/enums"
)

type C interface {
	APIType() enums.ClientAPI
	Get(ctx context.Context, id string) (*atom.A, error)
	Query(ctx context.Context, q query.Q) ([]*atom.A, error)
}
