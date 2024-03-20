package presenter

import (
	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type TelegramChangePresenter struct{}

func NewTelegramChangePresenter() *TelegramChangePresenter {
	return &TelegramChangePresenter{}
}

func (p *TelegramChangePresenter) RenderChange(change entity.RecordChange) (adapters.TelegramTextResponse, error) {

}
