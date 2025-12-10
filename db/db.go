// Package db encapsulates logic necessary to query for media data across
// multiple clients.
package db

import (
	"context"

	"github.com/minkezhang/bene-api/client"
	"github.com/minkezhang/bene-api/db/atom"
	"github.com/minkezhang/bene-api/db/generator"
	"github.com/minkezhang/bene-api/db/node"
	"github.com/minkezhang/bene-api/db/query"

	cq "github.com/minkezhang/bene-api/client/query"
	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

// g is an ID generator
type g interface {
	Generate() string
}

type O struct {
	Data    []*node.N
	Clients []client.C
}

type DB struct {
	data    map[epb.Type]map[string]*node.N
	clients map[epb.API]client.C

	g g // ID generator
}

func New(ctx context.Context, o O) (*DB, error) {
	ids := []string{}
	db := &DB{
		data:    map[epb.Type]map[string]*node.N{},
		clients: map[epb.API]client.C{},
	}
	for _, n := range o.Data {
		db.Add(ctx, n)
		ids = append(ids, n.ID())
	}
	db.g = generator.New(generator.O{
		IDs: ids,
	})

	for _, c := range o.Clients {
		db.clients[c.APIType()] = c
	}

	return db, nil
}

// Add will add a given node to the DB.
//
// A Node ID will be generated if no Node ID is provided and returned.
func (db *DB) Add(ctx context.Context, n *node.N) (string, error) {
	n.SetIsAuthoritative(true)
	if n.ID() == "" {
		n = node.New(node.O{
			ID:              db.g.Generate(),
			AtomType:        n.AtomType(),
			IsAuthoritative: n.IsAuthoritative(),
			IsQueued:        n.IsQueued(),
			Atoms:           n.Atoms(),
		})
	}
	if _, ok := db.data[n.AtomType()]; !ok {
		db.data[n.AtomType()] = map[string]*node.N{}
	}
	db.data[n.AtomType()][n.ID()] = n
	return n.ID(), nil
}

func (db *DB) Remove(ctx context.Context, id string) error {
	n, err := db.Get(ctx, id)
	if err != nil {
		return err
	}
	if n == nil {
		return nil
	}

	delete(db.data[n.AtomType()], n.ID())
	return nil
}

// Get returns a specific node by the Node ID.
func (db *DB) Get(ctx context.Context, id string) (*node.N, error) {
	for _, ns := range db.data {
		if n, ok := ns[id]; ok {
			return n, nil
		}
	}
	return nil, nil
}

// Query returns a set of nodes which match the input query.
//
// If q.APIs is set, all matching clients will be queried.
func (db *DB) Query(ctx context.Context, q *query.Q) ([]*node.N, error) {
	res := []*node.N{}
	if q.IsSupportedAPI(epb.API_API_BENE) {
		for atomType := range db.data {
			for _, n := range db.data[atomType] {
				match, err := cq.RegExp(q.Q, n.Virtual())
				if err != nil {
					return nil, err
				}
				if match > 0 {
					res = append(res, n.Copy())
				}
			}
		}
	}

	for _, c := range db.clients {
		if q.IsSupportedAPI(c.APIType()) {
			atoms, err := c.Query(ctx, q.Q)
			if err != nil {
				return nil, err
			}
			for _, a := range atoms {
				res = append(res, node.New(node.O{
					IsAuthoritative: false,
					AtomType:        a.AtomType(),
					Atoms:           []*atom.A{a},
				}))
			}
		}
	}

	return res, nil
}
