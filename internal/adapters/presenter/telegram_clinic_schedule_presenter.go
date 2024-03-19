package presenter

import (
	"fmt"
	"net/url"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"gopkg.in/telebot.v3"
)

type telegramClinicSchedulePresenter struct {
	calendarWebAppUrl           adapters.CalendarWebAppUrl
	calendarInputRequestOptions string
}

func newTelegramClinicSchedulePresenter(
	calendarWebAppUrl adapters.CalendarWebAppUrl,
	calendarWebHandlerUrl adapters.CalendarWebHandlerUrl,
) telegramClinicSchedulePresenter {
	return telegramClinicSchedulePresenter{
		calendarWebAppUrl:           calendarWebAppUrl,
		calendarInputRequestOptions: fmt.Sprintf(`{"url":"%s"}`, calendarWebHandlerUrl),
	}
}

func (p *telegramClinicSchedulePresenter) scheduleButtons(schedule entity.Schedule) []telebot.InlineButton {
	buttons := make([]telebot.InlineButton, 0, 3)
	if schedule.PrevDate != nil {
		buttons = append(buttons, *adapters.PreviousClinicScheduleBtn.With(schedule.PrevDate.Format(time.DateOnly)))
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
		buttons = append(buttons, *adapters.NextClinicScheduleBtn.With(schedule.NextDate.Format(time.DateOnly)))
	}
	return buttons
}

type TelegramClinicScheduleTextPresenter struct {
	telegramClinicSchedulePresenter
}

func NewTelegramClinicScheduleTextPresenter(
	calendarWebAppUrl adapters.CalendarWebAppUrl,
	calendarWebHandlerUrl adapters.CalendarWebHandlerUrl,
) *TelegramClinicScheduleTextPresenter {
	return &TelegramClinicScheduleTextPresenter{
		telegramClinicSchedulePresenter: newTelegramClinicSchedulePresenter(
			calendarWebAppUrl,
			calendarWebHandlerUrl,
		),
	}
}

func (p *TelegramClinicScheduleTextPresenter) RenderSchedule(schedule entity.Schedule) (adapters.TelegramTextResponse, error) {
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

type TelegramClinicScheduleQueryPresenter struct {
	telegramClinicSchedulePresenter
}

func NewTelegramClinicScheduleQueryPresenter(
	calendarWebAppUrl adapters.CalendarWebAppUrl,
	calendarWebHandlerUrl adapters.CalendarWebHandlerUrl,
) *TelegramClinicScheduleQueryPresenter {
	return &TelegramClinicScheduleQueryPresenter{
		telegramClinicSchedulePresenter: newTelegramClinicSchedulePresenter(
			calendarWebAppUrl,
			calendarWebHandlerUrl,
		),
	}
}

func (p *TelegramClinicScheduleQueryPresenter) RenderSchedule(schedule entity.Schedule) (adapters.TelegramQueryResponse, error) {
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
