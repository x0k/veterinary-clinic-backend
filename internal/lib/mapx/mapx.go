//go:build !modern

package mapx

// Clone returns a copy of m.  This is a shallow clone:
// the new keys and values are set using ordinary assignment.
func Clone[M ~map[K]V, K comparable, V any](m M) M {
	cloned := make(M, len(m))
	for k, v := range m {
		cloned[k] = v
	}
	return cloned
}

func Copy[M ~map[K]V, K comparable, V any](dst M, src M) {
	for k, v := range src {
		dst[k] = v
	}
}
