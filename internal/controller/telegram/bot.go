package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/config"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"gopkg.in/telebot.v3"
)

type Bot struct {
	log *logger.Logger
	bot *telebot.Bot
}

var (
	// Universal markup builders.
	menu     = &telebot.ReplyMarkup{ResizeKeyboard: true}
	selector = &telebot.ReplyMarkup{}

	// Reply buttons.
	btnHelp     = menu.Text("ℹ Help")
	btnSettings = menu.Text("⚙ Settings")

	// Inline buttons.
	//
	// Pressing it will cause the client to
	// send the bot a callback.
	//
	// Make sure Unique stays unique as per button kind
	// since it's required for callback routing to work.
	//
	btnPrev = selector.Data("⬅", "prev")
	btnNext = selector.Data("➡", "next")
)

func NewBot(cfg *config.Config, log *logger.Logger) (*Bot, error) {
	const op = "controller.telegram.NewBot"
	settings := telebot.Settings{
		Token:  cfg.TelegramToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := telebot.NewBot(settings)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to start telegram bot: %w", op, err)
	}

	menu.Reply(
		menu.Row(btnHelp),
		menu.Row(btnSettings),
	)
	selector.Inline(
		selector.Row(btnPrev, btnNext),
	)

	b.Handle("/start", func(c telebot.Context) error {
		return c.Send("Hello!", menu)
	})

	// On reply button pressed (message)
	b.Handle(&btnHelp, func(c telebot.Context) error {
		return c.Edit("Here is some help: ...")
	})

	// On inline button pressed (callback)
	b.Handle(&btnPrev, func(c telebot.Context) error {
		return c.Respond()
	})

	return &Bot{
		log: log.With(
			slog.String("component", "bot.Bot"),
		),
		bot: b,
	}, nil
}

func (b *Bot) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go b.bot.Start()
	context.AfterFunc(ctx, func() {
		defer wg.Done()
		b.bot.Stop()
	})
}
