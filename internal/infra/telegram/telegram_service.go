package adapters_telegram

import (
	"context"
	"log/slog"

	"github.com/x0k/veterinary-clinic-backend/internal/infra"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

func NewService(bot *telebot.Bot, log *logger.Logger) infra.Starter {
	return infra.Starter(func(ctx context.Context) error {
		bot.Use(
			middleware.Logger(slog.NewLogLogger(log.Logger.Handler(), slog.LevelInfo)),
			middleware.AutoRespond(),
			middleware.Recover(),
		)
		context.AfterFunc(ctx, func() {
			bot.Stop()
		})
		bot.Start()
		return nil
	})
}
