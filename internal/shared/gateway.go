package shared

import "context"

type Sender[R any] func(context.Context, R) error
