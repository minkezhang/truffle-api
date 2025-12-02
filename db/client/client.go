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
}

type C[T atom.A[T]] interface {
	API() atom.ClientAPI

	// IsSupported returns if the given atom type (e.g. TV, Movie, etc.) is
	// supported for this client.
	IsSupported(t atom.AtomType) bool

	// Get returns a single (potentially nil) atom with the given
	// API-specific ID.
	Get(ctx context.Context, id string) (T, error)

	// Query returns all (potentially none) atoms which match the input
	// query.
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

func New(o O) (*Base, error) {
	c := &Base{
		api:       o.API,
		supported: map[atom.AtomType]interface{}{},
	}
	for _, t := range o.SupportedTypes {
		c.supported[t] = struct{}{}
	}
	return c, nil
}

func (c *Base) API() atom.ClientAPI { return c.api }

func (c *Base) IsSupported(t atom.AtomType) bool {
	_, ok := c.supported[t]
	return ok
}
