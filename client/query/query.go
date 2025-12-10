package query

import (
	"math"
	"regexp"
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/hbollon/go-edlib"
	"github.com/minkezhang/bene-api/db/atom"

	epb "github.com/minkezhang/bene-api/proto/go/enums"
)

const (
	reward = 0.5
)

type G struct {
	AtomType epb.Type
	ID       string
}

type O struct {
	AtomTypes []epb.Type
	Title     string
}

type Q struct {
	types map[epb.Type]bool
	title string
}

func New(o O) *Q {
	q := &Q{
		types: map[epb.Type]bool{},
		title: o.Title,
	}
	for _, t := range o.AtomTypes {
		q.types[t] = true
	}
	return q
}

func (q *Q) AtomTypes() []epb.Type {
	var ts []epb.Type
	for t := range q.types {
		ts = append(ts, t)
	}
	return ts
}

func (q *Q) IsSupportedType(v epb.Type) bool { return q.types[v] }
func (q *Q) Title() string                   { return q.title }

func RegExp(q *Q, a *atom.A) (float64, error) {
	pattern, err := regexp.Compile(strings.ToLower(q.Title()))
	if err != nil {
		return 0, err
	}

	if !q.IsSupportedType(a.AtomType()) {
		return 0, nil
	}

	for _, t := range a.Titles() {
		if pattern.MatchString(strings.ToLower(t.Title)) {
			return 1, nil
		}
	}
	return 0, nil

}

func Jaccard(q *Q, a *atom.A) (float64, error) {
	if !q.IsSupportedType(a.AtomType()) {
		return 0, nil
	}

	score := float64(0)
	j := metrics.NewJaccard()
	j.CaseSensitive = false
	for _, t := range a.Titles() {
		score = math.Max(score, strutil.Similarity(q.Title(), t.Title, j))
	}

	if h, _ := RegExp(q, a); h > 0 {
		score = math.Min(score+reward, 1)
	}

	return score, nil
}

func Hamming(q *Q, a *atom.A) (float64, error) {
	if !q.IsSupportedType(a.AtomType()) {
		return 0, nil
	}

	score := float64(0)
	h := metrics.NewHamming()
	h.CaseSensitive = false
	for _, t := range a.Titles() {
		score = math.Max(score, strutil.Similarity(q.Title(), t.Title, h))
	}

	if h, _ := RegExp(q, a); h > 0 {
		score = math.Min(score+reward, 1)
	}

	return score, nil
}

func LCS(q *Q, a *atom.A) (float64, error) {
	if !q.IsSupportedType(a.AtomType()) {
		return 0, nil
	}

	score := float64(0)
	for _, t := range a.Titles() {
		score = math.Max(
			score,
			(2*float64(edlib.LCS(q.Title(), t.Title)))/float64(len(q.Title())+len(t.Title)))
	}

	if h, _ := RegExp(q, a); h > 0 {
		score = math.Min(score+reward, 1)
	}

	return score, nil
}
