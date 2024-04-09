package appointment_module

import (
	"github.com/jomei/notionapi"
	adapters_production_calendar "github.com/x0k/veterinary-clinic-backend/internal/adapters/production_calendar"
	adapters_web_calendar "github.com/x0k/veterinary-clinic-backend/internal/adapters/web_calendar"
)

type Notion struct {
	ServicesDatabaseId notionapi.DatabaseID `yaml:"services_database_id" env:"APPOINTMENT_NOTION_SERVICES_DATABASE_ID" env-required:"true"`
	RecordsDatabaseId  notionapi.DatabaseID `yaml:"records_database_id" env:"APPOINTMENT_NOTION_RECORDS_DATABASE_ID" env-required:"true"`
	BreaksDatabaseId   notionapi.DatabaseID `yaml:"breaks_database_id" env:"APPOINTMENT_NOTION_BREAKS_DATABASE_ID" env-required:"true"`
}

type ProductionCalendar struct {
	Url adapters_production_calendar.Url `yaml:"url" env:"APPOINTMENT_PRODUCTION_CALENDAR_URL" env-required:"true"`
}

type WebCalendar struct {
	AppUrl     adapters_web_calendar.AppUrl     `yaml:"app_url" env:"APPOINTMENT_WEB_CALENDAR_APP_URL" env-required:"true"`
	HandlerUrl adapters_web_calendar.HandlerUrl `yaml:"handler_url" env:"APPOINTMENT_WEB_CALENDAR_HANDLER_URL" env-required:"true"`
}

type Config struct {
	Notion             `yaml:"notion"`
	ProductionCalendar `yaml:"production_calendar"`
	WebCalendar        `yaml:"web_calendar"`
}
