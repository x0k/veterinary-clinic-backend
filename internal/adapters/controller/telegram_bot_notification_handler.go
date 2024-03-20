package controller

import (
	"context"
	"log/slog"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"gopkg.in/telebot.v3"
)

func StartTelegramBotNotificationHandler(
	ctx context.Context,
	log *logger.Logger,
	bot *telebot.Bot,
	notifications <-chan entity.NotificationMessage[adapters.TelegramTextResponse],
) {
	l := log.With(slog.String("component", "adapters.controller.StartTelegramBotNotificationHandler"))
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-notifications:
			tgId, err := entity.UserIdToTelegramUserId(msg.UserId)
			if err != nil {
				l.Error(ctx, "failed to convert user id to telegram user id", sl.Err(err))
				continue
			}
			if _, err = bot.Send(
				&telebot.User{ID: int64(tgId)},
				msg.Message.Text,
				msg.Message.Options,
			); err != nil {
				l.Error(ctx, "failed to send message", sl.Err(err))
			}
		}
	}
}
