package infra

import "github.com/x0k/veterinary-clinic-backend/internal/lib/containers"

type MemoryExpirableStateContainer[S any] struct {
	keys   *containers.ExpirationQueue[string]
	values map[string]S
}

func NewMemoryExpirableStateContainer[S any]() *MemoryExpirableStateContainer[S] {
	return &MemoryExpirableStateContainer[S]{
		keys:   containers.NewExpirationQueue[string](),
		values: make(map[string]S),
	}
}
