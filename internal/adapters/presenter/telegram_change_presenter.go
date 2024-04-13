package presenter

import (
	"fmt"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
	"gopkg.in/telebot.v3"
)

type TelegramChangePresenter struct{}

func NewTelegramChangePresenter() *TelegramChangePresenter {
	return &TelegramChangePresenter{}
}

func (p *TelegramChangePresenter) RenderChange(change shared.RecordChange) (adapters.TelegramTextResponse, error) {
	switch change.Type {
	case shared.RecordCreated:
		return adapters.TelegramTextResponse{
			Text: fmt.Sprintf(
				"Новая запись: %s, %s",
				adapters.EscapeTelegramMarkdownString(change.Record.Service.Title),
				adapters.EscapeTelegramMarkdownString(
					shared.DateTimeToGoTime(change.Record.DateTimePeriod.Start).
						Format("02.01.06 15:04"),
				),
			),
			Options: &telebot.SendOptions{
				ParseMode: telebot.ModeMarkdownV2,
			},
		}, nil
	case shared.RecordStatusChanged:
		statusName, err := shared.RecordStatusName(change.Record.Status)
		if err != nil {
			return adapters.TelegramTextResponse{}, err
		}
		return adapters.TelegramTextResponse{
			Text: fmt.Sprintf("Статус записи изменен: %s", statusName),
			Options: &telebot.SendOptions{
				ParseMode: telebot.ModeMarkdownV2,
			},
		}, nil
	case shared.RecordDateTimeChanged:
		return adapters.TelegramTextResponse{
			Text: fmt.Sprintf(
				"Время записи изменено: %s",
				adapters.EscapeTelegramMarkdownString(
					shared.DateTimeToGoTime(change.Record.DateTimePeriod.Start).
						Format("02.01.06 15:04"),
				),
			),
			Options: &telebot.SendOptions{
				ParseMode: telebot.ModeMarkdownV2,
			},
		}, nil
	case shared.RecordRemoved:
		return adapters.TelegramTextResponse{
			Text: fmt.Sprintf(
				"Запись удалена: %s, %s",
				adapters.EscapeTelegramMarkdownString(
					change.Record.Service.Title,
				),
				adapters.EscapeTelegramMarkdownString(
					shared.DateTimeToGoTime(change.Record.DateTimePeriod.Start).
						Format("02.01.06 15:04"),
				),
			),
			Options: &telebot.SendOptions{
				ParseMode: telebot.ModeMarkdownV2,
			},
		}, nil
	default:
		return adapters.TelegramTextResponse{}, fmt.Errorf("unexpected change type: %d", change.Type)
	}
}
