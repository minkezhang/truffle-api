package merge

import (
	"cmp"
	"strings"

	"github.com/minkezhang/truffle-api/util/slice"

	dpb "github.com/minkezhang/truffle-api/proto/go/data"
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

var (
	localization_priority = map[string]int{
		"en": 0,
		"":   1,
		"ja": 2,
	}
)

// is_higher_priority returns if the source API u is of higher priority than v.
func is_higher_priority(u, v epb.SourceAPI) bool { return u < v }

func Prioritize[T comparable](a epb.SourceAPI, u T, b epb.SourceAPI, v T) T {
	var zero T
	if u == zero || is_higher_priority(b, a) {
		return v
	}
	return u
}

func DeduplicateStrings(vs []string) []string {
	return slice.DeduplicateFunc(
		slice.Apply(vs, func(v string) string {
			return strings.TrimSpace(v)
		}),
		strings.Compare,
		func(a, b string) bool {
			return strings.ToLower(a) == strings.ToLower(b)
		},
	)
}

func DeduplicateTitles(titles []*dpb.Title) []*dpb.Title {
	return slice.DeduplicateFunc(
		titles,
		func(u, v *dpb.Title) int {
			if u.GetLocalization() == v.GetLocalization() {
				return strings.Compare(u.GetTitle(), v.GetTitle())
			}
			p, ok := localization_priority[u.GetLocalization()]
			if !ok {
				p = int(^uint(0) >> 1)
			}
			q, ok := localization_priority[v.GetLocalization()]
			if !ok {
				q = int(^uint(0) >> 1)
			}
			return cmp.Compare(p, q)
		},
		func(u, v *dpb.Title) bool {
			return u.GetTitle() == v.GetTitle()
		},
	)
}
