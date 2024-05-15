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
	ServicesDatabaseId  notionapi.DatabaseID `js:"servicesDatabaseId"`
	RecordsDatabaseId   notionapi.DatabaseID `js:"recordsDatabaseId"`
	BreaksDatabaseId    notionapi.DatabaseID `js:"breaksDatabaseId"`
	CustomersDatabaseId notionapi.DatabaseID `js:"customersDatabaseId"`
}

type Config struct {
	SchedulingService  SchedulingServiceConfig  `js:"schedulingService"`
	Notion             NotionConfig             `js:"notion"`
	ProductionCalendar ProductionCalendarConfig `js:"productionCalendar"`
}
