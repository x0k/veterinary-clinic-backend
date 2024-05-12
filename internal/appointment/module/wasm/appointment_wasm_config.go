//go:build js && wasm

package appointment_wasm_module

import (
	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_production_calendar_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/production_calendar"
)

type SchedulingServiceConfig struct {
	SampleRateInMinutes appointment.SampleRateInMinutes `js:"sampleRateInMinutes"`
}

type ProductionCalendarConfig struct {
	Url appointment_production_calendar_adapters.Url `js:"url"`
}

type NotionConfig struct {
	ServicesDatabaseId  notionapi.DatabaseID `yaml:"services_database_id" env:"APPOINTMENT_NOTION_SERVICES_DATABASE_ID" env-required:"true"`
	RecordsDatabaseId   notionapi.DatabaseID `yaml:"records_database_id" env:"APPOINTMENT_NOTION_RECORDS_DATABASE_ID" env-required:"true"`
	BreaksDatabaseId    notionapi.DatabaseID `yaml:"breaks_database_id" env:"APPOINTMENT_NOTION_BREAKS_DATABASE_ID" env-required:"true"`
	CustomersDatabaseId notionapi.DatabaseID `yaml:"customers_database_id" env:"APPOINTMENT_NOTION_CUSTOMERS_DATABASE_ID" env-required:"true"`
}

type Config struct {
	SchedulingService  SchedulingServiceConfig  `js:"schedulingService"`
	Notion             NotionConfig             `js:"notion"`
	ProductionCalendar ProductionCalendarConfig `js:"productionCalendar"`
}
