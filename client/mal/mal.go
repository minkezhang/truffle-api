package mal

import (
	"context"
	"math"
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
	var res []*atom.A

	// MAL API sometimes returns duplicate entries
	unique := map[string]bool{}

	f := func(resp *mal.Response) ([]mal.Anime, *mal.Response, error) {
		if resp != nil && resp.NextOffset == 0 {
			return nil, nil, nil
		}
		var offset int
		if resp != nil {
			offset = resp.NextOffset
		}
		results, resp, err := c.mal.Anime.List(
			ctx,
			q.Title(),
			// Union TV and Movie atom type results
			fields[epb.Type_TYPE_TV],
			mal.Limit(math.Min(100, float64(c.results))),
			mal.Offset(offset),
			mal.NSFW(c.nsfw),
		)
		return results, resp, err
	}

	for _, t := range []epb.Type{epb.Type_TYPE_TV, epb.Type_TYPE_MOVIE} {
		if q.IsSupportedType(t) {
			var results []mal.Anime

			var page []mal.Anime
			var resp *mal.Response
			var err error

			// Aggregate all results
			for page, resp, err = f(nil); err == nil && page != nil && len(results) <= int(c.results); page, resp, err = f(resp) {
				results = append(results, page...)
			}
			if err != nil {
				return nil, err
			}

			for _, r := range results {
				// Trim obscure series
				if popularity := c.cutoff; popularity >= 0 && r.Popularity >= int(popularity) {
					continue
				}
				if !types[t][r.MediaType] {
					continue
				}

				if !unique[strconv.FormatInt(int64(r.ID), 10)] {
					unique[strconv.FormatInt(int64(r.ID), 10)] = true
					res = append(res, FromAnime(t, &r))
				}
			}
		}
	}

	return res, nil
}

type transport struct {
	id string
}

func (t transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-MAL-CLIENT-ID", t.id)
	return http.DefaultTransport.RoundTrip(req)
}
