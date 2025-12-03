package db

import (
	"context"

	"github.com/minkezhang/bene-api/client"
	"github.com/minkezhang/bene-api/client/bene"
	"github.com/minkezhang/bene-api/db/enums"
	"github.com/minkezhang/bene-api/db/node"
)

type O struct {
	Data []*node.N
}

type Q struct {
	APIs      []enums.ClientAPI
	AtomTypes []enums.AtomType
	Title     string
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
	for _, a := range n.Atoms() {
		// Cache atoms of all APIs
		db.clients[enums.ClientAPIBene].(*bene.C).Add(a)
	}
}

func (db *DB) RemoveNode(ctx context.Context, id string) error {
	n, err := db.Get(ctx, id)
	if err != nil {
		return err
	}
	if n == nil {
		return nil
	}

	for _, a := range n.Atoms() {
		db.clients[enums.ClientAPIBene].(*bene.C).Remove(a.APIID())
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
func (db *DB) Query(ctx context.Context, q Q) ([]*node.N, error) {
	res := []*node.N{}

	types := map[enums.AtomType]interface{}{}
	apis := map[enums.ClientAPI]interface{}{}
	for _, t := range q.AtomTypes {
		types[t] = struct{}{}
	}
	for _, api := range q.APIs() {
		apis[t] = struct{}{}
	}

	if _, ok := apis[enums.ClientAPIBene] {

	}

	return nil, nil /* unimplemented */
}
