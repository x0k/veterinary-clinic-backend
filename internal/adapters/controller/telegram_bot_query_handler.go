package controller

import (
	"context"
	"log/slog"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
	"gopkg.in/telebot.v3"
)

func StartTelegramBotQueryHandler(
	ctx context.Context,
	log *logger.Logger,
	bot *telebot.Bot,
	query <-chan shared.DialogMessage[adapters.TelegramQueryResponse],
) {
	l := log.With(slog.String("component", "adapters.controller.RunTelegramBotQueryHandler"))
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-query:
			l.Debug(ctx, "received query", slog.String("query_id", string(msg.DialogId)))
			_, err := bot.AnswerWebApp(
				&telebot.Query{
					ID: string(msg.DialogId),
				},
				msg.Message.Result,
			)
			if err != nil {
				l.Error(ctx, "failed to answer query", slog.String("query_id", string(msg.DialogId)), sl.Err(err))
			}
		}
	}
}
