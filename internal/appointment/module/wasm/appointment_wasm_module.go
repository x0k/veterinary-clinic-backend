//go:build js && wasm

package appointment_wasm_module

import (
	"syscall/js"

	"github.com/jomei/notionapi"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_js_controller "github.com/x0k/veterinary-clinic-backend/internal/appointment/controller/js"
	appointment_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
)

func New(
	cfg *Config,
	log *logger.Logger,
	notion *notionapi.Client,
) js.Value {
	m := js_adapters.ObjectConstructor.New()

	schedulingService := appointment.NewSchedulingService(
		log,
		cfg.SchedulingService.SampleRateInMinutes,
	)

	appointment_js_controller.NewSchedule(
		m,
		appointment_use_case.NewScheduleUseCase(
			log,
			schedulingService,
			_,
		),
	)
	return m
}
