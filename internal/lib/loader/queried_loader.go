package loader

import "context"

type QueriedLoader[Q any, T any] func(context.Context, Q) (T, error)

type KeyedLoader[K comparable, T any] QueriedLoader[K, T]
