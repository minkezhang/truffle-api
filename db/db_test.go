package db

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/minkezhang/truffle-api/client/option"
	"github.com/minkezhang/truffle-api/data/node"
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
	config = &cpb.Config{
		Mal: &cpb.MAL{
			ClientId:   MALClientID,
			MaxResults: 2,
		},
		Truffle: &cpb.Truffle{},
	}

	frieren_truffle = &dpb.Source{
		NodeId: "frieren",
		Header: &dpb.SourceHeader{
			Api:  epb.SourceAPI_SOURCE_API_TRUFFLE,
			Type: epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
			Id:   "bar",
		},
		Titles: []*dpb.Title{
			&dpb.Title{Title: "Frieren"},
		},
	}

	frieren_mal_header_only = &dpb.Source{
		NodeId: "frieren",
		Header: &dpb.SourceHeader{
			Api:  epb.SourceAPI_SOURCE_API_MAL,
			Type: epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
			Id:   "126287", // Frieren
		},
		Titles: []*dpb.Title{
			&dpb.Title{Title: "Frieren: Beyond Journey's End", Localization: "en"},
		},
	}

	frieren_mal_with_data = &dpb.Source{
		NodeId: "frieren",
		Header: frieren_mal_header_only.GetHeader(),
		Titles: []*dpb.Title{
			&dpb.Title{Title: "Frieren: Beyond Journey's End", Localization: "en"},
			&dpb.Title{Title: "Sousou no Frieren", Localization: ""},
			&dpb.Title{Title: "葬送のフリーレン", Localization: "ja"},
		},
		Authors:      []string{"Kanehito Yamada"},
		Illustrators: []string{"Tsukasa Abe"},
		Genres:       []string{"Adventure", "Award Winning", "Drama", "Fantasy", "Shounen"},
	}

	conan_mal_header_only = &dpb.Source{
		NodeId: "",
		Header: &dpb.SourceHeader{
			Api:  epb.SourceAPI_SOURCE_API_MAL,
			Type: epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
			Id:   "1061", // Detective Conan
		},
		Titles: []*dpb.Title{
			&dpb.Title{Title: "Meitantei Conan"},
		},
	}

	conan_mal_with_data = &dpb.Source{
		NodeId: "",
		Header: &dpb.SourceHeader{
			Api:  epb.SourceAPI_SOURCE_API_MAL,
			Type: epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
			Id:   "1061", // Detective Conan
		},
		Titles: []*dpb.Title{
			&dpb.Title{Title: "Case Closed", Localization: "en"},
			&dpb.Title{Title: "Meitantei Conan"},
			&dpb.Title{Title: "名探偵コナン", Localization: "ja"},
		},
		Authors:      []string{"Gosho Aoyama"},
		Illustrators: []string{"Gosho Aoyama"},
		Genres:       []string{"Adventure", "Award Winning", "Comedy", "Detective", "Mystery", "Shounen"},
	}

	spyxfamily_mal_header_only = &dpb.Source{
		NodeId: "spyxfamily",
		Header: &dpb.SourceHeader{
			Api:  epb.SourceAPI_SOURCE_API_MAL,
			Type: epb.SourceType_SOURCE_TYPE_SERIES_ANIME,
			Id:   "50265", // Spy x Family
		},
		Titles: []*dpb.Title{
			&dpb.Title{Title: "Spy x Family"},
		},
	}

	data = &dpb.Database{
		Nodes: []*dpb.Node{
			&dpb.Node{
				Header: &dpb.NodeHeader{
					Id:   "frieren",
					Type: epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
				},
			},
			&dpb.Node{
				Header: &dpb.NodeHeader{
					Id:   "spyxfamily",
					Type: epb.SourceType_SOURCE_TYPE_SERIES_ANIME,
				},
			},
		},
		Sources: []*dpb.Source{
			frieren_truffle,
			frieren_mal_header_only,
			conan_mal_header_only,
			spyxfamily_mal_header_only,
		},
	}
)

func TestGetNode(t *testing.T) {
	db := New(context.Background(), config, data)

	got, err := db.GetNode(context.Background(), node.Make(&dpb.Node{
		Header: &dpb.NodeHeader{
			Id:   "frieren",
			Type: epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
		},
	}).Header(), true)
	if err != nil {
		t.Errorf("GetNode() returned non-nil error: %v", err)
	}

	want := node.Make(&dpb.Node{
		Header: &dpb.NodeHeader{
			Type: epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
			Id:   "frieren",
		},
	}).WithSources([]source.S{
		source.Make(frieren_truffle),
		source.Make(&dpb.Source{
			NodeId: "frieren",
			Header: frieren_mal_header_only.GetHeader(),
			Titles: []*dpb.Title{
				&dpb.Title{Title: "Frieren: Beyond Journey's End", Localization: "en"},
				&dpb.Title{Title: "Sousou no Frieren", Localization: ""},
				&dpb.Title{Title: "葬送のフリーレン", Localization: "ja"},
			},
			Authors:      []string{"Kanehito Yamada"},
			Illustrators: []string{"Tsukasa Abe"},
			Genres:       []string{"Adventure", "Award Winning", "Drama", "Fantasy", "Shounen"},
		}),
	})

	if diff := cmp.Diff(
		want,
		got,
		cmp.AllowUnexported(node.N{}, source.S{}),
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
		t.Errorf("GetNode() mismatch (-want +got):\n%v", diff)
	}
}

func TestGet(t *testing.T) {
	t.Run("NoRefresh", func(t *testing.T) {
		db := New(context.Background(), config, data)

		got, err := db.Get(context.Background(), source.Make(frieren_mal_header_only).Header(), false)
		if err != nil {
			t.Errorf("Get() returned non-nil error: %v", err)
		}

		want := source.Make(frieren_mal_header_only)
		if diff := cmp.Diff(
			want,
			got,
			cmp.AllowUnexported(source.S{}),
			protocmp.Transform(),
		); diff != "" {
			t.Errorf("Get() mismatch (-want +got):\n%v", diff)
		}
	})
	t.Run("Refresh", func(t *testing.T) {
		db := New(context.Background(), config, data)

		got, err := db.Get(context.Background(), source.Make(frieren_mal_header_only).Header(), true)
		if err != nil {
			t.Errorf("Get() returned non-nil error: %v", err)
		}

		want := source.Make(frieren_mal_with_data)
		if diff := cmp.Diff(
			want,
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
			t.Errorf("Get() mismatch (-want +got):\n%v", diff)
		}
	})
}

func TestPut(t *testing.T) {
	t.Run("AddSource/NewNode", func(t *testing.T) {
		db := New(
			context.Background(),
			config,
			&dpb.Database{
				Nodes:   []*dpb.Node{},
				Sources: []*dpb.Source{},
			},
		)

		h, err := db.Put(context.Background(), source.Make(conan_mal_header_only))
		if err != nil {
			t.Errorf("Put() returned non-nil error: %v", err)
		}

		s, err := db.Get(context.Background(), h, option.Remote(false))
		if err != nil {
			t.Errorf("Get() returned non-nil error: %v", err)
		}

		got, err := db.GetNode(context.Background(), node.Make(&dpb.Node{
			Header: &dpb.NodeHeader{
				Id:   s.NodeID(),
				Type: s.Header().Type(),
			}}).Header(), option.Remote(false))
		if err != nil {
			t.Errorf("GetNode() returned non-nil error: %v", err)
		}

		want := node.Make(&dpb.Node{
			Header: &dpb.NodeHeader{
				Id:   s.NodeID(),
				Type: s.Header().Type(),
			},
		}).WithSources([]source.S{s})

		if diff := cmp.Diff(
			want,
			got,
			cmp.AllowUnexported(node.N{}, source.S{}),
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
			t.Errorf("GetNode() mismatch (-want +got):\n%v", diff)
		}

	})
	t.Run("AddSource/LinkNode", func(t *testing.T) {
		db := New(
			context.Background(),
			config,
			&dpb.Database{
				Nodes: []*dpb.Node{
					&dpb.Node{
						Header: &dpb.NodeHeader{
							Id:   "frieren",
							Type: epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
						},
					},
				},
				Sources: []*dpb.Source{
					frieren_truffle,
				},
			},
		)

		h, err := db.Put(context.Background(), source.Make(frieren_mal_header_only).WithNodeID("frieren"))
		if err != nil {
			t.Errorf("Put() returned non-nil error: %v", err)
		}

		s, err := db.Get(context.Background(), h, option.Remote(false))
		if err != nil {
			t.Errorf("Get() returned non-nil error: %v", err)
		}

		got, err := db.GetNode(context.Background(), node.Make(&dpb.Node{
			Header: &dpb.NodeHeader{
				Id:   s.NodeID(),
				Type: s.Header().Type(),
			}}).Header(), option.Remote(false))
		if err != nil {
			t.Errorf("GetNode() returned non-nil error: %v", err)
		}

		want := node.Make(&dpb.Node{
			Header: &dpb.NodeHeader{
				Id:   s.NodeID(),
				Type: s.Header().Type(),
			},
		}).WithSources([]source.S{source.Make(frieren_truffle), s})

		if diff := cmp.Diff(
			want,
			got,
			cmp.AllowUnexported(node.N{}, source.S{}),
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
			t.Errorf("GetNode() mismatch (-want +got):\n%v", diff)
		}
	})

	t.Run("RelinkSource/NoRemoveNode", func(t *testing.T) {
		db := New(
			context.Background(),
			config,
			&dpb.Database{
				Nodes: []*dpb.Node{
					&dpb.Node{
						Header: &dpb.NodeHeader{
							Id:   "frieren",
							Type: epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
						},
					},
				},
				Sources: []*dpb.Source{
					frieren_truffle,
					frieren_mal_header_only,
				},
			},
		)

		h, err := db.Put(context.Background(), source.Make(frieren_mal_header_only).WithNodeID(""))
		if err != nil {
			t.Errorf("Put() returned non-nil error: %v", err)
		}

		s, err := db.Get(context.Background(), h, option.Remote(false))
		if err != nil {
			t.Errorf("Get() returned non-nil error: %v", err)
		}

		got, err := db.GetNode(context.Background(), node.Make(&dpb.Node{
			Header: &dpb.NodeHeader{
				Id:   s.NodeID(),
				Type: s.Header().Type(),
			}}).Header(), option.Remote(false))
		if err != nil {
			t.Errorf("GetNode() returned non-nil error: %v", err)
		}

		want := node.Make(&dpb.Node{
			Header: &dpb.NodeHeader{
				Id:   s.NodeID(),
				Type: s.Header().Type(),
			},
		}).WithSources([]source.S{s})

		if diff := cmp.Diff(
			want,
			got,
			cmp.AllowUnexported(node.N{}, source.S{}),
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
			t.Errorf("GetNode() mismatch (-want +got):\n%v", diff)
		}
	})

	t.Run("RelinkSource/RemoveNode", func(t *testing.T) {
		db := New(
			context.Background(),
			config,
			&dpb.Database{
				Nodes: []*dpb.Node{
					&dpb.Node{
						Header: &dpb.NodeHeader{
							Id:   "frieren",
							Type: epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
						},
					},
				},
				Sources: []*dpb.Source{
					frieren_truffle,
				},
			},
		)

		h, err := db.Put(context.Background(), source.Make(frieren_truffle).WithNodeID(""))
		if err != nil {
			t.Errorf("Put() returned non-nil error: %v", err)
		}

		s, err := db.Get(context.Background(), h, option.Remote(false))
		if err != nil {
			t.Errorf("Get() returned non-nil error: %v", err)
		}

		got, err := db.GetNode(context.Background(), node.Make(&dpb.Node{
			Header: &dpb.NodeHeader{
				Id:   "frieren",
				Type: s.Header().Type(),
			}}).Header(), option.Remote(false))
		if err != nil {
			t.Errorf("GetNode() returned non-nil error: %v", err)
		}

		want := node.N{}

		if diff := cmp.Diff(
			want,
			got,
			cmp.AllowUnexported(node.N{}, source.S{}),
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
			t.Errorf("GetNode() mismatch (-want +got):\n%v", diff)
		}
	})
}

func TestSearch(t *testing.T) {
	t.Run("NoRemote", func(t *testing.T) {
		db := New(
			context.Background(),
			config,
			&dpb.Database{
				Nodes: []*dpb.Node{
					&dpb.Node{
						Header: &dpb.NodeHeader{
							Id:   "frieren",
							Type: epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
						},
					},
				},
				Sources: []*dpb.Source{
					frieren_truffle,
					frieren_mal_header_only,
					&dpb.Source{
						NodeId: "",
						Header: &dpb.SourceHeader{
							Api:  epb.SourceAPI_SOURCE_API_MAL,
							Type: epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
							Id:   "1061", // Detective Conan
						},
					},
				},
			},
		)

		got, err := db.Search(context.Background(), "frieren", map[epb.SourceAPI][]option.O{
			epb.SourceAPI_SOURCE_API_MAL: []option.O{option.Remote(false)},
		})
		if err != nil {
			t.Errorf("Search() returned non-nil error: %v", err)
		}

		want := []source.S{source.Make(frieren_mal_header_only)}
		if diff := cmp.Diff(
			want,
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
			t.Errorf("Search() mismatch (-want +got):\n%v", diff)
		}
	})
	t.Run("Truffle", func(t *testing.T) {
		db := New(context.Background(), config, data)

		got, err := db.Search(context.Background(), "frieren", map[epb.SourceAPI][]option.O{
			epb.SourceAPI_SOURCE_API_TRUFFLE: []option.O{option.Remote(false)},
		})
		if err != nil {
			t.Errorf("Search() returned non-nil error: %v", err)
		}

		want := []source.S{source.Make(frieren_truffle)}
		if diff := cmp.Diff(
			want,
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
			t.Errorf("Search() mismatch (-want +got):\n%v", diff)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("WithDeleteNode", func(t *testing.T) {
		db := New(context.Background(), config, data)

		if err := db.Delete(context.Background(), source.Make(spyxfamily_mal_header_only).Header()); err != nil {
			t.Errorf("Delete() returned non-nil error: %v", err)
		}

		got, err := db.GetNode(context.Background(), node.Make(&dpb.Node{
			Header: &dpb.NodeHeader{
				Id:   "spyxfamily",
				Type: epb.SourceType_SOURCE_TYPE_SERIES_ANIME,
			}}).Header(),
			option.Remote(true),
		)
		if err != nil {
			t.Errorf("GetNode() returned non-nil error: %v", err)
		}

		want := node.N{}
		if diff := cmp.Diff(
			want,
			got,
			cmp.AllowUnexported(source.S{}, node.N{}),
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
			t.Errorf("GetNode() mismatch (-want +got):\n%v", diff)
		}
	})
	t.Run("NoWithDeleteNode", func(t *testing.T) {
		db := New(context.Background(), config, data)

		if err := db.Delete(context.Background(), source.Make(frieren_truffle).Header()); err != nil {
			t.Errorf("Delete() returned non-nil error: %v", err)
		}

		got, err := db.GetNode(context.Background(), node.Make(&dpb.Node{
			Header: &dpb.NodeHeader{
				Id:   "frieren",
				Type: epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
			}}).Header(),
			option.Remote(true),
		)
		if err != nil {
			t.Errorf("GetNode() returned non-nil error: %v", err)
		}

		want := node.Make(&dpb.Node{
			Header: &dpb.NodeHeader{
				Id:   "frieren",
				Type: epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
			},
		}).WithSources([]source.S{
			source.Make(frieren_mal_with_data),
		})
		if diff := cmp.Diff(
			want,
			got,
			cmp.AllowUnexported(source.S{}, node.N{}),
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
			t.Errorf("GetNode() mismatch (-want +got):\n%v", diff)
		}
	})
}

func TestDeleteNode(t *testing.T) {
	db := New(context.Background(), config, data)

	if err := db.DeleteNode(context.Background(), node.Make(&dpb.Node{
		Header: &dpb.NodeHeader{
			Id:   "frieren",
			Type: epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
		}}).Header(),
	); err != nil {
		t.Errorf("DeleteNode() returned non-nil error: %v", err)
	}

	got_node, err := db.GetNode(context.Background(), node.Make(&dpb.Node{
		Header: &dpb.NodeHeader{
			Id:   "frieren",
			Type: epb.SourceType_SOURCE_TYPE_BOOK_MANGA,
		},
	}).Header(), option.Remote(true))
	if err != nil {
		t.Errorf("GetNode() returned non-nil error: %v", err)
	}

	want_node := node.N{}
	if diff := cmp.Diff(
		want_node,
		got_node,
		cmp.AllowUnexported(source.S{}, node.N{}),
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
		t.Errorf("GetNode() mismatch (-want +got):\n%v", diff)
	}

	got, err := db.Get(context.Background(), source.Make(frieren_truffle).Header(), option.Remote(true))
	if err != nil {
		t.Errorf("Get() returned non-nil error: %v", err)
	}

	want := source.S{}
	if got != want {
		t.Errorf("Get() = %v, want = %v", got, want)
	}
}
