package merge_utils

import (
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

func Distinct[T comparable](ts []T, us []T) []T {
	m := map[T]interface{}{}
	for _, t := range append(append([]T{}, ts...), us...) {
		m[t] = struct{}{}
	}
	res := []T{}
	for k := range m {
		res = append(res, k)
	}
	return res
}

type V[T comparable] struct {
	API epb.API
	V   T
}

// lower compares two API values and returns if the first value should be of
// lower priority than the second (and therefore should be overwritten).
func lower(a epb.API, b epb.API) bool {
	return a >= b
}

func Prioritize[T comparable](a V[T], b V[T]) T {
	var zero T
	if a.V == zero {
		return b.V
	}
	if lower(a.API, b.API) {
		return b.V
	}
	return a.V
}
