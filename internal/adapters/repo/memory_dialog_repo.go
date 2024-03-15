package repo

import (
	"context"
	"sync"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type Memory struct {
	mu      sync.RWMutex
	dialogs map[entity.DialogId]entity.UserId
}

func NewMemoryDialog() *Memory {
	return &Memory{
		dialogs: make(map[entity.DialogId]entity.UserId),
	}
}

func (r *Memory) SaveDialog(ctx context.Context, dialog entity.Dialog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.dialogs[dialog.Id] = dialog.UserId
	return nil
}
