package adapters_web_calendar

import (
	"fmt"
	"net/url"
	"strings"
)

const HandlerPath = "/web-calendar"

const AppInputValidationSchema = `{"type":"object","properties":{"selectedDates":{"type":"array","minItems":1}},"required":["selectedDates"]}`

const AppOptionsTemplate = `{"date":{"min":"%s"},"settings":{"selected":{"dates":["%s"]}}}`

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

type HandlerUrlRoot string

func (h HandlerUrlRoot) String() string {
	return string(h)
}

type HandlerUrl string

func NewHandlerUrl(root HandlerUrlRoot) HandlerUrl {
	path := HandlerPath
	if strings.HasSuffix(root.String(), "/") {
		path = path[1:]
	}
	return HandlerUrl(fmt.Sprintf("%s%s", root, path))
}

type HandlerAddress string

func (h HandlerAddress) String() string {
	return string(h)
}

type AppResultResponse struct {
	Data struct {
		SelectedDates []string `json:"selectedDates"`
	} `json:"data"`
	WebAppInitData string `json:"webAppInitData"`
	State          string `json:"state"`
}
