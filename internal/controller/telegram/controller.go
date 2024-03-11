package telegram

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
	"gopkg.in/telebot.v3"
)

func UseRouter(ctx context.Context, bot *telebot.Bot, clinic *usecase.ClinicUseCase[string]) {
	bot.Handle("/start", func(c telebot.Context) error {
		return c.Send("Hello!", &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdownV2,
		})
	})

	bot.Handle("/services", func(c telebot.Context) error {
		services, err := clinic.Services(ctx)
		if err != nil {
			return err
		}
		return c.Send(services, telebot.ModeMarkdownV2)
	})
}
