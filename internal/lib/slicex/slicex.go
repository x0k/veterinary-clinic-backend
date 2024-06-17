package slicex

func Map[T any, U any](mapper func(T) U) func(slice []T) []U {
	return func(slice []T) []U {
		result := make([]U, 0, len(slice))
		for _, item := range slice {
			result = append(result, mapper(item))
		}
		return result
	}
}

func MapEx[A ~[]I, R ~[]O, I any, O any](mapper func(I) (O, error)) func(slice A) (R, error) {
	return func(slice A) (R, error) {
		result := make(R, 0, len(slice))
		for _, item := range slice {
			u, e := mapper(item)
			if e != nil {
				return nil, e
			}
			result = append(result, u)
		}
		return result, nil
	}
}

func MapE[T any, U any](mapper func(T) (U, error)) func(slice []T) ([]U, error) {
	return MapEx[[]T, []U](mapper)
}
