package slice

import (
	"slices"
	"strings"
)

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
