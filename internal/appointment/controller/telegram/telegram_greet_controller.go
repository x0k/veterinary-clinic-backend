package appointment_telegram_controller

import (
	"context"

	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	appointment_telegram_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/module"
	"gopkg.in/telebot.v3"
)

func NewGreet(
	bot *telebot.Bot,
	greetUseCase *appointment_telegram_use_case.GreetUseCase[telegram_adapters.TextResponses],
	createAppointment bool,
) module.Hook {
	return module.NewHook(
		"appointment_telegram_controller.NewGreet",
		func(ctx context.Context) error {
			bot.Handle("/start", func(c telebot.Context) error {
				res, err := greetUseCase.Greet(ctx)
				if err != nil {
					return err
				}
				return res.Send(c)
			})
			commands := []telebot.Command{
				{
					Text:        "/start",
					Description: "Приветствие",
				},
				{
					Text:        "/services",
					Description: "Список услуг",
				},
				{
					Text:        "/schedule",
					Description: "График работы",
				},
			}
			if createAppointment {
				commands = append(commands, telebot.Command{
					Text:        "/appointment",
					Description: "Запись на прием",
				})
			}
			return bot.SetCommands(commands)
		},
	)
}
