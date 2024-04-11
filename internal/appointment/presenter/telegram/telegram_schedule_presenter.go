package appointment_telegram_presenter

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	adapters_telegram "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	adapters_web_calendar "github.com/x0k/veterinary-clinic-backend/internal/adapters/web_calendar"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"gopkg.in/telebot.v3"
)

type schedulePresenter struct {
	calendarWebAppUrl           adapters_web_calendar.AppUrl
	calendarInputRequestOptions string
}

func newSchedulePresenter(
	webCalendarAppUrl adapters_web_calendar.AppUrl,
	webCalendarHandlerUrl adapters_web_calendar.HandlerUrl,
) schedulePresenter {
	return schedulePresenter{
		calendarWebAppUrl:           webCalendarAppUrl,
		calendarInputRequestOptions: fmt.Sprintf(`{"url":"%s"}`, webCalendarHandlerUrl),
	}
}

func (p *schedulePresenter) scheduleButtons(now time.Time, schedule appointment.Schedule) []telebot.InlineButton {
	buttons := make([]telebot.InlineButton, 0, 3)
	if now.Add(-24 * time.Hour).Before(schedule.PrevDate) {
		buttons = append(buttons, *adapters_telegram.PreviousScheduleBtn.With(schedule.PrevDate.Format(time.DateOnly)))
	}
	webAppParams := url.Values{}
	webAppParams.Add("r", p.calendarInputRequestOptions)
	webAppParams.Add("v", adapters_web_calendar.AppInputValidationSchema)
	webAppParams.Add("w", fmt.Sprintf(
		adapters_web_calendar.AppOptionsTemplate,
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
	buttons = append(buttons, *adapters_telegram.NextScheduleBtn.With(schedule.NextDate.Format(time.DateOnly)))
	return buttons
}

type ScheduleTextPresenter struct {
	schedulePresenter
}

func NewScheduleTextPresenter(
	calendarWebAppUrl adapters_web_calendar.AppUrl,
	calendarWebHandlerUrl adapters_web_calendar.HandlerUrl,
) *ScheduleTextPresenter {
	return &ScheduleTextPresenter{
		schedulePresenter: newSchedulePresenter(
			calendarWebAppUrl,
			calendarWebHandlerUrl,
		),
	}
}

func (p *ScheduleTextPresenter) RenderSchedule(now time.Time, schedule appointment.Schedule) (adapters_telegram.TextResponses, error) {
	sb := strings.Builder{}
	writeSchedule(&sb, schedule)
	return adapters_telegram.TextResponses{{
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
	calendarWebAppUrl adapters_web_calendar.AppUrl,
	calendarWebHandlerUrl adapters_web_calendar.HandlerUrl,
) *ScheduleQueryPresenter {
	return &ScheduleQueryPresenter{
		schedulePresenter: newSchedulePresenter(
			calendarWebAppUrl,
			calendarWebHandlerUrl,
		),
	}
}

func (p *ScheduleQueryPresenter) RenderSchedule(now time.Time, schedule appointment.Schedule) (adapters_telegram.QueryResponse, error) {
	sb := strings.Builder{}
	writeSchedule(&sb, schedule)
	return adapters_telegram.QueryResponse{
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
