// Package client is a collection of API hooks to distinct metadata aggregation
// sources.
//
// Bene will query these sources for data.
package client

import (
	"context"

	"github.com/minkezhang/bene-api/db/atom"
)

type Q[T atom.A[T]] struct {
	Title string
	ID    string
}

type C[T atom.A[T]] interface {
	API() atom.ClientAPI
	IsSupported(t atom.AtomType) bool
	Query(ctx context.Context, q Q[T]) ([]T, error)
}

type Base struct {
	api       atom.ClientAPI
	supported map[atom.AtomType]interface{}
}

type O struct {
	API            atom.ClientAPI
	SupportedTypes []atom.AtomType
}

func New(o O) *Base {
	c := &Base{
		api:       o.API,
		supported: map[atom.AtomType]interface{}{},
	}
	for _, t := range o.SupportedTypes {
		c.supported[t] = struct{}{}
	}
	return c
}

func (c *Base) API() atom.ClientAPI { return c.api }

func (c *Base) IsSupported(t atom.AtomType) bool {
	_, ok := c.supported[t]
	return ok
}
