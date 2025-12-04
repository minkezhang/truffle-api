package client

import (
	"context"
	"regexp"

	"github.com/minkezhang/bene-api/client/query"
	"github.com/minkezhang/bene-api/db/atom"

	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

type C interface {
	APIType() epb.API

	// Get returns a single atom given an ID associated with the API.
	Get(ctx context.Context, id string) (*atom.A, error)

	// Query returns a (potentially empty) list of atoms with the given
	// input.
	Query(ctx context.Context, q query.Q) ([]*atom.A, error)
}

func Match(q query.Q, a *atom.A) (bool, error) {
	pattern, err := regexp.Compile(q.Title)
	if err != nil {
		return false, err
	}

	for _, t := range a.Titles() {
		if pattern.MatchString(t.Title) {
			return true, nil
		}
	}
	return false, nil

}
