//go:build !modern

package slicex

import "sort"

func Clone[T any](s []T) []T {
	cloned := make([]T, len(s))
	copy(cloned, s)
	return cloned
}

func SortFunc[T any](arr []T, cmp func(a, b T) int) {
	sort.Slice(arr, func(i, j int) bool {
		return cmp(arr[i], arr[j]) < 0
	})
}

func Index[T comparable](arr []T, item T) int {
	for i, v := range arr {
		if v == item {
			return i
		}
	}
	return -1
}

func Delete[T any](arr []T, i, j int) []T {
	return append(arr[:i], arr[j:]...)
}
