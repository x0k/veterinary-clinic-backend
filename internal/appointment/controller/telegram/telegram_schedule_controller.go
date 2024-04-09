package appointment_telegram_controller

import (
	"context"
	"time"

	adapters_telegram "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	appointment_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case"
	"gopkg.in/telebot.v3"
)

func NewSchedule(
	bot *telebot.Bot,
	scheduleUseCase *appointment_use_case.ScheduleUseCase[adapters_telegram.TextResponse],
) func(context.Context) error {
	return func(ctx context.Context) error {
		scheduleHandler := func(c telebot.Context) error {
			now := time.Now()
			res, err := scheduleUseCase.Schedule(ctx, now, now)
			if err != nil {
				return err
			}
			return c.Send(res.Text, res.Options)
		}

		bot.Handle("/schedule", scheduleHandler)
		bot.Handle(adapters_telegram.ScheduleBtn, scheduleHandler)

		bot.Handle(adapters_telegram.NextScheduleBtn, func(c telebot.Context) error {
			date, err := time.Parse(time.DateOnly, c.Data())
			if err != nil {
				return err
			}
			res, err := scheduleUseCase.Schedule(ctx, time.Now(), date)
			if err != nil {
				return err
			}
			return c.Edit(res.Text, res.Options)
		})
		return nil
	}
}
