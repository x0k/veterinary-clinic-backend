package usecase

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type DialogRepo interface {
	SaveDialog(ctx context.Context, dialog entity.Dialog) error
}

type DialogPresenter[ScheduleDialog any] interface {
	RenderScheduleDialog(dialog entity.Dialog) (ScheduleDialog, error)
}

type ClinicDialogUseCase[ScheduleDialog any] struct {
	dialogRepo      DialogRepo
	dialogPresenter DialogPresenter[ScheduleDialog]
}

func NewClinicDialogUseCase[ScheduleDialog any](
	dialogRepo DialogRepo,
	dialogPresenter DialogPresenter[ScheduleDialog],
) *ClinicDialogUseCase[ScheduleDialog] {
	return &ClinicDialogUseCase[ScheduleDialog]{
		dialogRepo:      dialogRepo,
		dialogPresenter: dialogPresenter,
	}
}

func (u *ClinicDialogUseCase[D]) StartScheduleDialog(ctx context.Context, dialog entity.Dialog) (D, error) {
	if err := u.dialogRepo.SaveDialog(ctx, dialog); err != nil {
		return *new(D), err
	}
	return u.dialogPresenter.RenderScheduleDialog(dialog)
}
