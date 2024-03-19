package presenter

const CalendarInputValidationSchema = `{"type":"object","properties":{"selectedDates":{"type":"array","minItems":1}},"required":["selectedDates"]}`

const CalendarWebAppOptionsTemplate = `{"date":{"min":"%s"},"settings":{"selected":{"dates":["%s"]}}}`
