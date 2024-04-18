package appointment_telegram_controller

import (
	"context"

	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	appointment_telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/telegram"
	appointment_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/module"
	"gopkg.in/telebot.v3"
)

func NewServices(
	bot *telebot.Bot,
	servicesUseCase *appointment_use_case.ServicesUseCase[telegram_adapters.TextResponses],
) module.Hook {
	return module.NewHook(
		"appointment_telegram_controller.NewServices", func(ctx context.Context) error {
			servicesHandler := func(c telebot.Context) error {
				res, err := servicesUseCase.Services(ctx)
				if err != nil {
					return err
				}
				return res.Send(c)
			}
			bot.Handle("/services", servicesHandler)
			bot.Handle(appointment_telegram_adapters.ServicesBtn, servicesHandler)
			return nil
		},
	)
}
