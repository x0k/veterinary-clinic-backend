package appointment_module

import (
	"crypto/tls"
	"net/http"

	"github.com/jomei/notionapi"
	adapters_telegram "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_telegram_controller "github.com/x0k/veterinary-clinic-backend/internal/appointment/controller/telegram"
	appointment_telegram_presenter "github.com/x0k/veterinary-clinic-backend/internal/appointment/presenter/telegram"
	appointment_http_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/http"
	appointment_notion_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/notion"
	appointment_static_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/static"
	appointment_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/module"
	"gopkg.in/telebot.v3"
)

func New(
	cfg *Config,
	log *logger.Logger,
	bot *telebot.Bot,
	notion *notionapi.Client,
) (*module.Module, error) {
	m := module.New(log.Logger, "appointment")

	servicesRepository := appointment_notion_repository.NewService(
		notion,
		cfg.Notion.ServicesDatabaseId,
	)
	m.Append(servicesRepository)

	servicesController := adapters_telegram.NewController("services_controller", appointment_telegram_controller.NewServices(
		bot,
		appointment_use_case.NewServicesUseCase(
			servicesRepository,
			appointment_telegram_presenter.NewServices(),
		),
	))
	m.Append(servicesController)

	appointmentRepository := appointment_notion_repository.NewAppointment(
		log,
		notion,
		cfg.Notion.RecordsDatabaseId,
	)

	productionCalendarRepository := appointment_http_repository.NewProductionCalendar(
		log,
		cfg.ProductionCalendar.Url,
		&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: cfg.ProductionCalendar.TLSInsecureSkipVerify,
				},
			},
		},
	)
	m.Append(productionCalendarRepository)

	workingHoursRepository := appointment_static_repository.NewWorkingHoursRepository()

	workBreaksRepository := appointment_notion_repository.NewWorkBreaks(
		log,
		notion,
		cfg.Notion.BreaksDatabaseId,
	)
	m.Append(workBreaksRepository)

	schedulingService := appointment.NewSchedulingService(
		log,
		appointmentRepository,
		appointmentRepository,
		productionCalendarRepository,
		workingHoursRepository,
		appointmentRepository,
		workBreaksRepository,
	)

	scheduleController := adapters_telegram.NewController("schedule_controller", appointment_telegram_controller.NewSchedule(
		bot,
		appointment_use_case.NewScheduleUseCase(
			schedulingService,
			appointment_telegram_presenter.NewScheduleTextPresenter(
				cfg.WebCalendar.AppUrl,
				cfg.WebCalendar.HandlerUrl,
			),
			appointment_telegram_presenter.NewErrorPresenter(),
		),
	))
	m.Append(scheduleController)

	return m, nil
}
