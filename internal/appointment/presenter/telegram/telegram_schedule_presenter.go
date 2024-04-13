package appointment_telegram_presenter

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/telegram"
	web_calendar_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/web_calendar"
	"gopkg.in/telebot.v3"
)

type schedulePresenter struct {
	calendarWebAppUrl           web_calendar_adapters.AppUrl
	calendarInputRequestOptions string
}

func newSchedulePresenter(
	webCalendarAppUrl web_calendar_adapters.AppUrl,
	webCalendarHandlerUrl web_calendar_adapters.HandlerUrl,
) schedulePresenter {
	return schedulePresenter{
		calendarWebAppUrl:           webCalendarAppUrl,
		calendarInputRequestOptions: fmt.Sprintf(`{"url":"%s"}`, webCalendarHandlerUrl),
	}
}

func (p *schedulePresenter) scheduleButtons(now time.Time, schedule appointment.Schedule) []telebot.InlineButton {
	buttons := make([]telebot.InlineButton, 0, 3)
	if now.Add(-24 * time.Hour).Before(schedule.PrevDate) {
		buttons = append(buttons, *appointment_telegram_adapters.PreviousScheduleBtn.With(schedule.PrevDate.Format(time.DateOnly)))
	}
	webAppParams := url.Values{}
	webAppParams.Add("r", p.calendarInputRequestOptions)
	webAppParams.Add("v", web_calendar_adapters.AppInputValidationSchema)
	webAppParams.Add("w", fmt.Sprintf(
		web_calendar_adapters.AppOptionsTemplate,
		time.Now().Format(time.DateOnly),
		schedule.Date.Format(time.DateOnly),
	))
	url := fmt.Sprintf("%s?%s", p.calendarWebAppUrl, webAppParams.Encode())
	buttons = append(buttons, telebot.InlineButton{
		Text: "📅",
		WebApp: &telebot.WebApp{
			URL: url,
		},
	})
	buttons = append(buttons, *appointment_telegram_adapters.NextScheduleBtn.With(schedule.NextDate.Format(time.DateOnly)))
	return buttons
}

type ScheduleTextPresenter struct {
	schedulePresenter
}

func NewScheduleTextPresenter(
	calendarWebAppUrl web_calendar_adapters.AppUrl,
	calendarWebHandlerUrl web_calendar_adapters.HandlerUrl,
) *ScheduleTextPresenter {
	return &ScheduleTextPresenter{
		schedulePresenter: newSchedulePresenter(
			calendarWebAppUrl,
			calendarWebHandlerUrl,
		),
	}
}

func (p *ScheduleTextPresenter) RenderSchedule(now time.Time, schedule appointment.Schedule) (telegram_adapters.TextResponses, error) {
	sb := strings.Builder{}
	writeSchedule(&sb, schedule)
	return telegram_adapters.TextResponses{{
		Text: sb.String(),
		Options: &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdownV2,
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					p.scheduleButtons(now, schedule),
				},
			},
		},
	}}, nil
}

type ScheduleQueryPresenter struct {
	schedulePresenter
}

func NewScheduleQueryPresenter(
	calendarWebAppUrl web_calendar_adapters.AppUrl,
	calendarWebHandlerUrl web_calendar_adapters.HandlerUrl,
) *ScheduleQueryPresenter {
	return &ScheduleQueryPresenter{
		schedulePresenter: newSchedulePresenter(
			calendarWebAppUrl,
			calendarWebHandlerUrl,
		),
	}
}

func (p *ScheduleQueryPresenter) RenderSchedule(now time.Time, schedule appointment.Schedule) (telegram_adapters.QueryResponse, error) {
	sb := strings.Builder{}
	writeSchedule(&sb, schedule)
	return telegram_adapters.QueryResponse{
		Result: &telebot.ArticleResult{
			ResultBase: telebot.ResultBase{
				ID:        fmt.Sprintf("%p", &schedule),
				Type:      "article",
				ParseMode: telebot.ModeMarkdownV2,
				ReplyMarkup: &telebot.ReplyMarkup{
					InlineKeyboard: [][]telebot.InlineButton{
						p.scheduleButtons(now, schedule),
					},
				},
			},
			Title: "График работы",
			Text:  sb.String(),
		},
	}, nil
}
