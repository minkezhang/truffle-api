package client

import (
	"context"

	"github.com/minkezhang/bene-api/db/atom"
	"github.com/minkezhang/bene-api/db/enums"
)

type Q struct {
	Title string
}

type C interface {
	APIType() enums.ClientAPI
	Get(ctx context.Context, id string) (*atom.A, error)
	Query(ctx context.Context, q Q) ([]*atom.A, error)
}
