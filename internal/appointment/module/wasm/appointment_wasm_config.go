//go:build js && wasm

package appointment_wasm_module

import (
	"github.com/jomei/notionapi"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_production_calendar_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/production_calendar"
	appointment_js_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/js"
)

type SchedulingServiceConfig struct {
	SampleRateInMinutes appointment.SampleRateInMinutes `js:"sampleRateInMinutes"`
}

type ProductionCalendarRepositoryConfig struct {
	Url   appointment_production_calendar_adapters.Url `js:"url"`
	Cache *js_adapters.SimpleCacheConfig               `js:"cache"`
}

type ServicesRepositoryConfig struct {
	ServicesCache *js_adapters.SimpleCacheConfig `js:"servicesCache"`
	ServiceCache  *js_adapters.KeyedCacheConfig  `js:"serviceCache"`
}

type WorkBreaksRepositoryConfig struct {
	Cache *js_adapters.SimpleCacheConfig `js:"cache"`
}

type NotionConfig struct {
	ServicesDatabaseId  notionapi.DatabaseID `js:"servicesDatabaseId"`
	RecordsDatabaseId   notionapi.DatabaseID `js:"recordsDatabaseId"`
	BreaksDatabaseId    notionapi.DatabaseID `js:"breaksDatabaseId"`
	CustomersDatabaseId notionapi.DatabaseID `js:"customersDatabaseId"`
}

type Config struct {
	Notion                       NotionConfig                                                  `js:"notion"`
	ServicesRepository           ServicesRepositoryConfig                                      `js:"servicesRepository"`
	WorkBreaksRepository         WorkBreaksRepositoryConfig                                    `js:"workBreaksRepository"`
	ProductionCalendarRepository ProductionCalendarRepositoryConfig                            `js:"productionCalendar"`
	SchedulingService            SchedulingServiceConfig                                       `js:"schedulingService"`
	DateTimeLocksRepository      appointment_js_repository.DateTimePeriodLocksRepositoryConfig `js:"dateTimeLocksRepository"`
}
