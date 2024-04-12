package telegram_adapters

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/module"
	"gopkg.in/telebot.v3"
)

func NewService(name string, bot *telebot.Bot) module.Service {
	return module.NewService(name, func(ctx context.Context) error {
		context.AfterFunc(ctx, func() {
			bot.Stop()
		})
		bot.Start()
		return nil
	})
}
