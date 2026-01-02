// Package db implements a Truffle database.
package db

import (
	"cmp"
	"context"
	"fmt"
	"slices"

	"github.com/minkezhang/truffle-api/client"
	"github.com/minkezhang/truffle-api/client/mal"
	"github.com/minkezhang/truffle-api/client/option"
	"github.com/minkezhang/truffle-api/client/truffle"
	"github.com/minkezhang/truffle-api/data/node"
	"github.com/minkezhang/truffle-api/data/source"
	"github.com/minkezhang/truffle-api/data/source/util"
	"github.com/minkezhang/truffle-api/util/generator"
	"github.com/minkezhang/truffle-api/util/match"

	cpb "github.com/minkezhang/truffle-api/proto/go/config"
	dpb "github.com/minkezhang/truffle-api/proto/go/data"
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

func New(ctx context.Context, config *cpb.Config, db *dpb.Database) *DB {
	mal_client := client.New(
		ctx,
		mal.Make(config.GetMal()),
		util.Filter(db.GetSources(), func(v *dpb.Source) bool {
			return v.GetHeader().GetApi() != epb.SourceAPI_SOURCE_API_MAL
		}),
	)

	nodes := map[string]node.N{}
	for _, n := range db.GetNodes() {
		nodes[n.GetHeader().GetId()] = node.Make(n)
	}

	return &DB{
		clients: map[epb.SourceAPI]*client.Cache{
			epb.SourceAPI_SOURCE_API_TRUFFLE: client.New(
				ctx,
				truffle.New(
					config.GetTruffle(),
					util.Filter(db.GetSources(), func(v *dpb.Source) bool {
						return v.GetHeader().GetApi() != epb.SourceAPI_SOURCE_API_TRUFFLE
					}),
				),
				util.Filter(db.GetSources(), func(v *dpb.Source) bool {
					return v.GetHeader().GetApi() != epb.SourceAPI_SOURCE_API_TRUFFLE
				}),
			),
			epb.SourceAPI_SOURCE_API_MAL:               mal_client,
			epb.SourceAPI_SOURCE_API_MAL_ANIME_PARTIAL: mal_client,
			epb.SourceAPI_SOURCE_API_MAL_MANGA_PARTIAL: mal_client,
		},
		nodes: nodes,
		generator: generator.New(generator.O{
			IDs: util.Apply(db.GetNodes(), func(v *dpb.Node) string { return v.GetHeader().GetId() }),
		}),
	}
}

type DB struct {
	clients   map[epb.SourceAPI]*client.Cache
	nodes     map[string]node.N
	generator *generator.G
}

func (db *DB) DeleteNode(ctx context.Context, header node.H) error {
	n, err := db.GetNode(ctx, header, option.Remote(false))
	if err != nil {
		return err
	}

	for _, s := range n.Sources() {
		if err := db.clients[s.Header().API()].Delete(ctx, s.Header()); err != nil {
			return err
		}
	}
	delete(db.nodes, n.Header().ID())
	return nil
}

func (db *DB) GetNode(ctx context.Context, header node.H, refresh option.Remote) (node.N, error) {
	n, ok := db.nodes[header.ID()]
	if !ok {
		return node.N{}, nil
	}

	sources := []source.S{}
	for _, api := range []epb.SourceAPI{
		epb.SourceAPI_SOURCE_API_TRUFFLE,
		epb.SourceAPI_SOURCE_API_MAL,
	} {
		_sources, err := db.clients[api].SearchByNodeID(ctx, header, refresh)
		if err != nil {
			return node.N{}, err
		}
		sources = append(sources, _sources...)
	}

	return n.WithSources(sources), nil
}

// Put saves a source into the Truffle DB.
//
//  1. If the source does not exist, save it.
//  1. If the source does not specify a node ID, generate one.
//  1. If the source specifies a node ID different from its current designation,
//     attempt to link the new ID.
func (db *DB) Put(ctx context.Context, s source.S) (source.H, error) {
	var t source.S // Pre-existing source if source exists
	var n node.N   // Associated with s

	fmt.Println(s.NodeID())
	// Create node if the source is not yet linked.
	if _, ok := db.nodes[s.NodeID()]; !ok || s.NodeID() == "" {
		s = s.WithNodeID(db.generator.Generate())
		n = node.Make(&dpb.Node{
			Header: &dpb.NodeHeader{
				Id:   s.NodeID(),
				Type: s.Header().Type(),
			},
		})
		db.nodes[n.Header().ID()] = n
		fmt.Println(n.Header().Type().String())
	}

	// Get old source if it exists
	t, err := db.Get(ctx, s.Header(), option.Remote(false))
	if err != nil {
		return source.H{}, err
	}

	header, err := db.clients[s.Header().API()].Put(ctx, s)
	if err != nil {
		return header, err
	}
	s = s.WithHeader(header)

	// The source is being relinked to a new node.
	if t != (source.S{}) && s.NodeID() != t.NodeID() {
		// Check if we can actually link to this node.
		n, err := db.GetNode(ctx, node.Make(&dpb.Node{
			Header: &dpb.NodeHeader{
				Id:   s.NodeID(),
				Type: s.Header().Type(),
			},
		}).Header(), option.Remote(false))
		if err != nil {
			return source.H{}, err
		}
		if n.Header().Type() != s.Header().Type() {
			return source.H{}, fmt.Errorf("mismatching node and source media types: %v != %v", n.Header().Type(), s.Header().Type())
		}

		// Check if we need to delete old node if no remaining links
		// exist.
		m, err := db.GetNode(ctx, node.Make(&dpb.Node{
			Header: &dpb.NodeHeader{
				Id:   t.NodeID(),
				Type: t.Header().Type(),
			},
		}).Header(), option.Remote(false))
		if err != nil {
			return source.H{}, err
		}
		if len(m.Sources()) == 0 {
			if err := db.DeleteNode(ctx, m.Header()); err != nil {
				return source.H{}, err
			}
		}

	}

	return s.Header(), nil
}

func (db *DB) Delete(ctx context.Context, header source.H) error {
	s, err := db.clients[header.API()].Get(ctx, header, option.Remote(false))
	if err != nil {
		return err
	}

	if err := db.clients[header.API()].Delete(ctx, header); err != nil {
		return err
	}

	// Delete old node if it exists and no other connections remain.
	if s.NodeID() != "" {
		n, err := db.GetNode(ctx, node.Make(&dpb.Node{
			Header: &dpb.NodeHeader{
				Id:   s.NodeID(),
				Type: s.Header().Type(),
			},
		}).Header(), option.Remote(false))
		if err != nil {
			return err
		}

		if len(n.Sources()) == 0 {
			return db.DeleteNode(ctx, n.Header())
		}
	}
	return nil
}

func (db *DB) Get(ctx context.Context, header source.H, refresh option.Remote) (source.S, error) {
	if _, ok := db.clients[header.API()]; !ok {
		return source.S{}, nil
	}

	s, err := db.clients[header.API()].Get(ctx, header, refresh)
	if err != nil {
		return source.S{}, err
	}

	// If source is not found, attempt to force getting the remote version.
	if s == (source.S{}) && !refresh {
		return db.clients[header.API()].Get(ctx, header, option.Remote(true))
	}
	return s, nil
}

func (db *DB) Search(ctx context.Context, query string, opts map[epb.SourceAPI][]option.O) ([]source.S, error) {
	results := []source.S{}
	for api, _opts := range opts {
		if _, ok := db.clients[api]; ok {
			if api == epb.SourceAPI_SOURCE_API_TRUFFLE {
				_opts = append(_opts, option.Remote(true))
			}
			sources, err := db.clients[api].Search(
				ctx,
				query,
				_opts...,
			)
			if err != nil {
				return nil, err
			}
			results = append(results, sources...)
		}
	}

	slices.SortStableFunc(results, func(a, b source.S) int {
		u, _ := match.Hamming(query, a)
		v, _ := match.Hamming(query, b)
		return -1 * cmp.Compare(u, v)
	})

	return results, nil
}

func (db *DB) PB() *dpb.Database {
	sources := append(
		[]*dpb.Source{},
		append(
			db.clients[epb.SourceAPI_SOURCE_API_TRUFFLE].PB(),
			db.clients[epb.SourceAPI_SOURCE_API_MAL].PB()...,
		)...,
	)
	nodes := []*dpb.Node{}
	for _, n := range db.nodes {
		nodes = append(nodes, n.PB())
	}
	return &dpb.Database{
		Nodes:   nodes,
		Sources: sources,
	}
}
