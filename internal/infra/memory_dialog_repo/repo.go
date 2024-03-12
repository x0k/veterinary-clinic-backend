package memory_dialog_repo

import (
	"context"
	"sync"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type MemoryDialogRepo struct {
	mu      sync.RWMutex
	dialogs map[entity.DialogId]entity.UserId
}

func New() *MemoryDialogRepo {
	return &MemoryDialogRepo{
		dialogs: make(map[entity.DialogId]entity.UserId),
	}
}

func (r *MemoryDialogRepo) SaveDialog(ctx context.Context, dialog entity.Dialog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.dialogs[dialog.Id] = dialog.UserId
	return nil
}
