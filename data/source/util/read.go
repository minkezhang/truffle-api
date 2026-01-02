package util

import (
	"cmp"
	"strings"

	dpb "github.com/minkezhang/truffle-api/proto/go/data"
)

var (
	localization_priority = map[string]int{
		"en": 0,
		"":   1,
		"ja": 2,
	}
)

func DeduplicateTitles(titles []*dpb.Title) []*dpb.Title {
	return DeduplicateFunc(
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
