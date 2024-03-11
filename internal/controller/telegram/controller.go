package telegram

import (
	"gopkg.in/telebot.v3"
)

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

func UseRouter(bot *telebot.Bot) {
	menu.Reply(
		menu.Row(btnHelp),
		menu.Row(btnSettings),
	)
	selector.Inline(
		selector.Row(btnPrev, btnNext),
	)

	bot.Handle("/start", func(c telebot.Context) error {
		return c.Send("Hello!", menu)
	})

	// On reply button pressed (message)
	bot.Handle(&btnHelp, func(c telebot.Context) error {
		return c.Edit("Here is some help: ...")
	})

	// On inline button pressed (callback)
	bot.Handle(&btnPrev, func(c telebot.Context) error {
		return c.Respond()
	})
}
