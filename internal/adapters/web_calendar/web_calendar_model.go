package adapters_web_calendar

import (
	"fmt"
	"net/url"
)

type AppUrl string

func (u AppUrl) String() string {
	return string(u)
}

type AppOrigin string

func NewAppOrigin(appUrl AppUrl) (AppOrigin, error) {
	u, err := url.Parse(appUrl.String())
	if err != nil {
		return "", err
	}
	return AppOrigin(fmt.Sprintf("%s://%s", u.Scheme, u.Host)), nil
}

func (o AppOrigin) String() string {
	return string(o)
}

type HandlerUrl string

type HandlerAddress string

func (h HandlerAddress) String() string {
	return string(h)
}

const HandlerPath = "/web-calendar"

const AppInputValidationSchema = `{"type":"object","properties":{"selectedDates":{"type":"array","minItems":1}},"required":["selectedDates"]}`

const AppOptionsTemplate = `{"date":{"min":"%s"},"settings":{"selected":{"dates":["%s"]}}}`

type AppResultResponse struct {
	Data struct {
		SelectedDates []string `json:"selectedDates"`
	} `json:"data"`
	WebAppInitData string `json:"webAppInitData"`
	State          string `json:"state"`
}
