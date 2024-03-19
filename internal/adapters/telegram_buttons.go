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
	ServicesBtn = &telebot.InlineButton{
		Text:   "Услуги",
		Unique: "services",
	}
	ScheduleBtn = &telebot.InlineButton{
		Text:   "График работы",
		Unique: "schedule",
	}
	AppointmentBtn = &telebot.InlineButton{
		Text:   "Записаться",
		Unique: "appointment",
	}
	BotMenu = &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{*ScheduleBtn},
			{*ServicesBtn},
			{*AppointmentBtn},
		},
	}
	NextMakeAppointmentDateBtn = &telebot.InlineButton{
		Text:   "➡",
		Unique: "nx-mk-app-dt",
	}
	PrevMakeAppointmentDateBtn = &telebot.InlineButton{
		Text:   "⬅",
		Unique: "nx-mk-app-dt",
	}
	SelectMakeAppointmentDateBtn = &telebot.InlineButton{
		Text:   "Продолжить",
		Unique: "slc-mk-app-dt",
	}
)
