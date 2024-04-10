package appointment_telegram_presenter

import (
	"fmt"

	adapters_telegram "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"gopkg.in/telebot.v3"
)

const errorText = "Что-то пошло не так."

type ErrorTextPresenter struct{}

func NewErrorTextPresenter() *ErrorTextPresenter {
	return &ErrorTextPresenter{}
}

func (p *ErrorTextPresenter) RenderError(err error) (adapters_telegram.TextResponse, error) {
	// TODO: Handle domain errors
	return adapters_telegram.TextResponse{
		Text: errorText,
		Options: &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdownV2,
		},
	}, nil
}

type ErrorQueryPresenter struct{}

func NewErrorQueryPresenter() *ErrorQueryPresenter {
	return &ErrorQueryPresenter{}
}

func (p *ErrorQueryPresenter) RenderError(err error) (adapters_telegram.QueryResponse, error) {
	// TODO: Handle domain errors
	return adapters_telegram.QueryResponse{
		Result: &telebot.ArticleResult{
			ResultBase: telebot.ResultBase{
				ID:        fmt.Sprintf("%p", err),
				Type:      "article",
				ParseMode: telebot.ModeMarkdownV2,
			},
			Title: "Ошибка",
			Text:  errorText,
		},
	}, nil
}
