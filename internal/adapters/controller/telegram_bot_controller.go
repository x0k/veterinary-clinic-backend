package controller

import (
	"context"
	"errors"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase/make_appointment"
	"gopkg.in/telebot.v3"
)

var ErrUnexpectedMessageType = errors.New("unexpected message type")
var ErrUnknownService = errors.New("unknown service")
var ErrUnknownDatePickerState = errors.New("unknown date picker state")

func UseTelegramBotRouter(
	ctx context.Context,
	bot *telebot.Bot,
	greet *usecase.GreetUseCase[adapters.TelegramTextResponse],
	services *usecase.ServicesUseCase[adapters.TelegramTextResponse],
	schedule *usecase.ScheduleUseCase[adapters.TelegramTextResponse],
	makeAppointmentServicePicker *make_appointment.ServicePickerUseCase[adapters.TelegramTextResponse],
	serviceIdLoader adapters.StateLoader[entity.ServiceId],
	makeAppointmentDatePicker *make_appointment.DatePickerUseCase[adapters.TelegramTextResponse],
	datePickerStateLoader adapters.StateLoader[adapters.TelegramDatePickerState],
	makeAppointmentTimePicker *make_appointment.TimeSlotPickerUseCase[adapters.TelegramTextResponse],
	makeAppointmentConfirmation *make_appointment.AppointmentConfirmationUseCase[adapters.TelegramTextResponse],
	makeAppointment *make_appointment.MakeAppointmentUseCase[adapters.TelegramTextResponse],
	cancelAppointment *usecase.CancelAppointmentUseCase[adapters.TelegramCallbackResponse],
) error {
	bot.Handle("/start", func(c telebot.Context) error {
		res, err := greet.GreetUser(ctx)
		if err != nil {
			return err
		}
		return c.Send(res.Text, res.Options)
	})

	serviceHandler := func(c telebot.Context) error {
		res, err := services.Services(ctx)
		if err != nil {
			return err
		}
		return c.Send(res.Text, res.Options)
	}
	bot.Handle("/services", serviceHandler)
	bot.Handle(adapters.ServicesBtn, serviceHandler)

	scheduleHandler := func(c telebot.Context) error {
		now := time.Now()
		res, err := schedule.Schedule(ctx, now, now)
		if err != nil {
			return err
		}
		return c.Send(res.Text, res.Options)
	}
	bot.Handle("/schedule", scheduleHandler)
	bot.Handle(adapters.ScheduleBtn, scheduleHandler)

	bot.Handle(adapters.NextScheduleBtn, func(c telebot.Context) error {
		date, err := time.Parse(time.DateOnly, c.Data())
		if err != nil {
			return err
		}
		res, err := schedule.Schedule(ctx, time.Now(), date)
		if err != nil {
			return err
		}
		return c.Edit(res.Text, res.Options)
	})

	makeAppointmentServicePickerHandler := func(c telebot.Context) error {
		servicePicker, err := makeAppointmentServicePicker.ServicesPicker(
			ctx,
			entity.TelegramUserIdToUserId(entity.TelegramUserId(c.Sender().ID)),
		)
		if err != nil {
			return err
		}
		return c.Send(servicePicker.Text, servicePicker.Options)
	}
	bot.Handle("/appointment", makeAppointmentServicePickerHandler)
	bot.Handle(adapters.AppointmentBtn, makeAppointmentServicePickerHandler)

	bot.Handle(adapters.MakeAppointmentServiceCallback, func(c telebot.Context) error {
		serviceId, ok := serviceIdLoader.Load(
			adapters.StateId(c.Callback().Data),
		)
		if !ok {
			return ErrUnknownService
		}
		now := time.Now()
		datePicker, err := makeAppointmentDatePicker.DatePicker(
			ctx,
			serviceId,
			now,
			now,
		)
		if err != nil {
			return err
		}
		return c.Edit(datePicker.Text, datePicker.Options)
	})

	makeAppointmentNextDatePickerHandler := func(c telebot.Context) error {
		state, ok := datePickerStateLoader.Load(
			adapters.StateId(c.Callback().Data),
		)
		if !ok {
			return ErrUnknownDatePickerState
		}
		datePicker, err := makeAppointmentDatePicker.DatePicker(
			ctx,
			state.ServiceId,
			time.Now(),
			state.Date,
		)
		if err != nil {
			return err
		}
		return c.Edit(datePicker.Text, datePicker.Options)
	}
	bot.Handle(adapters.NextMakeAppointmentDateBtn, makeAppointmentNextDatePickerHandler)

	bot.Handle(adapters.CancelMakeAppointmentDateBtn, func(c telebot.Context) error {
		servicePicker, err := makeAppointmentServicePicker.ServicesPicker(
			ctx,
			entity.TelegramUserIdToUserId(entity.TelegramUserId(c.Sender().ID)),
		)
		if err != nil {
			return err
		}
		return c.Edit(servicePicker.Text, servicePicker.Options)
	})

	makeAppointmentTimePickerHandler := func(c telebot.Context) error {
		state, ok := datePickerStateLoader.Load(
			adapters.StateId(c.Callback().Data),
		)
		if !ok {
			return ErrUnknownDatePickerState
		}
		timePicker, err := makeAppointmentTimePicker.TimePicker(
			ctx,
			state.ServiceId,
			time.Now(),
			state.Date,
		)
		if err != nil {
			return err
		}
		return c.Edit(timePicker.Text, timePicker.Options)
	}
	bot.Handle(adapters.SelectMakeAppointmentDateBtn, makeAppointmentTimePickerHandler)

	bot.Handle(adapters.MakeAppointmentTimeCallback, func(c telebot.Context) error {
		state, ok := datePickerStateLoader.Load(
			adapters.StateId(c.Callback().Data),
		)
		if !ok {
			return ErrUnknownDatePickerState
		}
		confirmation, err := makeAppointmentConfirmation.Confirmation(
			ctx,
			state.ServiceId,
			state.Date,
		)
		if err != nil {
			return err
		}
		return c.Edit(confirmation.Text, confirmation.Options)
	})

	bot.Handle(adapters.CancelMakeAppointmentTimeBtn, makeAppointmentNextDatePickerHandler)

	bot.Handle(adapters.ConfirmMakeAppointmentBtn, func(c telebot.Context) error {
		state, ok := datePickerStateLoader.Load(
			adapters.StateId(c.Callback().Data),
		)
		if !ok {
			return ErrUnknownDatePickerState
		}
		res, err := makeAppointment.Make(
			ctx,
			entity.NewTelegramUser(
				entity.TelegramUserId(c.Sender().ID),
				entity.TelegramUsername(c.Sender().Username),
				entity.TelegramFirstName(c.Sender().FirstName),
				entity.TelegramLastName(c.Sender().LastName),
			),
			state.ServiceId,
			state.Date,
		)
		if err != nil {
			return err
		}
		return c.Edit(res.Text, res.Options)
	})

	bot.Handle(adapters.CancelConfirmationAppointmentBtn, makeAppointmentTimePickerHandler)

	bot.Handle(adapters.CancelAppointmentBtn, func(c telebot.Context) error {
		recordId := entity.RecordId(c.Callback().Data)
		res, err := cancelAppointment.Cancel(ctx, recordId)
		if err != nil {
			return err
		}
		if err := c.Respond(res.Response); err != nil {
			return err
		}
		return c.Delete()
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
			Description: "Запись на прием",
		},
	})
}
