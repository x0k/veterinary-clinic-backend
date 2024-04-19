//go:build js && wasm

package wasm_appointment_module

import (
	"syscall/js"

	"github.com/jomei/notionapi"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_js_controller "github.com/x0k/veterinary-clinic-backend/internal/appointment/controller/js"
	appointment_js_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/js"
	appointment_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
)

func New(
	cfg *Config,
	log *logger.Logger,
	notion *notionapi.Client,
) (js.Value, error) {
	m := js_adapters.ObjectConstructor.New()

	recordsRepository := appointment_js_repository.NewRecordRepository(
		cfg.RecordsRepository,
	)
	productionCalendarRepository := appointment_js_repository.NewProductionCalendarRepository(
		cfg.ProductionCalendarRepository,
	)

	schedulingService := appointment.NewSchedulingService(
		log,
		cfg.SchedulingService.SampleRateInMinutes,
		recordsRepository.CreateRecord,
		productionCalendarRepository.ProductionCalendar,
	)

	appointment_js_controller.NewSchedule(
		m,
		appointment_use_case.NewScheduleUseCase(
			log,
			schedulingService,
			_,
		),
	)
	return m, nil
}
