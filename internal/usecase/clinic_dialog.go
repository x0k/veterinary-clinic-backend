package usecase

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type DialogRepo interface {
	SaveDialog(ctx context.Context, dialog entity.Dialog) error
}

type DialogPresenter[R any] interface {
	RenderGreeting() (R, error)
	RenderScheduleDialog(dialog entity.Dialog) (R, error)
}

type ClinicDialogUseCase[R any] struct {
	dialogRepo      DialogRepo
	dialogPresenter DialogPresenter[R]
}

func NewClinicDialogUseCase[R any](
	dialogRepo DialogRepo,
	dialogPresenter DialogPresenter[R],
) *ClinicDialogUseCase[R] {
	return &ClinicDialogUseCase[R]{
		dialogRepo:      dialogRepo,
		dialogPresenter: dialogPresenter,
	}
}

func (u *ClinicDialogUseCase[R]) GreetUser(ctx context.Context) (R, error) {
	return u.dialogPresenter.RenderGreeting()
}

func (u *ClinicDialogUseCase[R]) StartScheduleDialog(ctx context.Context, dialog entity.Dialog) (R, error) {
	if err := u.dialogRepo.SaveDialog(ctx, dialog); err != nil {
		return *new(R), err
	}
	return u.dialogPresenter.RenderScheduleDialog(dialog)
}
