package appointment_telegram_controller

import (
	"context"

	appointment_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case"
	adapters_telegram "github.com/x0k/veterinary-clinic-backend/internal/infra/telegram"
	"gopkg.in/telebot.v3"
)

func UseServices(
	ctx context.Context,
	bot *telebot.Bot,
	servicesUseCase *appointment_use_case.ServicesUseCase[adapters_telegram.TextResponse],
) error {
	servicesHandler := func(c telebot.Context) error {
		res, err := servicesUseCase.Services(ctx)
		if err != nil {
			return err
		}
		return c.Send(res.Text, res.Options)
	}
	bot.Handle("/services", servicesHandler)
	bot.Handle(adapters_telegram.ServicesBtn, servicesHandler)
	return nil
}
