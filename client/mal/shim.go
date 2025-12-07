package mal

import (
	"fmt"
	"strconv"

	"github.com/minkezhang/bene-api/db/atom"
	"github.com/minkezhang/bene-api/db/atom/metadata"
	"github.com/minkezhang/bene-api/db/atom/metadata/movie"
	"github.com/minkezhang/bene-api/db/atom/metadata/tv"
	"github.com/nstratos/go-myanimelist/mal"

	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

var (
	// MAL returns movies and TV shows through the same endpoint, so we will
	// need to check.
	types = map[epb.Type]map[string]bool{
		epb.Type_TYPE_MOVIE: map[string]bool{
			"movie": true,
		},
		epb.Type_TYPE_TV: map[string]bool{
			"tv":      true,
			"ova":     true,
			"special": true,
			"ona":     true,
		},
	}

	fields = map[epb.Type]mal.Fields{
		epb.Type_TYPE_MOVIE: mal.Fields{
			"media_type",
			"title",
			"alternative_titles",
			"mean",
			"studios",
			"genres",
		},
		epb.Type_TYPE_TV: mal.Fields{
			"media_type",
			"title",
			"alternative_titles",
			"mean",
			"studios",
			"genres",
		},
		epb.Type_TYPE_BOOK: mal.Fields{
			"media_type",
			"title",
			"alternative_titles",
			"mean",
			"authors{first_name,last_name}",
		},
	}
)

func FromAnime(t epb.Type, r *mal.Anime) *atom.A {
	if !types[t][r.MediaType] {
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
		AtomType:   t,
		Metadata:   m,
	})
}
