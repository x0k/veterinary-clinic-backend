package controller

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase/clinic_make_appointment"
	"gopkg.in/telebot.v3"
)

var ErrUnexpectedMessageType = errors.New("unexpected message type")
var ErrUnknownService = errors.New("unknown service")

type TelegramClinicServiceIdDecoder interface {
	Decode(id string) (entity.ServiceId, bool)
}

func UseTelegramBotRouter(
	ctx context.Context,
	bot *telebot.Bot,
	clinicGreet *usecase.ClinicGreetUseCase[adapters.TelegramTextResponse],
	clinicServices *usecase.ClinicServicesUseCase[adapters.TelegramTextResponse],
	clinicSchedule *usecase.ClinicScheduleUseCase[adapters.TelegramTextResponse],
	clinicMakeAppointmentServicePicker *clinic_make_appointment.ServicePickerUseCase[adapters.TelegramTextResponse],
	clinicServiceIdDecoder TelegramClinicServiceIdDecoder,
) error {
	bot.Handle("/start", func(c telebot.Context) error {
		res, err := clinicGreet.GreetUser(ctx)
		if err != nil {
			return err
		}
		return c.Send(res.Text, res.Options)
	})

	clinicServiceHandler := func(c telebot.Context) error {
		res, err := clinicServices.Services(ctx)
		if err != nil {
			return err
		}
		return c.Send(res.Text, res.Options)
	}
	bot.Handle("/services", clinicServiceHandler)
	bot.Handle(adapters.ClinicServiceBtn, clinicServiceHandler)

	clinicScheduleHandler := func(c telebot.Context) error {
		now := time.Now()
		res, err := clinicSchedule.Schedule(ctx, now, now)
		if err != nil {
			return err
		}
		return c.Send(res.Text, res.Options)
	}
	bot.Handle("/schedule", clinicScheduleHandler)
	bot.Handle(adapters.ClinicScheduleBtn, clinicScheduleHandler)

	bot.Handle(adapters.NextClinicScheduleBtn, func(c telebot.Context) error {
		date, err := time.Parse(time.DateOnly, c.Data())
		if err != nil {
			return err
		}
		res, err := clinicSchedule.Schedule(ctx, time.Now(), date)
		if err != nil {
			return err
		}
		return c.Edit(res.Text, res.Options)
	})

	clinicAppointmentHandler := func(c telebot.Context) error {
		picker, err := clinicMakeAppointmentServicePicker.ServicesPicker(ctx)
		if err != nil {
			return err
		}
		return c.Send(picker.Text, picker.Options)
	}
	bot.Handle("/appointment", clinicAppointmentHandler)
	bot.Handle(adapters.ClinicAppointmentBtn, clinicAppointmentHandler)

	bot.Handle(adapters.ClinicMakeAppointmentServiceCallback, func(c telebot.Context) error {
		serviceId, ok := clinicServiceIdDecoder.Decode(
			c.Callback().Data,
		)
		if !ok {
			return ErrUnknownService
		}
		fmt.Println(serviceId)
		return nil
	})

	return bot.SetCommands([]telebot.Command{
		{
			Text:        "/start",
			Description: "Показать приветствие",
		},
		{
			Text:        "/services",
			Description: "Список услуг",
		},
		{
			Text:        "/schedule",
			Description: "График работы",
		},
		{
			Text:        "/appointment",
			Description: "Записаться на прием",
		},
	})
}

func StartTelegramBotQueryHandler(
	ctx context.Context,
	log *logger.Logger,
	bot *telebot.Bot,
	query <-chan entity.DialogMessage[adapters.TelegramQueryResponse],
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
