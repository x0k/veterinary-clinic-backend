//go:build js && wasm

package appointment_wasm_module

import "github.com/x0k/veterinary-clinic-backend/internal/appointment"

type SchedulingServiceConfig struct {
	SampleRateInMinutes appointment.SampleRateInMinutes
}

type Config struct {
	SchedulingService SchedulingServiceConfig
}
