package mal

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/minkezhang/bene-api/client"
	"github.com/minkezhang/bene-api/client/query"
	"github.com/minkezhang/bene-api/db/atom"
	"github.com/minkezhang/bene-api/db/atom/metadata/tv"
	"github.com/nstratos/go-myanimelist/mal"

	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

var (
	_ client.C = &C{}

	fields = map[epb.Type]mal.Fields{
		epb.Type_TYPE_TV: mal.Fields{
			"title",
			"mean",
			"studios",
			"alternative_titles",
			"genres",
		},
		epb.Type_TYPE_BOOK: mal.Fields{
			"media_type",
			"popularity",
			"title",
			"alternative_titles",
			"mean",
			"authors{first_name,last_name}",
		},
	}
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

	switch t := g.AtomType; t {
	case epb.Type_TYPE_TV:
		result, _, err := c.mal.Anime.Details(ctx, int(id), fields[epb.Type_TYPE_TV])
		if err != nil {
			return nil, err
		}

		titles := []atom.T{
			atom.T{
				Title: result.Title,
			},
		}
		if t := result.AlternativeTitles.En; t != "" {
			titles = append(titles, atom.T{
				Title:        t,
				Localization: "en",
			})
		}

		genres := []string{}
		for _, g := range result.Genres {
			genres = append(genres, g.Name)
		}

		studios := []string{}
		for _, s := range result.Studios {
			studios = append(studios, s.Name)
		}

		m := tv.New(tv.O{
			IsAnimated: true,
			IsAnime:    true,
			Genres:     genres,
			Studios:    studios,
		})

		a := atom.New(atom.O{
			APIType:    c.APIType(),
			APIID:      strconv.FormatInt(int64(result.ID), 10),
			Titles:     titles,
			PreviewURL: result.MainPicture.Large,
			Score:      int64(result.Mean * 10),
			AtomType:   t,
			Metadata:   m,
		})
		return a, nil
	case epb.Type_TYPE_MOVIE:
		return nil, fmt.Errorf("unimplemented")
	case epb.Type_TYPE_BOOK:
		return nil, fmt.Errorf("unimplemented")
	}

	return nil, nil
}

func (c *C) Query(ctx context.Context, q *query.Q) ([]*atom.A, error) {
	return nil, nil /* unimplmented */
}

type transport struct {
	id string
}

func (t transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-MAL-CLIENT-ID", t.id)
	return http.DefaultTransport.RoundTrip(req)
}
