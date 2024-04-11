package appointment_telegram_presenter

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	adapters_web_calendar "github.com/x0k/veterinary-clinic-backend/internal/adapters/web_calendar"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/telegram"
	"gopkg.in/telebot.v3"
)

type datePickerPresenter struct {
	webCalendarAppUrl              adapters_web_calendar.AppUrl
	webCalendarInputRequestOptions string
	stateSaver                     adapters.StateSaver[appointment_telegram_adapters.AppointmentSate]
}

func newDatePickerPresenter(
	webCalendarAppUrl adapters_web_calendar.AppUrl,
	webCalendarHandlerUrl adapters_web_calendar.DatePickerUrl,
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
		buttons = append(buttons, *adapters.PrevMakeAppointmentDateBtn.With(string(
			p.stateSaver.Save(appointment_telegram_adapters.AppointmentSate{
				ServiceId: serviceId,
				Date:      schedule.PrevDate,
			}),
		)))
	}
	webAppParams := url.Values{}
	webAppParams.Add("r", p.webCalendarInputRequestOptions)
	webAppParams.Add("v", adapters_web_calendar.AppInputValidationSchema)
	webAppParams.Add("w", fmt.Sprintf(
		adapters_web_calendar.AppOptionsTemplate,
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
	buttons = append(buttons, *adapters.NextMakeAppointmentDateBtn.With(string(
		p.stateSaver.Save(appointment_telegram_adapters.AppointmentSate{
			ServiceId: serviceId,
			Date:      schedule.NextDate,
		}),
	)))
	return [][]telebot.InlineButton{
		buttons,
		{
			*adapters.CancelMakeAppointmentDateBtn,
			*adapters.SelectMakeAppointmentDateBtn.With(string(
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
	webCalendarAppUrl adapters_web_calendar.AppUrl,
	webCalendarHandlerUrl adapters_web_calendar.DatePickerUrl,
	stateSaver adapters.StateSaver[appointment_telegram_adapters.AppointmentSate],
) *DatePickerTextPresenter {
	return &DatePickerTextPresenter{
		datePickerPresenter: newDatePickerPresenter(
			webCalendarAppUrl,
			webCalendarHandlerUrl,
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
	webCalendarAppUrl adapters_web_calendar.AppUrl,
	webCalendarHandlerUrl adapters_web_calendar.DatePickerUrl,
	stateSaver adapters.StateSaver[appointment_telegram_adapters.AppointmentSate],
) *DatePickerQueryPresenter {
	return &DatePickerQueryPresenter{
		datePickerPresenter: newDatePickerPresenter(
			webCalendarAppUrl,
			webCalendarHandlerUrl,
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
