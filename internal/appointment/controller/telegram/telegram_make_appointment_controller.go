package appointment_telegram_controller

import (
	"context"

	adapters_telegram "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	appointment_telegram_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"gopkg.in/telebot.v3"
)

func NewMakeAppointment(
	bot *telebot.Bot,
	startMakeAppointmentDialogUseCase *appointment_telegram_use_case.StartMakeAppointmentDialogUseCase[adapters_telegram.TextResponse],
) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		startMakeAppointmentHandler := func(c telebot.Context) error {
			res, err := startMakeAppointmentDialogUseCase.StartMakeAppointmentDialog(
				ctx,
				entity.NewTelegramUserId(c.Sender().ID),
			)
			if err != nil {
				return err
			}
			return c.Send(res.Text, res.Options)
		}
		bot.Handle("/appointment", startMakeAppointmentHandler)
		bot.Handle(adapters_telegram.StartMakeAppointmentDialogBtn, startMakeAppointmentHandler)
		return nil
	}
}
