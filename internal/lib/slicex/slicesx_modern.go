//go:build modern

package slicex

import "slices"

func Clone[T any](s []T) []T {
	return slices.Clone(s)
}

func SortFunc[T any](arr []T, cmp func(a, b T) int) {
	slices.SortFunc(arr, cmp)
}

func Index[T comparable](arr []T, item T) int {
	return slices.Index(arr, item)
}

func Delete[T any](arr []T, i, j int) []T {
	return slices.Delete(arr, i, j)
}
