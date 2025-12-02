package db

import (
	"context"

	"github.com/minkezhang/bene-api/db/atom"
	"github.com/minkezhang/bene-api/db/atom/tv"
	"github.com/minkezhang/bene-api/db/client"
	"github.com/minkezhang/bene-api/db/node"
)

type es[T atom.A[T]] struct {
	nodes  map[string]*node.N[T]
	client client.C[T]
}

type DB struct {
	nodeTV es[*tv.T]
}

func Create[T atom.A[T]](db *DB /* data */) (*node.N[T], error) { return nil, nil /* unimplemented */ }
func Get[T atom.A[T]](db *DB /* query */) (*node.N[T], error)   { return nil, nil /* unimplemented */ }
func Update[T atom.A[T]](db *DB /* query */) error              { return nil /* unimplemented */ }
func Delete[T atom.A[T]](db *DB, n *node.N[T]) error            { return nil /* unimplemented */ }

func Query[T atom.A[T]](ctx context.Context, db *DB, q client.Q[T]) ([]*node.N[T], error) {
	atoms, err := db.nodeTV.client.Query(ctx, (client.Q[*tv.T]).q)
	if err != nil {
		return nil, err
	}
	res := []*node.N[T]{}
	for _, a := range atoms {
		res = append(res, node.New(node.O[T]{
			Type:  a.Type(),
			ID:    a.ID(),
			Atoms: []T{a},
		}))
	}
	return res, nil
}
