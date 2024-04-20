//go:build js && wasm

package mapx

func Clone[M ~map[K]V, K comparable, V any](m M) M {
	clone := make(M, len(m))
	for k, v := range m {
		clone[k] = v
	}
	return clone
}
