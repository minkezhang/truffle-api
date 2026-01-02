package util

import (
	"slices"
	"testing"
)

func TestDeduplicateStrings(t *testing.T) {
	configs := []struct {
		name string
		vs   []string
		cmp  func(a, b string) int
		eq   func(a, b string) bool
		want []string
	}{
		{
			name: "Trivial",
			vs:   []string{},
			want: []string{},
		},
		{
			name: "Trivial/SingleElement",
			vs:   []string{"foo"},
			want: []string{"foo"},
		},
		{
			name: "Deduplicate",
			vs:   []string{"foo", "bar", "foo"},
			want: []string{"bar", "foo"},
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := DeduplicateStrings(c.vs); slices.Compare(got, c.want) != 0 {
				t.Errorf("DeduplicateStrings() = %v, want = %v", got, c.want)
			}
		})
	}
}
