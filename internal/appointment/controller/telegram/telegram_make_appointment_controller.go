package appointment_telegram_controller

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/telegram"
	appointment_telegram_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"gopkg.in/telebot.v3"
)

func NewMakeAppointment(
	bot *telebot.Bot,
	startMakeAppointmentDialogUseCase *appointment_telegram_use_case.StartMakeAppointmentDialogUseCase[telegram_adapters.TextResponses],
	appointmentDatePickerUseCase *appointment_telegram_use_case.AppointmentDatePickerUseCase[telegram_adapters.TextResponses],
	appointmentTimePickerUseCase *appointment_telegram_use_case.AppointmentTimePickerUseCase[telegram_adapters.TextResponses],
	errorPresenter appointment.ErrorPresenter[telegram_adapters.TextResponses],
	serviceIdLoader adapters.StateLoader[appointment.ServiceId],
	appointmentStateLoader adapters.StateLoader[appointment_telegram_adapters.AppointmentSate],
) func(context.Context) error {
	return func(ctx context.Context) error {
		bot.Handle(appointment_telegram_adapters.MakeAppointmentServiceCallback, func(c telebot.Context) error {
			serviceId, ok := serviceIdLoader.Load(adapters.NewStateId(c.Callback().Data))
			if !ok {
				res, err := errorPresenter.RenderError(appointment_telegram_adapters.ErrUnknownState)
				if err != nil {
					return err
				}
				return res.Send(c)
			}
			now := time.Now()
			datePicker, err := appointmentDatePickerUseCase.DatePicker(ctx, serviceId, now, now)
			if err != nil {
				return err
			}
			return datePicker.Edit(c)
		})

		appointmentNextDatePickerHandler := func(c telebot.Context) error {
			state, ok := appointmentStateLoader.Load(
				adapters.NewStateId(c.Callback().Data),
			)
			if !ok {
				res, err := errorPresenter.RenderError(appointment_telegram_adapters.ErrUnknownState)
				if err != nil {
					return err
				}
				return res.Send(c)
			}
			datePicker, err := appointmentDatePickerUseCase.DatePicker(
				ctx,
				state.ServiceId,
				time.Now(),
				state.Date,
			)
			if err != nil {
				return err
			}
			return datePicker.Edit(c)
		}
		bot.Handle(appointment_telegram_adapters.NextMakeAppointmentDateBtn, appointmentNextDatePickerHandler)

		bot.Handle(appointment_telegram_adapters.CancelMakeAppointmentDateBtn, func(c telebot.Context) error {
			res, err := startMakeAppointmentDialogUseCase.StartMakeAppointmentDialog(
				ctx,
				entity.NewTelegramUserId(c.Sender().ID),
			)
			if err != nil {
				return err
			}
			return res.Edit(c)
		})

		appointmentTimePickerHandler := func(c telebot.Context) error {
			state, ok := appointmentStateLoader.Load(
				adapters.NewStateId(c.Callback().Data),
			)
			if !ok {
				res, err := errorPresenter.RenderError(appointment_telegram_adapters.ErrUnknownState)
				if err != nil {
					return err
				}
				return res.Send(c)
			}
			timePicker, err := appointmentTimePickerUseCase.TimePicker(
				ctx,
				state.ServiceId,
				time.Now(),
				state.Date,
			)
			if err != nil {
				return err
			}
			return timePicker.Edit(c)
		}
		bot.Handle(appointment_telegram_adapters.SelectMakeAppointmentDateBtn, appointmentTimePickerHandler)

		// TODO: Handler MakeAppointmentTimeCallback

		bot.Handle(appointment_telegram_adapters.CancelMakeAppointmentTimeBtn, appointmentNextDatePickerHandler)

		return nil
	}
}
