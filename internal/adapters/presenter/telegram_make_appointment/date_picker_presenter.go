package telegram_make_appointment

import (
	"fmt"
	"net/url"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters/presenter"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"gopkg.in/telebot.v3"
)

type telegramDatePickerPresenter struct {
	calendarWebAppUrl           adapters.CalendarWebAppUrl
	calendarInputRequestOptions string
	stateSaver                  adapters.StateSaver[adapters.TelegramDatePickerState]
}

func newTelegramDatePickerPresenter(
	calendarWebAppUrl adapters.CalendarWebAppUrl,
	makeAppointmentDatePickerHandlerUrl adapters.MakeAppointmentDatePickerHandlerUrl,
	stateSaver adapters.StateSaver[adapters.TelegramDatePickerState],
) telegramDatePickerPresenter {
	return telegramDatePickerPresenter{
		calendarWebAppUrl:           calendarWebAppUrl,
		calendarInputRequestOptions: fmt.Sprintf(`{"url":"%s"}`, makeAppointmentDatePickerHandlerUrl),
		stateSaver:                  stateSaver,
	}
}

func (p *telegramDatePickerPresenter) buttons(serviceId entity.ServiceId, schedule entity.Schedule) [][]telebot.InlineButton {
	buttons := make([]telebot.InlineButton, 0, 3)
	if schedule.PrevDate != nil {
		buttons = append(buttons, *adapters.PrevMakeAppointmentDateBtn.With(string(
			p.stateSaver.Save(adapters.TelegramDatePickerState{
				ServiceId: serviceId,
				Date:      *schedule.PrevDate,
			}),
		)))
	}
	webAppParams := url.Values{}
	webAppParams.Add("r", p.calendarInputRequestOptions)
	webAppParams.Add("v", presenter.CalendarInputValidationSchema)
	webAppParams.Add("w", fmt.Sprintf(
		presenter.CalendarWebAppOptionsTemplate,
		time.Now().Format(time.DateOnly),
		schedule.Date.Format(time.DateOnly),
	))
	webAppParams.Add("s", string(serviceId))
	url := fmt.Sprintf("%s?%s", p.calendarWebAppUrl, webAppParams.Encode())
	buttons = append(buttons, telebot.InlineButton{
		Text: "📅",
		WebApp: &telebot.WebApp{
			URL: url,
		},
	})
	if schedule.NextDate != nil {
		buttons = append(buttons, *adapters.NextMakeAppointmentDateBtn.With(string(
			p.stateSaver.Save(adapters.TelegramDatePickerState{
				ServiceId: serviceId,
				Date:      *schedule.NextDate,
			}),
		)))
	}
	return [][]telebot.InlineButton{
		buttons,
		{
			*adapters.CancelMakeAppointmentDateBtn,
			*adapters.SelectMakeAppointmentDateBtn.With(string(
				p.stateSaver.Save(adapters.TelegramDatePickerState{
					ServiceId: serviceId,
					Date:      schedule.Date,
				}),
			)),
		},
	}
}

type TelegramDatePickerTextPresenter struct {
	telegramDatePickerPresenter
}

func NewTelegramDatePickerTextPresenter(
	calendarWebAppUrl adapters.CalendarWebAppUrl,
	makeAppointmentDatePickerHandlerUrl adapters.MakeAppointmentDatePickerHandlerUrl,
	stateSaver adapters.StateSaver[adapters.TelegramDatePickerState],
) *TelegramDatePickerTextPresenter {
	return &TelegramDatePickerTextPresenter{
		telegramDatePickerPresenter: newTelegramDatePickerPresenter(
			calendarWebAppUrl,
			makeAppointmentDatePickerHandlerUrl,
			stateSaver,
		),
	}
}

func (p *TelegramDatePickerTextPresenter) RenderDatePicker(serviceId entity.ServiceId, schedule entity.Schedule) (adapters.TelegramTextResponse, error) {
	return adapters.TelegramTextResponse{
		Text: presenter.RenderSchedule(schedule),
		Options: &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdownV2,
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: p.buttons(serviceId, schedule),
			},
		},
	}, nil
}

type TelegramDatePickerQueryPresenter struct {
	telegramDatePickerPresenter
}

func NewTelegramDatePickerQueryPresenter(
	calendarWebAppUrl adapters.CalendarWebAppUrl,
	makeAppointmentDatePickerHandlerUrl adapters.MakeAppointmentDatePickerHandlerUrl,
	stateSaver adapters.StateSaver[adapters.TelegramDatePickerState],
) *TelegramDatePickerQueryPresenter {
	return &TelegramDatePickerQueryPresenter{
		telegramDatePickerPresenter: newTelegramDatePickerPresenter(
			calendarWebAppUrl,
			makeAppointmentDatePickerHandlerUrl,
			stateSaver,
		),
	}
}

func (p *TelegramDatePickerQueryPresenter) RenderDatePicker(serviceId entity.ServiceId, schedule entity.Schedule) (adapters.TelegramQueryResponse, error) {
	return adapters.TelegramQueryResponse{
		Result: &telebot.ArticleResult{
			ResultBase: telebot.ResultBase{
				ID:        fmt.Sprintf("%p", &schedule),
				Type:      "article",
				ParseMode: telebot.ModeMarkdownV2,
				ReplyMarkup: &telebot.ReplyMarkup{
					InlineKeyboard: p.buttons(serviceId, schedule),
				},
			},
			Title: "Выберите дату:",
			Text:  presenter.RenderSchedule(schedule),
		},
	}, nil
}
