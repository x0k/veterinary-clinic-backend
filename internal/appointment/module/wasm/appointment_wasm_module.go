//go:build js && wasm

package appointment_wasm_module

import (
	"context"
	"net/http"
	"syscall/js"

	"github.com/jomei/notionapi"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	pubsub_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/pubsub"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_js_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/js"
	appointment_js_controller "github.com/x0k/veterinary-clinic-backend/internal/appointment/controller/js"
	appointment_js_presenter "github.com/x0k/veterinary-clinic-backend/internal/appointment/presenter/js"
	appointment_http_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/http"
	appointment_js_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/js"
	appointment_notion_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/notion"
	appointment_static_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/static"
	appointment_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case"
	appointment_js_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case/js"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/loader"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/slicex"
)

func New(
	ctx context.Context,
	cfg *Config,
	log *logger.Logger,
	httpClient *http.Client,
	notion *notionapi.Client,
) js.Value {
	m := js_adapters.ObjectConstructor.New()

	publisher := pubsub_adapters.NewNullPublisher[appointment.EventType]()

	// Schedule controller
	appointmentRepository := appointment_notion_repository.NewAppointment(
		log,
		notion,
		cfg.Notion.RecordsDatabaseId,
		cfg.Notion.ServicesDatabaseId,
		cfg.Notion.CustomersDatabaseId,
	)
	cachedServices := appointment.ServicesLoader(
		loader.WithCache(
			log, appointmentRepository.Services,
			js_adapters.NewSimpleCache(
				log, "appointment_wasm_module.services_cache",
				cfg.ServicesRepository.Cache,
				js_adapters.To(slicex.MapE(appointment_js_adapters.ServiceToDTO)),
				js_adapters.From(slicex.MapE(appointment_js_adapters.ServiceFromDTO)),
			),
		),
	)

	productionCalendarRepository := appointment_http_repository.NewProductionCalendar(
		cfg.ProductionCalendarRepository.Url,
		httpClient,
	)
	cachedProductionCalendar := appointment.ProductionCalendarLoader(
		loader.WithCache(
			log, productionCalendarRepository.ProductionCalendar,
			js_adapters.NewSimpleCache(
				log, "appointment_wasm_module.production_calendar_cache",
				cfg.ProductionCalendarRepository.Cache,
				js_adapters.To(appointment_js_adapters.ProductionCalendarToDTO),
				js_adapters.From(appointment_js_adapters.ProductionCalendarFromDTO),
			),
		),
	)

	workingHoursRepository := appointment_static_repository.NewWorkingHoursRepository()

	workBreaksRepository := appointment_notion_repository.NewWorkBreaks(
		log,
		notion,
		cfg.Notion.BreaksDatabaseId,
	)
	cachedWorkBreaks := appointment.WorkBreaksLoader(
		loader.WithCache(
			log, workBreaksRepository.WorkBreaks,
			js_adapters.NewSimpleCache(
				log, "appointment_wasm_module.work_breaks_cache",
				cfg.WorkBreaksRepository.Cache,
				js_adapters.To(
					slicex.MapEx[appointment.WorkBreaks, []appointment_js_adapters.WorkBreakDTO](
						appointment_js_adapters.WorkBreakToDTO,
					),
				),
				js_adapters.From(
					slicex.MapEx[[]appointment_js_adapters.WorkBreakDTO, appointment.WorkBreaks](
						appointment_js_adapters.WorkBreakFromDTO,
					),
				),
			),
		),
	)

	dateTimerPeriodLockRepository := appointment_js_repository.NewDateTimePeriodLocksRepository(
		cfg.DateTimeLocksRepository,
	)

	schedulingService := appointment.NewSchedulingService(
		log,
		cfg.SchedulingService.SampleRateInMinutes,
		dateTimerPeriodLockRepository.Lock,
		dateTimerPeriodLockRepository.UnLock,
		appointmentRepository.CreateAppointment,
		cachedProductionCalendar,
		workingHoursRepository.WorkingHours,
		appointmentRepository.BusyPeriods,
		cachedWorkBreaks,
		appointmentRepository.CustomerActiveAppointment,
		appointmentRepository.RemoveAppointment,
	)

	customerRepository := appointment_notion_repository.NewCustomer(
		notion,
		cfg.Notion.CustomersDatabaseId,
	)

	appointment_js_controller.NewAppointment(
		ctx, m,
		appointment_use_case.NewScheduleUseCase(
			log,
			schedulingService,
			appointment_js_presenter.SchedulePresenter,
			appointment_js_presenter.ErrorPresenter,
		),
		appointment_js_use_case.NewDayOrNextWorkingDayUseCase(
			log,
			cachedProductionCalendar,
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
		appointment_js_use_case.NewFreeTimeSlotsUseCase(
			log,
			schedulingService,
			appointmentRepository.Service,
			appointment_js_presenter.FreeTimeSlotsPresenter,
			appointment_js_presenter.ErrorPresenter,
		),
		appointment_js_use_case.NewActiveAppointmentUseCase(
			log,
			customerRepository.CustomerByIdentity,
			appointmentRepository.CustomerActiveAppointment,
			appointmentRepository.Service,
			appointment_js_presenter.AppointmentInfoPresenter,
			appointment_js_presenter.NotFoundPresenter,
			appointment_js_presenter.ErrorPresenter,
		),
		appointment_use_case.NewMakeAppointmentUseCase(
			log,
			schedulingService,
			customerRepository.CustomerByIdentity,
			appointmentRepository.Service,
			appointment_js_presenter.AppointmentInfoPresenter,
			appointment_js_presenter.ErrorPresenter,
			publisher,
		),
		appointment_use_case.NewCancelAppointmentUseCase(
			log,
			schedulingService,
			customerRepository.CustomerByIdentity,
			appointmentRepository.Service,
			appointment_js_presenter.OkPresenter,
			appointment_js_presenter.ErrorPresenter,
			publisher,
		),
		appointment_use_case.NewServicesUseCase(
			log,
			cachedServices,
			appointment_js_presenter.ServicesPresenter,
			appointment_js_presenter.ErrorPresenter,
		),
	)
	return m
}
