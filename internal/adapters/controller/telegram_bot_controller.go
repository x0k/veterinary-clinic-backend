package controller

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
	"gopkg.in/telebot.v3"
)

var ErrUnexpectedMessageType = errors.New("unexpected message type")

func UseTelegramBotRouter(
	ctx context.Context,
	wg *sync.WaitGroup,
	log *logger.Logger,
	bot *telebot.Bot,
	clinic *usecase.ClinicUseCase[adapters.TelegramResponse],
	clinicDialog *usecase.ClinicDialogUseCase[adapters.TelegramResponse],
) {
	l := log.With(slog.String("component", "controller.telegram.UseRouter"))

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-clinicDialog.Messages():
				queryResponse, ok := msg.Message.(adapters.TelegramQueryResponse)
				if !ok {
					l.Error(ctx, "unexpected message type", slog.Int("type", int(msg.Message.Type())))
					continue
				}
				if _, err := bot.AnswerWebApp(
					&telebot.Query{
						ID: string(msg.DialogId),
					},
					queryResponse.Result,
				); err != nil {
					l.Error(ctx, "failed to answer query", sl.Err(err))
				}
			}
		}
	}()

	send := func(c telebot.Context, response adapters.TelegramResponse) error {
		msg, ok := response.(adapters.TelegramTextResponse)
		if !ok {
			return ErrUnexpectedMessageType
		}
		return c.Send(msg.Text, msg.Options)
	}

	edit := func(c telebot.Context, response adapters.TelegramResponse) error {
		msg, ok := response.(adapters.TelegramTextResponse)
		if !ok {
			return ErrUnexpectedMessageType
		}
		return c.Edit(msg.Text, msg.Options)
	}

	bot.Handle("/start", func(c telebot.Context) error {
		res, err := clinicDialog.GreetUser(ctx)
		if err != nil {
			return err
		}
		return send(c, res)
	})

	bot.Handle("/services", func(c telebot.Context) error {
		res, err := clinic.Services(ctx)
		if err != nil {
			return err
		}
		return send(c, res)
	})

	bot.Handle("/s", func(c telebot.Context) error {
		res, err := clinicDialog.Schedule(ctx, time.Now())
		if err != nil {
			return err
		}
		return send(c, res)
	})

	bot.Handle(adapters.NextScheduleBtn, func(c telebot.Context) error {
		date, err := time.Parse(time.DateOnly, c.Data())
		if err != nil {
			return err
		}
		res, err := clinicDialog.Schedule(ctx, date)
		if err != nil {
			return err
		}
		return edit(c, res)
	})
}
