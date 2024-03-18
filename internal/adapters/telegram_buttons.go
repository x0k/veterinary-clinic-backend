package adapters

import "gopkg.in/telebot.v3"

var (
	NextClinicScheduleBtn = &telebot.InlineButton{
		Text:   "➡",
		Unique: "next-schedule",
	}
	PreviousClinicScheduleBtn = &telebot.InlineButton{
		Text:   "⬅",
		Unique: "next-schedule",
	}
	ClinicServiceBtn = &telebot.InlineButton{
		Text:   "Услуги",
		Unique: "clinic-services",
	}
	ClinicScheduleBtn = &telebot.InlineButton{
		Text:   "График работы",
		Unique: "clinic-schedule",
	}
	BotMenu = &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{*ClinicScheduleBtn},
			{*ClinicServiceBtn},
		},
	}
)
