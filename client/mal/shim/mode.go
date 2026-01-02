package shim

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/minkezhang/truffle-api/data/source/util"
	"github.com/nstratos/go-myanimelist/mal"
	"google.golang.org/protobuf/types/known/timestamppb"

	dpb "github.com/minkezhang/truffle-api/proto/go/data"
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

type Mode int

const (
	ModeAnime Mode = iota
	ModeManga
)

var (
	fields = map[Mode]mal.Fields{
		ModeAnime: mal.Fields{
			"media_type",
			"title",
			"alternative_titles",
			"main_picture",
			"mean",
			"synopsis",
			"genres",
			"related_anime", // MAL API does not return related_manga for anime entries
			"studios",
			"start_season",
		},
		ModeManga: mal.Fields{
			"media_type",
			"title",
			"alternative_titles",
			"main_picture",
			"mean",
			"synopsis",
			"genres",
			"related_manga", // MAL API does not return related_anime for manga entries
			"authors{last_name, first_name}",
			"updated_at",
		},
	}

	types = map[Mode]map[string]epb.SourceType{
		ModeAnime: map[string]epb.SourceType{
			"tv":         epb.SourceType_SOURCE_TYPE_SERIES_ANIME,
			"ova":        epb.SourceType_SOURCE_TYPE_SERIES_ANIME,
			"special":    epb.SourceType_SOURCE_TYPE_SERIES_ANIME,
			"tv_special": epb.SourceType_SOURCE_TYPE_SERIES_ANIME,
			"music":      epb.SourceType_SOURCE_TYPE_SERIES_ANIME,
			"ona":        epb.SourceType_SOURCE_TYPE_SERIES_ANIME,
			"movie":      epb.SourceType_SOURCE_TYPE_MOVIE_ANIME,
		},
		ModeManga: map[string]epb.SourceType{
			// MAL lists the "novel" type but experimentally, this is
			// "light_novel" instead.
			"light_novel": epb.SourceType_SOURCE_TYPE_BOOK_LIGHT_NOVEL,

			"manga":     epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
			"one_shot":  epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
			"doujinshi": epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
			"manhua":    epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
			"manhwa":    epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
			"oel":       epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
		},
	}
)

type Anime struct {
	mal.Anime
}

func (a Anime) PB() *dpb.Source {
	t, ok := types[ModeAnime][a.MediaType]
	if !ok {
		return nil
	}

	titles := []*dpb.Title{&dpb.Title{Title: a.Title}}

	if title := a.AlternativeTitles.En; title != "" {
		titles = append(titles, &dpb.Title{
			Title:        title,
			Localization: "en",
		})
	}
	if title := a.AlternativeTitles.Ja; title != "" {
		titles = append(titles, &dpb.Title{
			Title:        title,
			Localization: "ja",
		})
	}

	return &dpb.Source{
		Header: &dpb.SourceHeader{
			Api:  epb.SourceAPI_SOURCE_API_MAL,
			Type: t,
			Id:   strconv.FormatInt(int64(a.ID), 10),
		},
		Titles:     titles,
		PreviewUrl: a.MainPicture.Large,
		Score:      int64(a.Mean * 10),
		Synopsis:   a.Synopsis,
		Genres:     util.Apply(a.Genres, func(v mal.Genre) string { return v.Name }),
		RelatedHeaders: util.Apply(a.RelatedAnime, func(v mal.RelatedAnime) *dpb.SourceHeader {
			return &dpb.SourceHeader{
				Api: epb.SourceAPI_SOURCE_API_MAL_ANIME_PARTIAL,
				Id:  strconv.FormatInt(int64(v.Node.ID), 10),
			}
		}),
		Studios: util.Apply(a.Studios, func(v mal.Studio) string { return v.Name }),
		Seasons: []string{fmt.Sprintf("%v %d", strings.ToTitle(a.StartSeason.Season), a.StartSeason.Year)},
	}
}

type Manga struct {
	mal.Manga
}

func (m Manga) PB() *dpb.Source {
	t, ok := types[ModeManga][m.MediaType]
	if !ok {
		return nil
	}

	titles := []*dpb.Title{&dpb.Title{Title: m.Title}}

	if title := m.AlternativeTitles.En; title != "" {
		titles = append(titles, &dpb.Title{
			Title:        title,
			Localization: "en",
		})
	}
	if title := m.AlternativeTitles.Ja; title != "" {
		titles = append(titles, &dpb.Title{
			Title:        title,
			Localization: "ja",
		})
	}

	authors := []string{}
	illustrators := []string{}
	for _, s := range m.Authors {
		if strings.Contains(s.Role, "Story") {
			authors = append(authors, fmt.Sprintf("%s %s", s.Person.FirstName, s.Person.LastName))
		}
		if strings.Contains(s.Role, "Art") {
			illustrators = append(illustrators, fmt.Sprintf("%s %s", s.Person.FirstName, s.Person.LastName))
		}
	}

	return &dpb.Source{
		Header: &dpb.SourceHeader{
			Api:  epb.SourceAPI_SOURCE_API_MAL,
			Type: t,
			Id:   strconv.FormatInt(int64(m.ID), 10),
		},
		Titles:     titles,
		PreviewUrl: m.MainPicture.Large,
		Score:      int64(m.Mean * 10),
		Synopsis:   m.Synopsis,
		Genres:     util.Apply(m.Genres, func(v mal.Genre) string { return v.Name }),
		RelatedHeaders: util.Apply(m.RelatedManga, func(v mal.RelatedManga) *dpb.SourceHeader {
			return &dpb.SourceHeader{
				Api: epb.SourceAPI_SOURCE_API_MAL_MANGA_PARTIAL,
				Id:  strconv.FormatInt(int64(v.Node.ID), 10),
			}
		}),
		Authors:      authors,
		Illustrators: illustrators,
		LastUpdated:  timestamppb.New(m.UpdatedAt),
	}
}
