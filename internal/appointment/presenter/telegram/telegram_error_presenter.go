package appointment_telegram_presenter

import (
	"errors"
	"fmt"

	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/telegram"
	"gopkg.in/telebot.v3"
)

const errorText = "Что-то пошло не так."

type ErrorTextPresenter struct{}

func NewErrorTextPresenter() *ErrorTextPresenter {
	return &ErrorTextPresenter{}
}

func (p *ErrorTextPresenter) RenderError(err error) (telegram_adapters.TextResponses, error) {
	if errors.Is(err, appointment_telegram_adapters.ErrUnknownState) {
		return telegram_adapters.TextResponses{{
			Text: "Выбранное действие устарело\\.\nНачните весь процесс заново\\.",
			Options: &telebot.SendOptions{
				ParseMode: telebot.ModeMarkdownV2,
			},
		}}, nil
	}
	// TODO: Handle domain errors
	return telegram_adapters.TextResponses{{
		Text:    errorText,
		Options: &telebot.SendOptions{},
	}}, nil
}

type ErrorQueryPresenter struct{}

func NewErrorQueryPresenter() *ErrorQueryPresenter {
	return &ErrorQueryPresenter{}
}

func (p *ErrorQueryPresenter) RenderError(err error) (telegram_adapters.QueryResponse, error) {
	// TODO: Handle domain errors
	return telegram_adapters.QueryResponse{
		Result: &telebot.ArticleResult{
			ResultBase: telebot.ResultBase{
				ID:   fmt.Sprintf("%p", err),
				Type: "article",
			},
			Title: "Ошибка",
			Text:  errorText,
		},
	}, nil
}

type ErrorCallbackPresenter struct{}

func NewErrorCallbackPresenter() *ErrorCallbackPresenter {
	return &ErrorCallbackPresenter{}
}

func (p *ErrorCallbackPresenter) RenderError(err error) (telegram_adapters.CallbackResponse, error) {
	if errors.Is(err, appointment.ErrInvalidAppointmentStatusForCancel) {
		return telegram_adapters.CallbackResponse{
			Response: &telebot.CallbackResponse{
				Text: "Ваша запись не может быть отменена.",
			},
		}, nil
	}
	return telegram_adapters.CallbackResponse{
		Response: &telebot.CallbackResponse{
			Text: errorText,
		},
	}, nil
}
