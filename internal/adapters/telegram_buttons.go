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
	ClinicAppointmentBtn = &telebot.InlineButton{
		Text:   "Записаться",
		Unique: "clinic-appointment",
	}
	BotMenu = &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{*ClinicScheduleBtn},
			{*ClinicServiceBtn},
			{*ClinicAppointmentBtn},
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
