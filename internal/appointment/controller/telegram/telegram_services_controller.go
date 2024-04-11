package appointment_telegram_controller

import (
	"context"

	adapters_telegram "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	appointment_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case"
	"gopkg.in/telebot.v3"
)

func NewServices(
	bot *telebot.Bot,
	servicesUseCase *appointment_use_case.ServicesUseCase[adapters_telegram.TextResponses],
) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		servicesHandler := func(c telebot.Context) error {
			res, err := servicesUseCase.Services(ctx)
			if err != nil {
				return err
			}
			return adapters_telegram.Send(c, res)
		}
		bot.Handle("/services", servicesHandler)
		bot.Handle(adapters_telegram.ServicesBtn, servicesHandler)
		return nil
	}
}
