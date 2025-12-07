package mal

import (
	"context"
	"net/http"
	"strconv"

	"github.com/minkezhang/bene-api/client"
	"github.com/minkezhang/bene-api/client/query"
	"github.com/minkezhang/bene-api/db/atom"
	"github.com/nstratos/go-myanimelist/mal"

	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

var (
	_ client.C = &C{}
)

type O struct {
	ClientID         string
	PopularityCutoff int64
	MaxResults       int64
	NSFW             bool
}

func New(o O) *C {
	return &C{
		cutoff:  o.PopularityCutoff,
		results: o.MaxResults,
		nsfw:    o.NSFW,
		mal: *mal.NewClient(
			&http.Client{
				Transport: transport{id: o.ClientID},
			},
		),
	}
}

type C struct {
	cutoff  int64
	results int64
	nsfw    bool
	mal     mal.Client
}

func (c *C) APIType() epb.API { return epb.API_API_MAL }

func (c *C) Get(ctx context.Context, g query.G) (*atom.A, error) {
	id, err := strconv.ParseInt(g.ID, 10, 0)
	if err != nil {
		return nil, err
	}

	result, _, err := c.mal.Anime.Details(ctx, int(id), fields[g.AtomType])
	if err != nil {
		return nil, err
	}

	return FromAnime(g.AtomType, result), nil
}

func (c *C) Query(ctx context.Context, q *query.Q) ([]*atom.A, error) {
	return nil, nil
}

type transport struct {
	id string
}

func (t transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-MAL-CLIENT-ID", t.id)
	return http.DefaultTransport.RoundTrip(req)
}
