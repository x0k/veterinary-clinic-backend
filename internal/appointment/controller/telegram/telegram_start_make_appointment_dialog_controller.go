package appointment_telegram_controller

import (
	"context"
	"fmt"
	"strconv"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	adapters_telegram "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	appointment_telegram_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"gopkg.in/telebot.v3"
)

func NewStartMakeAppointmentDialog(
	bot *telebot.Bot,
	tgUserIdLoader adapters.StatePopper[entity.TelegramUserId],
	startMakeAppointmentDialogUseCase *appointment_telegram_use_case.StartMakeAppointmentDialogUseCase[adapters_telegram.TextResponse],
	registerCustomerUseCase *appointment_telegram_use_case.RegisterCustomerUseCase[adapters_telegram.TextResponse],
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
		bot.Handle(telebot.OnContact, func(c telebot.Context) error {
			cnt := c.Message().Contact
			if cnt == nil {
				return nil
			}
			tgUserId, ok := tgUserIdLoader.Pop(adapters.StateId(strconv.FormatInt(cnt.UserID, 10)))
			if !ok {
				return fmt.Errorf("user id %d is not found in registration queue", cnt.UserID)
			}
			isRegistered, res, err := registerCustomerUseCase.RegisterCustomer(
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
			if err = c.Send(res.Text, res.Options); err != nil {
				return err
			}
			if !isRegistered {
				return fmt.Errorf("failed to register customer: %v", cnt)
			}
			res, err = startMakeAppointmentDialogUseCase.StartMakeAppointmentDialog(ctx, tgUserId)
			if err != nil {
				return err
			}
			return c.Send(res.Text, res.Options)
		})
		bot.Handle(adapters_telegram.CancelRegisterTelegramCustomerBtn, func(c telebot.Context) error {
			return c.Send("Регистрация отменена", &telebot.SendOptions{
				ReplyMarkup: &telebot.ReplyMarkup{RemoveKeyboard: true},
			})
		})
		return nil
	}
}
