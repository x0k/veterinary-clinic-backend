package appointment_telegram_controller

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/telegram"
	appointment_telegram_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case/telegram"
	"gopkg.in/telebot.v3"
)

func NewMakeAppointment(
	bot *telebot.Bot,
	appointmentDatePickerUseCase *appointment_telegram_use_case.AppointmentDatePickerUseCase[telegram_adapters.TextResponses],
	errorPresenter appointment.ErrorPresenter[telegram_adapters.TextResponses],
	serviceIdLoader adapters.StateLoader[appointment.ServiceId],
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
		return nil
	}
}
