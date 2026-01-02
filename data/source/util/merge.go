package util

import (
	"slices"
	"strings"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
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

func Apply[T any, U any](vs []T, f func(v T) U) []U {
	results := []U{}
	for _, v := range vs {
		results = append(results, f(v))
	}
	return results
}

func DeduplicateFunc[T any](
	vs []T,
	cmp func(u, v T) int,
	eq func(u, v T) bool,
) []T {
	slices.SortStableFunc(vs, cmp)
	return slices.CompactFunc(vs, eq)
}

func DeduplicateStrings(vs []string) []string {
	return DeduplicateFunc(
		Apply(vs, func(v string) string {
			return strings.TrimSpace(v)
		}),
		strings.Compare,
		func(a, b string) bool {
			return strings.ToLower(a) == strings.ToLower(b)
		},
	)
}

func Filter[T any](vs []T, f func(v T) bool) []T {
	results := []T{}
	for _, v := range vs {
		if !f(v) {
			results = append(results, v)
		}
	}
	return results
}
