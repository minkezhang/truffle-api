package db

import (
	apb "github.com/minkezhang/bene-api/proto/go/api"
	dpb "github.com/minkezhang/bene-api/proto/go/data"
)

type O struct {
	Path string
	Data []*dpb.Node
}

type DB struct {
	nodes map[string]*Node
	atoms map[apb.API]map[string]*Atom

	path string
}

func NewDB(o O) (*DB, error) {
	db := &DB{
		nodes: map[string]*Node{},
		atoms: map[apb.API]map[string]*Atom{},
		path:  o.Path,
	}

	for _, n := range o.Data {
		db.nodes[n.GetId()] = &Node{
			db: db,
			pb: n,
		}
		for _, a := range n.GetAtoms() {
			if _, ok := db.atoms[a.GetApi()]; !ok {
				db.atoms[a.GetApi()] = map[string]*Atom{}
			}
			db.atoms[a.GetApi()][a.GetId()] = &Atom{
				db:     db,
				nodeID: n.GetId(),
				atomID: a.GetId(),
				pb:     a,
			}
		}
	}

	return db, nil
}

func (db *DB) Get(id string) *Node {
	if n, ok := db.nodes[id]; ok {
		return n
	}
	return nil
}

func (db *DB) GetAtoms(id string) *b

func (db *DB) Write() error {
	return nil
}
