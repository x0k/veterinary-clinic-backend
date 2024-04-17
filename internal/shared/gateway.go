package shared

import "context"

type Sender[T any] func(context.Context, T) error

type Saver[T any] func(context.Context, T) error

type Loader[T any] func(context.Context) (T, error)

type QueriedLoader[Q any, T any] func(context.Context, Q) (T, error)
