package adapters_web_calendar

type AppUrl string

type HandlerUrl string

type HandlerAddress string

const AppInputValidationSchema = `{"type":"object","properties":{"selectedDates":{"type":"array","minItems":1}},"required":["selectedDates"]}`

const AppOptionsTemplate = `{"date":{"min":"%s"},"settings":{"selected":{"dates":["%s"]}}}`
