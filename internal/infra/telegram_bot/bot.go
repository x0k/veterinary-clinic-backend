package telegram_bot

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters/controller"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type Bot struct {
	bot *telebot.Bot
	wg  sync.WaitGroup

	log            *logger.Logger
	telegramToken  adapters.TelegramToken
	clinicGreet    *usecase.ClinicGreetUseCase[adapters.TelegramTextResponse]
	clinicServices *usecase.ClinicServicesUseCase[adapters.TelegramTextResponse]
	clinicSchedule *usecase.ClinicScheduleUseCase[adapters.TelegramTextResponse]
	query          <-chan entity.DialogMessage[adapters.TelegramQueryResponse]
	pollerTimeout  time.Duration
}

func New(
	log *logger.Logger,
	telegramToken adapters.TelegramToken,
	clinicGreet *usecase.ClinicGreetUseCase[adapters.TelegramTextResponse],
	clinicServices *usecase.ClinicServicesUseCase[adapters.TelegramTextResponse],
	clinicSchedule *usecase.ClinicScheduleUseCase[adapters.TelegramTextResponse],
	pollerTimeout time.Duration,
) *Bot {
	return &Bot{
		log:            log,
		telegramToken:  telegramToken,
		clinicGreet:    clinicGreet,
		clinicServices: clinicServices,
		clinicSchedule: clinicSchedule,
		pollerTimeout:  pollerTimeout,
	}
}

func (b *Bot) Start(ctx context.Context) error {
	const op = "infra.telegram_bot.Bot.Start"

	if bot, err := telebot.NewBot(telebot.Settings{
		Token: string(b.telegramToken),
		Poller: &telebot.LongPoller{
			Timeout: b.pollerTimeout,
		},
	}); err != nil {
		return fmt.Errorf("%s failed to start: %w", op, err)
	} else {
		b.bot = bot
	}
	b.bot.Use(
		middleware.Logger(slog.NewLogLogger(b.log.Logger.Handler(), slog.LevelDebug)),
		middleware.AutoRespond(),
	)
	controller.UseTelegramBotRouter(
		ctx,
		b.bot,
		b.clinicGreet,
		b.clinicServices,
		b.clinicSchedule,
	)
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		controller.StartTelegramBotQueryHandler(ctx, b.log, b.bot, b.query)
	}()
	context.AfterFunc(ctx, func() {
		b.bot.Stop()
	})
	b.bot.Start()
	b.wg.Wait()
	return nil
}
