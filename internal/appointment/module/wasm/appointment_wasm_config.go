//go:build js && wasm

package wasm_appointment_module

import (
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_js_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/js"
)

type SchedulingServiceConfig struct {
	SampleRateInMinutes appointment.SampleRateInMinutes `js:"sampleRateInMinutes"`
}

type Config struct {
	SchedulingService SchedulingServiceConfig                          `js:"schedulingService"`
	RecordsRepository appointment_js_repository.RecordRepositoryConfig `js:"recordsRepository"`
}
