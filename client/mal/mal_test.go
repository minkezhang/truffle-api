package mal

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/minkezhang/bene-api/client/query"
	"github.com/minkezhang/bene-api/db/atom"
	"github.com/minkezhang/bene-api/db/atom/metadata/movie"
	"github.com/minkezhang/bene-api/db/atom/metadata/shared/video"
	"github.com/minkezhang/bene-api/db/atom/metadata/tv"

	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

const (
	// This is a publically-known API key for the MAL Android app.
	MALClientID = "6114d00ca681b7701d1e15fe11a4987e"
)

func TestGet(t *testing.T) {
	configs := []struct {
		name     string
		atomType epb.Type
		id       string
		want     *atom.A
	}{
		{
			name:     "TV",
			atomType: epb.Type_TYPE_TV,
			id:       "235", // Detective Conan
			want: atom.New(atom.O{
				APIType: epb.API_API_MAL,
				APIID:   "235",
				Titles: []atom.T{
					atom.T{Title: "Meitantei Conan"},
					atom.T{Title: "Case Closed", Localization: "en"},
					atom.T{Title: "名探偵コナン", Localization: "ja"},
				},
				PreviewURL: "https://cdn.myanimelist.net/images/anime/7/75199l.jpg",
				Score:      81,
				AtomType:   epb.Type_TYPE_TV,
				Metadata: tv.New(tv.O{
					IsAnimated: true,
					IsAnime:    true,
					Genres:     []string{"Adventure", "Comedy", "Detective", "Mystery", "Shounen"},
					Studios:    []string{"TMS Entertainment"},
				}),
			}),
		},
		{
			name:     "Movie",
			atomType: epb.Type_TYPE_MOVIE,
			id:       "28851", // Koe no Katachi
			want: atom.New(atom.O{
				APIType: epb.API_API_MAL,
				APIID:   "28851",
				Titles: []atom.T{
					atom.T{Title: "Koe no Katachi"},
					atom.T{Title: "A Silent Voice", Localization: "en"},
					atom.T{Title: "聲の形", Localization: "ja"},
				},
				PreviewURL: "https://cdn.myanimelist.net/images/anime/1122/96435l.webp",
				Score:      89,
				AtomType:   epb.Type_TYPE_MOVIE,
				Metadata: movie.New(movie.O{
					IsAnimated: true,
					IsAnime:    true,
					Genres:     []string{"Award Winning", "Drama", "Shounen"},
					Studios:    []string{"Kyoto Animation"},
				}),
			}),
		},
		{
			name:     "IncorrectType",
			atomType: epb.Type_TYPE_MOVIE,
			id:       "235",
			want:     nil,
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			client := New(O{
				ClientID: MALClientID,
			})

			got, err := client.Get(context.Background(), query.G{
				AtomType: c.atomType,
				ID:       c.id,
			})
			if err != nil {
				t.Errorf("Get() returned unexpected error: %v", err)
			}

			if diff := cmp.Diff(
				c.want,
				got,
				cmp.AllowUnexported(
					atom.A{},
					video.M{},
				),
				cmpopts.IgnoreFields(
					atom.A{},
					"synopsis",
					"previewURL",
				),
			); diff != "" {
				t.Errorf("Get() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestList(t *testing.T) {
	configs := []struct {
		name  string
		query *query.Q
		want  []*atom.A
	}{
		{
			name: "Base",
			query: query.New(query.O{
				AtomTypes: []epb.Type{
					epb.Type_TYPE_TV,
					epb.Type_TYPE_MOVIE,
				},
				Title: "Digimon Tamers",
			}),
			want: []*atom.A{
				atom.New(atom.O{
					APIType:  epb.API_API_MAL,
					APIID:    "874", // Digimon Tamers
					AtomType: epb.Type_TYPE_TV,
				}),
				atom.New(atom.O{
					APIType:  epb.API_API_MAL,
					APIID:    "3033", // Digimon Tamers: Runaway Locomon
					AtomType: epb.Type_TYPE_MOVIE,
				}),
			},
		},
		{
			name: "FilterType",
			query: query.New(query.O{
				AtomTypes: []epb.Type{
					epb.Type_TYPE_MOVIE,
				},
				Title: "Digimon Tamers",
			}),
			want: []*atom.A{
				atom.New(atom.O{
					APIType:  epb.API_API_MAL,
					APIID:    "3033", // Digimon Tamers: Runaway Locomon
					AtomType: epb.Type_TYPE_MOVIE,
				}),
				atom.New(atom.O{
					APIType:  epb.API_API_MAL,
					APIID:    "3032", // Digimon Tamers: Battle of Adventurers
					AtomType: epb.Type_TYPE_MOVIE,
				}),
			},
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {

			client := New(O{
				ClientID:         MALClientID,
				PopularityCutoff: 10000,
				MaxResults:       2,
				NSFW:             true,
			})
			got, err := client.Query(context.Background(), c.query)

			if err != nil {
				t.Errorf("Query() returned unexpected error: %v", err)
			}

			if diff := cmp.Diff(
				c.want,
				got,
				cmp.AllowUnexported(
					atom.A{},
					video.M{},
				),
				cmpopts.IgnoreFields(
					atom.A{},
					"synopsis",
					"previewURL",
					"metadata",
					"score",
					"titles",
				),
			); diff != "" {
				t.Errorf("Query() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
