package loader

import (
	"context"
)

type Loader[T any] func(context.Context) (T, error)
