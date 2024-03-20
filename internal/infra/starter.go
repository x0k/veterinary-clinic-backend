package infra

import "context"

type Starter func(ctx context.Context) error

func (s Starter) Start(ctx context.Context) error {
	return s(ctx)
}
