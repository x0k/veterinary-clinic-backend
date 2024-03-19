package adapters

type TelegramToken string

type CalendarWebAppOrigin string

type CalendarWebAppUrl string

type CalendarWebHandlerUrl string

type ProductionCalendarUrl string

const CalendarWebHandlerPath = "/calendar-input"

type MakeAppointmentDatePickerHandlerUrl string

const MakeAppointmentDatePickerHandlerPath = "/make-appointment-date"

const MakeAppointmentService = "mk-app-srv"

const MakeAppointmentServiceCallback = "\f" + MakeAppointmentService

const MakeAppointmentDate = "mk-app-dt"

const MakeAppointmentDateCallback = "\f" + MakeAppointmentDate
