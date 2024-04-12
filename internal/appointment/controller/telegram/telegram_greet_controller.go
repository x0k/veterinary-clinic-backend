package appointment_telegram_controller

import (
	"context"

	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	appointment_telegram_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case/telegram"
	"gopkg.in/telebot.v3"
)

func NewGreet(
	bot *telebot.Bot,
	greetUseCase *appointment_telegram_use_case.GreetUseCase[telegram_adapters.TextResponses],
) func(context.Context) error {
	return func(ctx context.Context) error {
		bot.Handle("/start", func(c telebot.Context) error {
			res, err := greetUseCase.Greet(ctx)
			if err != nil {
				return err
			}
			return res.Send(c)
		})
		return bot.SetCommands([]telebot.Command{
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
			{
				Text:        "/appointment",
				Description: "Запись на прием",
			},
		})
	}
}
