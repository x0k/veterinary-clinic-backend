package appointment_telegram_controller

import (
	"context"
	"time"

	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	appointment_telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/telegram"
	appointment_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/module"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
	"gopkg.in/telebot.v3"
)

func NewSchedule(
	bot *telebot.Bot,
	scheduleUseCase *appointment_use_case.ScheduleUseCase[telegram_adapters.TextResponses],
) module.Hook {
	return module.NewHook(
		"appointment_telegram_controller.NewSchedule",
		func(ctx context.Context) error {
			scheduleHandler := func(c telebot.Context) error {
				now := shared.NewUTCTime(time.Now())
				res, err := scheduleUseCase.Schedule(ctx, now, now)
				if err != nil {
					return err
				}
				return res.Send(c)
			}

			bot.Handle("/schedule", scheduleHandler)
			bot.Handle(appointment_telegram_adapters.ScheduleBtn, scheduleHandler)

			bot.Handle(appointment_telegram_adapters.NextScheduleBtn, func(c telebot.Context) error {
				date, err := time.Parse(time.DateOnly, c.Data())
				if err != nil {
					return err
				}
				now := shared.NewUTCTime(time.Now())
				utcDate := shared.NewUTCTime(date)
				res, err := scheduleUseCase.Schedule(ctx, now, utcDate)
				if err != nil {
					return err
				}
				return res.Edit(c)
			})
			return nil
		},
	)
}
