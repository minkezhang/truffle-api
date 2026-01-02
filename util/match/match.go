package match

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/minkezhang/truffle-api/data/source"
)

func RegExp(match string, s source.S) (float64, error) {
	pattern, err := regexp.Compile(fmt.Sprintf("(?i)%v", match))
	if err != nil {
		return 0, err
	}
	for _, t := range s.Titles() {
		if pattern.MatchString(strings.ToLower(t.Title())) {
			return 1, nil
		}
	}
	return 0, nil

}

func Hamming(match string, s source.S) (float64, error) {
	score := 0.0

	m := metrics.NewHamming()
	m.CaseSensitive = false
	for _, t := range s.Titles() {
		score = math.Max(score, strutil.Similarity(match, t.Title(), m))
	}

	h, _ := RegExp(match, s)

	return (score + h) / 2, nil
}
