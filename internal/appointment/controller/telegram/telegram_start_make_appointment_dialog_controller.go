package appointment_telegram_controller

import (
	"context"
	"fmt"
	"strconv"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/telegram"
	appointment_telegram_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"gopkg.in/telebot.v3"
)

func NewStartMakeAppointmentDialog(
	bot *telebot.Bot,
	tgUserIdLoader adapters.StatePopper[entity.TelegramUserId],
	startMakeAppointmentDialogUseCase *appointment_telegram_use_case.StartMakeAppointmentDialogUseCase[telegram_adapters.TextResponses],
	registerCustomerUseCase *appointment_telegram_use_case.RegisterCustomerUseCase[telegram_adapters.TextResponses],
	errorPresenter appointment.ErrorPresenter[telegram_adapters.TextResponses],
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
			return telegram_adapters.Send(c, res)
		}
		bot.Handle("/appointment", startMakeAppointmentHandler)
		bot.Handle(appointment_telegram_adapters.StartMakeAppointmentDialogBtn, startMakeAppointmentHandler)
		bot.Handle(telebot.OnContact, func(c telebot.Context) error {
			cnt := c.Message().Contact
			if cnt == nil {
				return nil
			}
			tgUserId, ok := tgUserIdLoader.Pop(adapters.StateId(strconv.FormatInt(cnt.UserID, 10)))
			if !ok {
				return fmt.Errorf("user id %d is not found in registration queue", cnt.UserID)
			}
			res, err := registerCustomerUseCase.RegisterCustomer(
				ctx,
				tgUserId,
				c.Sender().Username,
				cnt.FirstName,
				cnt.LastName,
				cnt.PhoneNumber,
			)
			if err != nil {
				return err
			}
			return telegram_adapters.Send(c, res)
		})
		bot.Handle(appointment_telegram_adapters.CancelRegisterTelegramCustomerBtn, func(c telebot.Context) error {
			return c.Send("Регистрация отменена", &telebot.SendOptions{
				ReplyMarkup: &telebot.ReplyMarkup{RemoveKeyboard: true},
			})
		})
		return nil
	}
}
