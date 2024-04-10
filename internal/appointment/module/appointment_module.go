package appointment_module

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	adapters_http "github.com/x0k/veterinary-clinic-backend/internal/adapters/http"
	adapters_telegram "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	adapters_web_calendar "github.com/x0k/veterinary-clinic-backend/internal/adapters/web_calendar"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_http_controller "github.com/x0k/veterinary-clinic-backend/internal/appointment/controller/http"
	appointment_telegram_controller "github.com/x0k/veterinary-clinic-backend/internal/appointment/controller/telegram"
	appointment_telegram_presenter "github.com/x0k/veterinary-clinic-backend/internal/appointment/presenter/telegram"
	appointment_http_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/http"
	appointment_notion_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/notion"
	appointment_static_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/static"
	appointment_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case"
	appointment_telegram_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/module"
	"gopkg.in/telebot.v3"
)

func New(
	cfg *Config,
	log *logger.Logger,
	bot *telebot.Bot,
	notion *notionapi.Client,
	telegramInitDataParser *adapters_telegram.InitDataParser,
) (*module.Module, error) {
	m := module.New(log.Logger, "appointment")

	servicesRepository := appointment_notion_repository.NewService(
		notion,
		cfg.Notion.ServicesDatabaseId,
	)
	m.Append(servicesRepository)

	errorPresenter := appointment_telegram_presenter.NewErrorTextPresenter()

	servicesController := adapters_telegram.NewController("services_controller", appointment_telegram_controller.NewServices(
		bot,
		appointment_use_case.NewServicesUseCase(
			log,
			servicesRepository,
			appointment_telegram_presenter.NewServices(),
			errorPresenter,
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

	webCalendarHandlerUrl := adapters_web_calendar.NewHandlerUrl(cfg.WebCalendar.HandlerUrlRoot)

	scheduleController := adapters_telegram.NewController("schedule_controller", appointment_telegram_controller.NewSchedule(
		bot,
		appointment_use_case.NewScheduleUseCase(
			log,
			schedulingService,
			appointment_telegram_presenter.NewScheduleTextPresenter(
				cfg.WebCalendar.AppUrl,
				webCalendarHandlerUrl,
			),
			errorPresenter,
		),
	))
	m.Append(scheduleController)

	webCalendarAppOrigin, err := adapters_web_calendar.NewAppOrigin(
		cfg.WebCalendar.AppUrl,
	)
	if err != nil {
		return nil, err
	}

	webCalendarService := adapters_http.NewService("web_calendar_server", &http.Server{
		Addr: cfg.WebCalendar.HandlerAddress.String(),
		Handler: adapters_http.Logging(
			log,
			appointment_http_controller.UseWebCalendarRouter(
				http.NewServeMux(), log, bot,
				webCalendarAppOrigin,
				telegramInitDataParser,
				appointment_use_case.NewScheduleUseCase(
					log,
					schedulingService,
					appointment_telegram_presenter.NewScheduleQueryPresenter(
						cfg.WebCalendar.AppUrl,
						webCalendarHandlerUrl,
					),
					appointment_telegram_presenter.NewErrorQueryPresenter(),
				),
			),
		),
	}, m)
	m.Append(webCalendarService)

	customerRepository := appointment_notion_repository.NewCustomer(
		notion,
		cfg.Notion.CustomersDatabaseId,
	)

	expirableServiceIdContainer := adapters.NewExpirableStateContainer[appointment.ServiceId](
		"expirable_service_id_container",
		uint64(time.Now().UnixNano()),
		10*time.Minute,
	)
	m.Append(expirableServiceIdContainer)

	makeAppointmentController := adapters_telegram.NewController("make_appointment_controller", appointment_telegram_controller.NewMakeAppointment(
		bot,
		appointment_telegram_use_case.NewStartMakeAppointmentDialogUseCase(
			log,
			customerRepository,
			servicesRepository,
			appointment_telegram_presenter.NewServicesPickerPresenter(
				expirableServiceIdContainer,
			),
			appointment_telegram_presenter.NewRegistrationPresenter(),
			errorPresenter,
		),
	))
	m.Append(makeAppointmentController)

	return m, nil
}
