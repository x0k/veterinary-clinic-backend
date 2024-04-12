package appointment_module

import (
	"github.com/jomei/notionapi"
	production_calendar_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/production_calendar"
	web_calendar_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/web_calendar"
)

type NotionConfig struct {
	ServicesDatabaseId  notionapi.DatabaseID `yaml:"services_database_id" env:"APPOINTMENT_NOTION_SERVICES_DATABASE_ID" env-required:"true"`
	RecordsDatabaseId   notionapi.DatabaseID `yaml:"records_database_id" env:"APPOINTMENT_NOTION_RECORDS_DATABASE_ID" env-required:"true"`
	BreaksDatabaseId    notionapi.DatabaseID `yaml:"breaks_database_id" env:"APPOINTMENT_NOTION_BREAKS_DATABASE_ID" env-required:"true"`
	CustomersDatabaseId notionapi.DatabaseID `yaml:"customers_database_id" env:"APPOINTMENT_NOTION_CUSTOMERS_DATABASE_ID" env-required:"true"`
}

type ProductionCalendarConfig struct {
	Url                   production_calendar_adapters.Url `yaml:"url" env:"APPOINTMENT_PRODUCTION_CALENDAR_URL" env-required:"true"`
	TLSInsecureSkipVerify bool                             `yaml:"tls_insecure_skip_verify" env:"APPOINTMENT_PRODUCTION_CALENDAR_TLS_INSECURE_SKIP_VERIFY" env-default:"false"`
}

type WebCalendarConfig struct {
	AppUrl         web_calendar_adapters.AppUrl         `yaml:"app_url" env:"APPOINTMENT_WEB_CALENDAR_APP_URL" env-required:"true"`
	HandlerAddress web_calendar_adapters.HandlerAddress `yaml:"handler_address" env:"APPOINTMENT_WEB_CALENDAR_HANDLER_ADDRESS" env-required:"true"`
	HandlerUrlRoot web_calendar_adapters.HandlerUrlRoot `yaml:"handler_url_root" env:"APPOINTMENT_WEB_CALENDAR_HANDLER_URL_ROOT" env-required:"true"`
}

type Config struct {
	Notion             NotionConfig             `yaml:"notion"`
	ProductionCalendar ProductionCalendarConfig `yaml:"production_calendar"`
	WebCalendar        WebCalendarConfig        `yaml:"web_calendar"`
}
