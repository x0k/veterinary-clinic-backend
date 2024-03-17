package adapters

import "gopkg.in/telebot.v3"

var (
	NextScheduleBtn = &telebot.InlineButton{
		Text:   "➡",
		Unique: "next-schedule",
	}
	PreviousScheduleBtn = &telebot.InlineButton{
		Text:   "⬅",
		Unique: "next-schedule",
	}
)
