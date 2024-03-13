package telegram

import (
	"context"
	"strconv"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
	"gopkg.in/telebot.v3"
)

type RouterConfig struct {
	CalendarHandlerUrl string
}

func UseRouter(
	ctx context.Context,
	bot *telebot.Bot,
	clinic *usecase.ClinicUseCase[shared.TelegramResponse],
	clinicDialog *usecase.ClinicDialogUseCase[shared.TelegramResponse],
	cfg *RouterConfig,
) error {
	bot.Handle("/start", func(c telebot.Context) error {
		greet, err := clinicDialog.GreetUser(ctx)
		if err != nil {
			return err
		}
		return c.Send(greet.Text, greet.Options)
	})

	bot.Handle("/services", func(c telebot.Context) error {
		res, err := clinic.Services(ctx)
		if err != nil {
			return err
		}
		return c.Send(res.Text, res.Options)
	})

	bot.Handle("/s", func(c telebot.Context) error {
		res, err := clinicDialog.StartScheduleDialog(ctx, entity.Dialog{
			Id:     entity.DialogId(strconv.FormatInt(c.Chat().ID, 10)),
			UserId: entity.UserId(strconv.FormatInt(c.Sender().ID, 10)),
		})
		if err != nil {
			return err
		}
		return c.Send(res.Text, res.Options)
	})

	return nil
}
