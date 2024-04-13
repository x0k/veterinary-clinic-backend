package presenter

import (
	"fmt"
	"net/url"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
	"gopkg.in/telebot.v3"
)

type telegramSchedulePresenter struct {
	calendarWebAppUrl           adapters.CalendarWebAppUrl
	calendarInputRequestOptions string
}

func newTelegramSchedulePresenter(
	calendarWebAppUrl adapters.CalendarWebAppUrl,
	calendarWebHandlerUrl adapters.CalendarWebHandlerUrl,
) telegramSchedulePresenter {
	return telegramSchedulePresenter{
		calendarWebAppUrl:           calendarWebAppUrl,
		calendarInputRequestOptions: fmt.Sprintf(`{"url":"%s"}`, calendarWebHandlerUrl),
	}
}

func (p *telegramSchedulePresenter) scheduleButtons(schedule shared.Schedule) []telebot.InlineButton {
	buttons := make([]telebot.InlineButton, 0, 3)
	if schedule.PrevDate != nil {
		buttons = append(buttons, *adapters.PreviousScheduleBtn.With(schedule.PrevDate.Format(time.DateOnly)))
	}
	webAppParams := url.Values{}
	webAppParams.Add("r", p.calendarInputRequestOptions)
	webAppParams.Add("v", CalendarInputValidationSchema)
	webAppParams.Add("w", fmt.Sprintf(
		CalendarWebAppOptionsTemplate,
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
	if schedule.NextDate != nil {
		buttons = append(buttons, *adapters.NextScheduleBtn.With(schedule.NextDate.Format(time.DateOnly)))
	}
	return buttons
}

type TelegramScheduleTextPresenter struct {
	telegramSchedulePresenter
}

func NewTelegramScheduleTextPresenter(
	calendarWebAppUrl adapters.CalendarWebAppUrl,
	calendarWebHandlerUrl adapters.CalendarWebHandlerUrl,
) *TelegramScheduleTextPresenter {
	return &TelegramScheduleTextPresenter{
		telegramSchedulePresenter: newTelegramSchedulePresenter(
			calendarWebAppUrl,
			calendarWebHandlerUrl,
		),
	}
}

func (p *TelegramScheduleTextPresenter) RenderSchedule(schedule shared.Schedule) (adapters.TelegramTextResponse, error) {
	return adapters.TelegramTextResponse{
		Text: RenderSchedule(schedule),
		Options: &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdownV2,
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					p.scheduleButtons(schedule),
				},
			},
		},
	}, nil
}

type TelegramScheduleQueryPresenter struct {
	telegramSchedulePresenter
}

func NewTelegramScheduleQueryPresenter(
	calendarWebAppUrl adapters.CalendarWebAppUrl,
	calendarWebHandlerUrl adapters.CalendarWebHandlerUrl,
) *TelegramScheduleQueryPresenter {
	return &TelegramScheduleQueryPresenter{
		telegramSchedulePresenter: newTelegramSchedulePresenter(
			calendarWebAppUrl,
			calendarWebHandlerUrl,
		),
	}
}

func (p *TelegramScheduleQueryPresenter) RenderSchedule(schedule shared.Schedule) (adapters.TelegramQueryResponse, error) {
	return adapters.TelegramQueryResponse{
		Result: &telebot.ArticleResult{
			ResultBase: telebot.ResultBase{
				ID:        fmt.Sprintf("%p", &schedule),
				Type:      "article",
				ParseMode: telebot.ModeMarkdownV2,
				ReplyMarkup: &telebot.ReplyMarkup{
					InlineKeyboard: [][]telebot.InlineButton{
						p.scheduleButtons(schedule),
					},
				},
			},
			Title: "График работы",
			Text:  RenderSchedule(schedule),
		},
	}, nil
}
