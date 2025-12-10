package manga

import (
	"context"
	"fmt"
	"strings"
	// "sort"
	"strconv"

	"github.com/minkezhang/bene-api/client/query"
	"github.com/minkezhang/bene-api/db/atom"
	"github.com/minkezhang/bene-api/db/atom/metadata"
	"github.com/minkezhang/bene-api/db/atom/metadata/book"
	"github.com/nstratos/go-myanimelist/mal"

	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

var (
	fields = mal.Fields{
		"media_type",
		"title",
		"alternative_titles",
		"main_picture",
		"synopsis",
		"mean",
		"popularity",
		"genres",
		"authors{last_name, first_name}",
	}

	types = map[string]bool{
		// MAL lists the "novel" type but experimentally, this is
		// "light_novel" instead.
		"light_novel": false,

		"manga":     true,
		"one_shot":  true,
		"doujinshi": true,
		"manhua":    true,
		"manhwa":    true,
		"oel":       true,
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

	result, _, err := c.MAL.Manga.Details(ctx, int(id), fields)
	if err != nil {
		return nil, err
	}

	return ToBene(*result, []epb.Type{g.AtomType}), nil
}

func ToBene(r mal.Manga, ts []epb.Type) *atom.A {
	t := epb.Type_TYPE_BOOK
	im, ok := types[r.MediaType]
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

	authors := []string{}
	illustrators := []string{}
	for _, s := range r.Authors {
		if strings.Contains(s.Role, "Story") {
			authors = append(authors, fmt.Sprintf("%s %s", s.Person.FirstName, s.Person.LastName))
		}
		if strings.Contains(s.Role, "Art") {

			illustrators = append(illustrators, fmt.Sprintf("%s %s", s.Person.FirstName, s.Person.LastName))
		}
	}

	var m metadata.M = book.New(book.O{
		Genres:        genres,
		Illustrators:  illustrators,
		Authors:       authors,
		IsIllustrated: im,
		IsManga:       true,
	})

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
