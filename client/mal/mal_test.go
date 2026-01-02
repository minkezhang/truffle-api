package mal

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/minkezhang/truffle-api/client/option"
	"github.com/minkezhang/truffle-api/data/source"
	"google.golang.org/protobuf/testing/protocmp"

	cpb "github.com/minkezhang/truffle-api/proto/go/config"
	dpb "github.com/minkezhang/truffle-api/proto/go/data"
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

const (
	// This is a publically-known API key for the MAL Android app.
	MALClientID = "6114d00ca681b7701d1e15fe11a4987e"
)

var (
	config = &cpb.MAL{
		ClientId:   MALClientID,
		MaxResults: 2,
	}
)

func TestGet(t *testing.T) {
	configs := []struct {
		name        string
		api_type    epb.SourceAPI
		source_type epb.SourceType
		id          string
		want        source.S
	}{
		{
			name:        "Book/Manga",
			api_type:    epb.SourceAPI_SOURCE_API_MAL,
			source_type: epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
			id:          "1061", // Detective Conan
			want: source.Make(&dpb.Source{
				Header: &dpb.SourceHeader{
					Api:  epb.SourceAPI_SOURCE_API_MAL,
					Type: epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
					Id:   "1061",
				},
				Titles: []*dpb.Title{
					&dpb.Title{Title: "Meitantei Conan"},
					&dpb.Title{Title: "Case Closed", Localization: "en"},
					&dpb.Title{Title: "名探偵コナン", Localization: "ja"},
				},
				Score:        82,
				Genres:       []string{"Adventure", "Award Winning", "Comedy", "Detective", "Mystery", "Shounen"},
				Authors:      []string{"Gosho Aoyama"},
				Illustrators: []string{"Gosho Aoyama"},
			}),
		},
		{
			name:        "Book/LightNovel",
			api_type:    epb.SourceAPI_SOURCE_API_MAL,
			source_type: epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
			id:          "86769", // Apothecary Diaries
			want: source.Make(&dpb.Source{
				Header: &dpb.SourceHeader{
					Api:  epb.SourceAPI_SOURCE_API_MAL,
					Id:   "86769",
					Type: epb.SourceType_SOURCE_TYPE_BOOK_LIGHT_NOVEL,
				},
				Titles: []*dpb.Title{
					&dpb.Title{Title: "Kusuriya no Hitorigoto"},
					&dpb.Title{Title: "The Apothecary Diaries", Localization: "en"},
					&dpb.Title{Title: "薬屋のひとりごと", Localization: "ja"},
				},
				Score:        88,
				Genres:       []string{"Drama", "Medical", "Mystery"},
				Authors:      []string{"Natsu Hyuuga"},
				Illustrators: []string{"Touko Shino"},
			}),
		},
		{
			name:        "Series",
			api_type:    epb.SourceAPI_SOURCE_API_MAL,
			source_type: epb.SourceType_SOURCE_TYPE_SERIES_ANIME,
			id:          "235", // Detective Conan
			want: source.Make(&dpb.Source{
				Header: &dpb.SourceHeader{
					Api:  epb.SourceAPI_SOURCE_API_MAL,
					Id:   "235",
					Type: epb.SourceType_SOURCE_TYPE_SERIES_ANIME,
				},
				Titles: []*dpb.Title{
					&dpb.Title{Title: "Meitantei Conan"},
					&dpb.Title{Title: "Case Closed", Localization: "en"},
					&dpb.Title{Title: "名探偵コナン", Localization: "ja"},
				},
				Score:   81,
				Genres:  []string{"Adventure", "Comedy", "Detective", "Mystery", "Shounen"},
				Studios: []string{"TMS Entertainment"},
			}),
		},
		{
			name:        "Movie",
			api_type:    epb.SourceAPI_SOURCE_API_MAL,
			source_type: epb.SourceType_SOURCE_TYPE_MOVIE_ANIME,
			id:          "28851", // Koe no Katachi
			want: source.Make(&dpb.Source{
				Header: &dpb.SourceHeader{
					Api:  epb.SourceAPI_SOURCE_API_MAL,
					Id:   "28851",
					Type: epb.SourceType_SOURCE_TYPE_MOVIE_ANIME,
				},
				Titles: []*dpb.Title{
					&dpb.Title{Title: "Koe no Katachi"},
					&dpb.Title{Title: "A Silent Voice", Localization: "en"},
					&dpb.Title{Title: "聲の形", Localization: "ja"},
				},
				Score:   89,
				Genres:  []string{"Award Winning", "Drama", "Shounen"},
				Studios: []string{"Kyoto Animation"},
			}),
		},
		{
			name:        "Movie/Partial",
			api_type:    epb.SourceAPI_SOURCE_API_MAL_ANIME_PARTIAL,
			source_type: epb.SourceType_SOURCE_TYPE_UNKNOWN,
			id:          "28851",
			want: source.Make(&dpb.Source{
				Header: &dpb.SourceHeader{
					Api:  epb.SourceAPI_SOURCE_API_MAL,
					Id:   "28851",
					Type: epb.SourceType_SOURCE_TYPE_MOVIE_ANIME,
				},
				Titles: []*dpb.Title{
					&dpb.Title{Title: "Koe no Katachi"},
					&dpb.Title{Title: "A Silent Voice", Localization: "en"},
					&dpb.Title{Title: "聲の形", Localization: "ja"},
				},
				Score:   89,
				Genres:  []string{"Award Winning", "Drama", "Shounen"},
				Studios: []string{"Kyoto Animation"},
			}),
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			client := Make(config)

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			got, err := client.Get(ctx, source.Make(
				&dpb.Source{
					Header: &dpb.SourceHeader{
						Api:  c.api_type,
						Type: c.source_type,
						Id:   c.id,
					},
				},
			).Header())
			if err != nil {
				t.Fatalf("Get() returned unexpected error: %v", err)
			}

			if diff := cmp.Diff(
				c.want,
				got,
				cmp.AllowUnexported(source.S{}),
				protocmp.Transform(),
				protocmp.IgnoreFields(
					&dpb.Source{},
					"synopsis",
					"last_updated",
					"preview_url",
					"seasons",
					"score",
					"related_headers",
				),
			); diff != "" {
				t.Errorf("Get() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	t.Run("NSFW", func(t *testing.T) {
		client := Make(&cpb.MAL{
			ClientId:   config.GetClientId(),
			MaxResults: 20,
		})
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		got, err := client.Search(ctx, "Citrus", option.NSFW(true)) // NSFW
		if err != nil {
			t.Errorf("Search() returned non-nil error: %v", err)
		}

		// Check that at least one NSFW result was returned.
		nsfw := false
		for _, s := range got {
			categories := map[string]bool{}
			for _, g := range s.Genres() {
				categories[strings.ToLower(g)] = true
			}
			if categories["hentai"] || categories["erotica"] {
				nsfw = true
			}
		}
		if !nsfw {
			t.Errorf("Search() returned no NSFW results")
		}
	})
	t.Run("NoNSFW", func(t *testing.T) {
		client := Make(&cpb.MAL{
			ClientId:   config.GetClientId(),
			MaxResults: 20,
		})
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		got, err := client.Search(ctx, "Citrus", option.NSFW(false)) // NSFW
		if err != nil {
			t.Errorf("Search() returned non-nil error: %v", err)
		}
		for _, s := range got {
			categories := map[string]bool{}
			for _, g := range s.Genres() {
				categories[strings.ToLower(g)] = true
			}
			if categories["hentai"] || categories["erotica"] {
				t.Errorf(
					"Search() returned a NSFW result: %v (%v:%v)",
					s.Title().Title(),
					s.Header().Type().String(),
					s.Header().ID(),
				)
			}
		}
	})
}
