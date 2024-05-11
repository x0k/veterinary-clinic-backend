//go:build js && wasm

package appointment_wasm_module

import (
	"syscall/js"

	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_js_controller "github.com/x0k/veterinary-clinic-backend/internal/appointment/controller/js"
	appointment_js_presenter "github.com/x0k/veterinary-clinic-backend/internal/appointment/presenter/js"
	appointment_js_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/js"
	appointment_static_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/static"
	appointment_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case"
	appointment_js_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case/js"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
)

func New(
	cfg *Config,
	log *logger.Logger,
) js.Value {
	m := js_adapters.ObjectConstructor.New()
	// Schedule controller
	recordsRepository := appointment_js_repository.NewRecordRepository(
		cfg.RecordsRepository,
	)
	productionCalendarRepository := appointment_js_repository.NewProductionCalendarRepository(
		cfg.ProductionCalendarRepository,
	)
	workingHoursRepository := appointment_static_repository.NewWorkingHoursRepository()

	workBreaksRepository := appointment_js_repository.NewWorkBreaksRepository(
		cfg.WorkBreaksRepository,
	)
	customerRepository := appointment_js_repository.NewCustomerRepository(
		cfg.CustomerRepository,
	)
	schedulingService := appointment.NewSchedulingService(
		log,
		cfg.SchedulingService.SampleRateInMinutes,
		recordsRepository.CreateRecord,
		productionCalendarRepository.ProductionCalendar,
		workingHoursRepository.WorkingHours,
		recordsRepository.BusyPeriods,
		workBreaksRepository.WorkBreaks,
		recordsRepository.CustomerActiveAppointment,
		recordsRepository.RemoveRecord,
	)
	appointment_js_controller.NewSchedule(
		m,
		appointment_use_case.NewScheduleUseCase(
			log,
			schedulingService,
			appointment_js_presenter.SchedulePresenter,
			appointment_js_presenter.ErrorPresenter,
		),
		appointment_js_use_case.NewDayOrNextWorkingDayUseCase(
			log,
			productionCalendarRepository.ProductionCalendar,
			appointment_js_presenter.DayPresenter,
			appointment_js_presenter.ErrorPresenter,
		),
		appointment_js_use_case.NewUpsertCustomerUseCase(
			log,
			customerRepository.CustomerByIdentity,
			customerRepository.CreateCustomer,
			customerRepository.UpdateCustomer,
			appointment_js_presenter.CustomerPresenter,
			appointment_js_presenter.ErrorPresenter,
		),
	)
	return m
}
