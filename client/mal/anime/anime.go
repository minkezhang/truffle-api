package anime

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/minkezhang/bene-api/client/query"
	"github.com/minkezhang/bene-api/db/atom"
	"github.com/minkezhang/bene-api/db/atom/metadata"
	"github.com/minkezhang/bene-api/db/atom/metadata/movie"
	"github.com/minkezhang/bene-api/db/atom/metadata/tv"
	"github.com/nstratos/go-myanimelist/mal"

	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

var (
	fields = mal.Fields{
		"media_type",
		"title",
		"alternative_titles",
		"mean",
		"studios",
		"synopsis",
		"genres",
		"popularity",
	}

	types = map[string]epb.Type{
		"tv":      epb.Type_TYPE_TV,
		"ova":     epb.Type_TYPE_TV,
		"special": epb.Type_TYPE_TV,
		"ona":     epb.Type_TYPE_TV,
		"movie":   epb.Type_TYPE_MOVIE,
	}
)

type C struct {
	Cutoff  int
	Results int
	NSFW    bool
	MAL     mal.Client
}

func (c *C) Get(ctx context.Context, g query.G) (*atom.A, error) {
	id, err := strconv.ParseInt(g.ID, 10, 0)
	if err != nil {
		return nil, err
	}

	result, _, err := c.MAL.Anime.Details(ctx, int(id), fields)
	if err != nil {
		return nil, err
	}

	return ToBene(*result, []epb.Type{g.AtomType}), nil
}

type candidate struct {
	similarity float64
	atom       *atom.A
}

func (c *C) Query(ctx context.Context, q *query.Q) ([]*atom.A, error) {
	f := func(resp *mal.Response) ([]*atom.A, *mal.Response, error) {
		if resp != nil && resp.NextOffset == 0 {
			return nil, nil, nil
		}
		var offset int
		if resp != nil {
			offset = resp.NextOffset
		}
		results, resp, err := c.MAL.Anime.List(
			ctx,
			q.Title(),
			fields,
			mal.Limit(100),
			mal.Offset(offset),
			mal.NSFW(c.NSFW),
		)
		page := []*atom.A{}
		for _, r := range results {
			a := ToBene(r, q.AtomTypes())
			if a != nil && r.Popularity <= c.Cutoff {
				page = append(page, a)
			}
		}
		return page, resp, err
	}

	// MAL API sometimes returns duplicate entries
	candidates := map[string]candidate{}

	var page []*atom.A
	var resp *mal.Response
	var err error

	// Aggregate all results
	for page, resp, err = f(nil); err == nil && page != nil; page, resp, err = f(resp) {
		for _, a := range page {
			if !q.IsSupportedType(a.AtomType()) {
				continue
			}
			h, err := query.Hamming(q, a)
			if err != nil {
				return nil, err
			}
			candidates[a.APIID()] = candidate{
				similarity: h,
				atom:       a,
			}
		}
		if len(candidates) >= c.Results {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	// Sort candidates by title similarity
	cl := []candidate{}
	for _, c := range candidates {
		cl = append(cl, c)
	}
	sort.Slice(cl, func(i, j int) bool {
		return cl[i].similarity > cl[j].similarity
	})
	res := []*atom.A{}
	for _, r := range cl {
		res = append(res, r.atom)
	}
	end := len(res)
	if c.Results < end {
		end = c.Results
	}
	return res[:end], nil
}

func ToBene(r mal.Anime, ts []epb.Type) *atom.A {
	t, ok := types[r.MediaType]
	if !ok {
		return nil
	}

	l := map[epb.Type]bool{}
	for _, u := range ts {
		l[u] = true
	}
	if !l[t] {
		return nil
	}

	titles := []atom.T{
		atom.T{Title: r.Title},
	}
	if title := r.AlternativeTitles.En; title != "" {
		titles = append(titles, atom.T{
			Title:        title,
			Localization: "en",
		})
	}
	if title := r.AlternativeTitles.Ja; title != "" {
		titles = append(titles, atom.T{
			Title:        title,
			Localization: "ja",
		})
	}

	genres := []string{}
	for _, g := range r.Genres {
		genres = append(genres, g.Name)
	}

	studios := []string{}
	for _, s := range r.Studios {
		studios = append(studios, s.Name)
	}

	var m metadata.M
	switch t {
	case epb.Type_TYPE_TV:
		m = tv.New(tv.O{
			IsAnimated: true,
			IsAnime:    true,
			Genres:     genres,
			Studios:    studios,
		})
	case epb.Type_TYPE_MOVIE:
		m = movie.New(movie.O{
			IsAnimated: true,
			IsAnime:    true,
			Genres:     genres,
			Studios:    studios,
		})
	default:
		panic(fmt.Errorf("invalid switch type: %v", t))
	}

	return atom.New(atom.O{
		APIType:    epb.API_API_MAL,
		APIID:      strconv.FormatInt(int64(r.ID), 10),
		Titles:     titles,
		PreviewURL: r.MainPicture.Large,
		Score:      int64(r.Mean * 10),
		Synopsis:   r.Synopsis,
		AtomType:   t,
		Metadata:   m,
	})
}
