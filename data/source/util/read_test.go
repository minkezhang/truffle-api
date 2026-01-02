package util

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	dpb "github.com/minkezhang/truffle-api/proto/go/data"
)

func TestDeduplicateTitles(t *testing.T) {
	configs := []struct {
		name   string
		titles []*dpb.Title
		want   []*dpb.Title
	}{
		{
			name:   "Trivial",
			titles: []*dpb.Title{},
			want:   []*dpb.Title{},
		},
		{
			name: "Trivial/SingleElement",
			titles: []*dpb.Title{
				&dpb.Title{Title: "foo"},
			},
			want: []*dpb.Title{
				&dpb.Title{Title: "foo"},
			},
		},
		{
			name: "Simple/Title",
			titles: []*dpb.Title{
				&dpb.Title{Title: "foo"},
				&dpb.Title{Title: "bar"},
			},
			want: []*dpb.Title{
				&dpb.Title{Title: "bar"},
				&dpb.Title{Title: "foo"},
			},
		},
		{
			name: "Simple/Localization",
			titles: []*dpb.Title{
				&dpb.Title{Title: "foo"},
				&dpb.Title{Title: "bar", Localization: "en"},
			},
			want: []*dpb.Title{
				&dpb.Title{Title: "bar", Localization: "en"},
				&dpb.Title{Title: "foo"},
			},
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := DeduplicateTitles(c.titles)
			if diff := cmp.Diff(c.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("DeduplicateTitles() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
