//go:build modern

package mapx

import "maps"

// Clone returns a copy of m.  This is a shallow clone:
// the new keys and values are set using ordinary assignment.
func Clone[M ~map[K]V, K comparable, V any](m M) M {
	return maps.Clone(m)
}

func Copy[M ~map[K]V, K comparable, V any](dst M, src M) {
	maps.Copy(dst, src)
}
