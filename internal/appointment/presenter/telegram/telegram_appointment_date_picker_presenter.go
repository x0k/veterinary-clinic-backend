package appointment_telegram_presenter

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/telegram"
	web_calendar_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/web_calendar"
	"gopkg.in/telebot.v3"
)

type datePickerPresenter struct {
	webCalendarAppUrl              web_calendar_adapters.AppUrl
	webCalendarInputRequestOptions string
	stateSaver                     adapters.StateSaver[appointment_telegram_adapters.AppointmentSate]
}

func newDatePickerPresenter(
	webCalendarAppUrl web_calendar_adapters.AppUrl,
	webCalendarHandlerUrl web_calendar_adapters.DatePickerUrl,
	stateSaver adapters.StateSaver[appointment_telegram_adapters.AppointmentSate],
) datePickerPresenter {
	return datePickerPresenter{
		webCalendarAppUrl:              webCalendarAppUrl,
		webCalendarInputRequestOptions: fmt.Sprintf(`{"url":"%s"}`, webCalendarHandlerUrl),
		stateSaver:                     stateSaver,
	}
}

func (p *datePickerPresenter) buttons(now time.Time, serviceId appointment.ServiceId, schedule appointment.Schedule) [][]telebot.InlineButton {
	buttons := make([]telebot.InlineButton, 0, 3)
	if now.Add(-24 * time.Hour).Before(schedule.PrevDate) {
		buttons = append(buttons, *appointment_telegram_adapters.PrevMakeAppointmentDateBtn.With(string(
			p.stateSaver.Save(appointment_telegram_adapters.AppointmentSate{
				ServiceId: serviceId,
				Date:      schedule.PrevDate,
			}),
		)))
	}
	webAppParams := url.Values{}
	webAppParams.Add("r", p.webCalendarInputRequestOptions)
	webAppParams.Add("v", web_calendar_adapters.AppInputValidationSchema)
	webAppParams.Add("w", fmt.Sprintf(
		web_calendar_adapters.AppOptionsTemplate,
		time.Now().Format(time.DateOnly),
		schedule.Date.Format(time.DateOnly),
	))
	webAppParams.Add("s", string(serviceId))
	url := fmt.Sprintf("%s?%s", p.webCalendarAppUrl, webAppParams.Encode())
	buttons = append(buttons, telebot.InlineButton{
		Text: "📅",
		WebApp: &telebot.WebApp{
			URL: url,
		},
	})
	buttons = append(buttons, *appointment_telegram_adapters.NextMakeAppointmentDateBtn.With(string(
		p.stateSaver.Save(appointment_telegram_adapters.AppointmentSate{
			ServiceId: serviceId,
			Date:      schedule.NextDate,
		}),
	)))
	return [][]telebot.InlineButton{
		buttons,
		{
			*appointment_telegram_adapters.CancelMakeAppointmentDateBtn,
			*appointment_telegram_adapters.SelectMakeAppointmentDateBtn.With(string(
				p.stateSaver.Save(appointment_telegram_adapters.AppointmentSate{
					ServiceId: serviceId,
					Date:      schedule.Date,
				}),
			)),
		},
	}
}

type DatePickerTextPresenter struct {
	datePickerPresenter
}

func NewDatePickerTextPresenter(
	webCalendarAppUrl web_calendar_adapters.AppUrl,
	webCalendarDatePickerUrl web_calendar_adapters.DatePickerUrl,
	stateSaver adapters.StateSaver[appointment_telegram_adapters.AppointmentSate],
) *DatePickerTextPresenter {
	return &DatePickerTextPresenter{
		datePickerPresenter: newDatePickerPresenter(
			webCalendarAppUrl,
			webCalendarDatePickerUrl,
			stateSaver,
		),
	}
}

func (p *DatePickerTextPresenter) RenderDatePicker(
	now time.Time,
	serviceId appointment.ServiceId,
	schedule appointment.Schedule,
) (telegram_adapters.TextResponses, error) {
	sb := strings.Builder{}
	writeSchedule(&sb, schedule)
	return telegram_adapters.TextResponses{{
		Text: sb.String(),
		Options: &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdownV2,
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: p.buttons(now, serviceId, schedule),
			},
		},
	}}, nil
}

type DatePickerQueryPresenter struct {
	datePickerPresenter
}

func NewDatePickerQueryPresenter(
	webCalendarAppUrl web_calendar_adapters.AppUrl,
	webCalendarDatePickerUrl web_calendar_adapters.DatePickerUrl,
	stateSaver adapters.StateSaver[appointment_telegram_adapters.AppointmentSate],
) *DatePickerQueryPresenter {
	return &DatePickerQueryPresenter{
		datePickerPresenter: newDatePickerPresenter(
			webCalendarAppUrl,
			webCalendarDatePickerUrl,
			stateSaver,
		),
	}
}

func (p *DatePickerQueryPresenter) RenderDatePicker(
	now time.Time,
	serviceId appointment.ServiceId,
	schedule appointment.Schedule,
) (telegram_adapters.QueryResponse, error) {
	sb := strings.Builder{}
	writeSchedule(&sb, schedule)
	return telegram_adapters.QueryResponse{
		Result: &telebot.ArticleResult{
			ResultBase: telebot.ResultBase{
				ID:        fmt.Sprintf("%p", &schedule),
				Type:      "article",
				ParseMode: telebot.ModeMarkdownV2,
				ReplyMarkup: &telebot.ReplyMarkup{
					InlineKeyboard: p.buttons(now, serviceId, schedule),
				},
			},
			Title: "Выберите дату:",
			Text:  sb.String(),
		},
	}, nil
}
