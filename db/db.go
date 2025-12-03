package db

import (
	"context"
	"regexp"

	"github.com/minkezhang/bene-api/client"
	"github.com/minkezhang/bene-api/db/enums"
	"github.com/minkezhang/bene-api/db/node"
	"github.com/minkezhang/bene-api/db/query"
)

type O struct {
	Data []*node.N
}

type DB struct {
	data    map[enums.AtomType]map[string]*node.N
	clients map[enums.ClientAPI]client.C
}

// AddNode will add a given node to the DB.
//
// A Node ID will be generated if no Node ID is provided.
func (db *DB) AddNode(n *node.N) {
	if n.ID() != "" {
		// ...
	}
	if _, ok := db.data[n.AtomType()]; !ok {
		db.data[n.AtomType()] = map[string]*node.N{}
	}
	db.data[n.AtomType()][n.ID()] = n
}

func (db *DB) RemoveNode(ctx context.Context, id string) error {
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
	pattern, err := regexp.Compile(q.Title())
	if err != nil {
		return nil, err
	}

	res := []*node.N{}
	if q.IsSupportedAPI(enums.ClientAPIBene) {
		for atomType := range db.data {
			if q.IsSupportedAtomType(atomType) {
				for _, n := range db.data[atomType] {
					a := n.Virtual()
					for _, t := range a.Titles() {
						if pattern.MatchString(t.Title) {
							res = append(res, n.Copy())
						}
					}
				}
			}
		}
	}
	// TODO(minkezhang): Implement client query logic
	return res, nil
}
