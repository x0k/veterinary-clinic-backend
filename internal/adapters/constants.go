package adapters

type TelegramToken string

type CalendarWebAppOrigin string

type CalendarWebAppUrl string

type CalendarWebHandlerUrl string

type ProductionCalendarUrl string

const CalendarWebHandlerPath = "/calendar-input"

type MakeAppointmentDatePickerHandlerUrl string

const MakeAppointmentDatePickerHandlerPath = "/make-appointment-date"

const ClinicMakeAppointmentService = "cl-mk-app-srv"

const ClinicMakeAppointmentServiceCallback = "\f" + ClinicMakeAppointmentService

const ClinicMakeAppointmentDate = "cl-mk-app-dt"

const ClinicMakeAppointmentDateCallback = "\f" + ClinicMakeAppointmentDate
