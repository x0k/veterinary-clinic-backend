package appointment_telegram_presenter

import (
	"fmt"

	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"gopkg.in/telebot.v3"
)

const errorText = "Что-то пошло не так."

type ErrorTextPresenter struct{}

func NewErrorTextPresenter() *ErrorTextPresenter {
	return &ErrorTextPresenter{}
}

func (p *ErrorTextPresenter) RenderError(err error) (telegram_adapters.TextResponses, error) {
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
