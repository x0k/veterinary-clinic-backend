//go:build js && wasm

package appointment_wasm_module

import (
	"context"
	"net/http"
	"syscall/js"

	"github.com/jomei/notionapi"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_js_controller "github.com/x0k/veterinary-clinic-backend/internal/appointment/controller/js"
	appointment_js_presenter "github.com/x0k/veterinary-clinic-backend/internal/appointment/presenter/js"
	appointment_http_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/http"
	appointment_notion_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/notion"
	appointment_static_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/static"
	appointment_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case"
	appointment_js_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case/js"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
)

func New(
	ctx context.Context,
	cfg *Config,
	log *logger.Logger,
	notion *notionapi.Client,
) js.Value {
	m := js_adapters.ObjectConstructor.New()
	// Schedule controller
	appointmentRepository := appointment_notion_repository.NewAppointment(
		log,
		notion,
		cfg.Notion.RecordsDatabaseId,
		cfg.Notion.ServicesDatabaseId,
		cfg.Notion.CustomersDatabaseId,
	)
	go appointmentRepository.Start(ctx)

	productionCalendarRepository := appointment_http_repository.NewProductionCalendar(
		log,
		cfg.ProductionCalendar.Url,
		&http.Client{},
	)
	go productionCalendarRepository.Start(ctx)

	workingHoursRepository := appointment_static_repository.NewWorkingHoursRepository()

	workBreaksRepository := appointment_notion_repository.NewWorkBreaks(
		log,
		notion,
		cfg.Notion.BreaksDatabaseId,
	)
	go workBreaksRepository.Start(ctx)

	schedulingService := appointment.NewSchedulingService(
		log,
		cfg.SchedulingService.SampleRateInMinutes,
		appointmentRepository.CreateAppointment,
		productionCalendarRepository.ProductionCalendar,
		workingHoursRepository.WorkingHours,
		appointmentRepository.BusyPeriods,
		workBreaksRepository.WorkBreaks,
		appointmentRepository.CustomerActiveAppointment,
		appointmentRepository.RemoveAppointment,
	)

	customerRepository := appointment_notion_repository.NewCustomer(
		notion,
		cfg.Notion.CustomersDatabaseId,
	)

	appointment_js_controller.NewSchedule(
		ctx, m,
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
