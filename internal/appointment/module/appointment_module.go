package appointment_module

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	http_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/http"
	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/telegram"
	web_calendar_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/web_calendar"
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
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
	"gopkg.in/telebot.v3"
)

func New(
	cfg *Config,
	log *logger.Logger,
	bot *telebot.Bot,
	notion *notionapi.Client,
	telegramInitDataParser telegram_adapters.InitDataParser,
) (*module.Module, error) {
	m := module.New(log.Logger, "appointment")

	greetController := telegram_adapters.NewController(
		"greet_controller",
		appointment_telegram_controller.NewGreet(
			bot,
			appointment_telegram_use_case.NewGreetUseCase(
				appointment_telegram_presenter.NewGreet(),
			),
		),
	)
	m.PostStart(greetController)

	errorPresenter := appointment_telegram_presenter.NewErrorTextPresenter()
	errorQueryPresenter := appointment_telegram_presenter.NewErrorQueryPresenter()
	errorCallbackPresenter := appointment_telegram_presenter.NewErrorCallbackPresenter()

	appointmentRepository := appointment_notion_repository.NewAppointment(
		log,
		notion,
		cfg.Notion.RecordsDatabaseId,
		cfg.Notion.ServicesDatabaseId,
	)
	m.Append(appointmentRepository)

	servicesController := telegram_adapters.NewController(
		"services_controller",
		appointment_telegram_controller.NewServices(
			bot,
			appointment_use_case.NewServicesUseCase(
				log,
				appointmentRepository,
				appointment_telegram_presenter.NewServices(),
				errorPresenter,
			),
		),
	)
	m.PostStart(servicesController)

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
		cfg.SchedulingService.SampleRateInMinutes,
		appointmentRepository,
		productionCalendarRepository,
		workingHoursRepository,
		appointmentRepository,
		workBreaksRepository,
		appointmentRepository,
		appointmentRepository,
	)

	webCalendarHandlerUrl := web_calendar_adapters.NewHandlerUrl(cfg.WebCalendar.HandlerUrlRoot)

	scheduleController := telegram_adapters.NewController(
		"schedule_controller",
		appointment_telegram_controller.NewSchedule(
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
		),
	)
	m.PostStart(scheduleController)

	webCalendarAppOrigin, err := web_calendar_adapters.NewAppOrigin(
		cfg.WebCalendar.AppUrl,
	)
	if err != nil {
		return nil, err
	}

	webCalendarDatePickerUrl := web_calendar_adapters.NewDatePickerUrl(
		cfg.WebCalendar.HandlerUrlRoot,
	)

	expirableAppointmentStateContainer := adapters.NewExpirableStateContainer[appointment_telegram_adapters.AppointmentSate](
		"expirable_appointment_state_container",
		uint64(time.Now().UnixNano()),
		5*time.Minute,
	)
	m.Append(expirableAppointmentStateContainer)

	webCalendarServerMux := http.NewServeMux()
	if err := appointment_http_controller.UseWebCalendarRouter(
		webCalendarServerMux,
		log,
		bot,
		webCalendarAppOrigin,
		telegramInitDataParser,
		appointment_use_case.NewScheduleUseCase(
			log,
			schedulingService,
			appointment_telegram_presenter.NewScheduleQueryPresenter(
				cfg.WebCalendar.AppUrl,
				webCalendarHandlerUrl,
			),
			errorQueryPresenter,
		),
	); err != nil {
		return nil, err
	}
	if err := appointment_http_controller.UseDatePickerRouter(
		webCalendarServerMux,
		log,
		bot,
		webCalendarAppOrigin,
		telegramInitDataParser,
		appointment_telegram_use_case.NewAppointmentDatePickerUseCase(
			log,
			schedulingService,
			appointment_telegram_presenter.NewDatePickerQueryPresenter(
				cfg.WebCalendar.AppUrl,
				webCalendarDatePickerUrl,
				expirableAppointmentStateContainer,
			),
			errorQueryPresenter,
		),
	); err != nil {
		return nil, err
	}

	webCalendarService := http_adapters.NewService(
		"web_calendar_server",
		&http.Server{
			Addr: cfg.WebCalendar.HandlerAddress.String(),
			Handler: http_adapters.Logging(
				log,
				webCalendarServerMux,
			),
		},
		m,
	)
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

	expirableTelegramUserIdContainer := adapters.NewExpirableStateContainer[shared.TelegramUserId](
		"expirable_telegram_user_id_container",
		uint64(time.Now().UnixNano()),
		3*time.Minute,
	)
	m.Append(expirableTelegramUserIdContainer)

	servicesPickerPresenter := appointment_telegram_presenter.NewServicesPickerPresenter(
		expirableServiceIdContainer,
	)

	appointmentInfoPresenter := appointment_telegram_presenter.NewAppointmentInfoPresenter()

	startMakeAppointmentDialogUseCase := appointment_telegram_use_case.NewStartMakeAppointmentDialogUseCase(
		log,
		customerRepository,
		appointmentRepository,
		appointmentRepository,
		appointmentInfoPresenter,
		servicesPickerPresenter,
		appointment_telegram_presenter.NewRegistrationPresenter(
			expirableTelegramUserIdContainer,
		),
		errorPresenter,
	)

	errorSender := appointment_telegram_adapters.NewErrorSender(errorPresenter)

	startMakeAppointmentDialogController := telegram_adapters.NewController(
		"start_make_appointment_dialog_controller",
		appointment_telegram_controller.NewStartMakeAppointmentDialog(
			bot,
			expirableTelegramUserIdContainer,
			startMakeAppointmentDialogUseCase,
			appointment_telegram_use_case.NewRegisterCustomerUseCase(
				log,
				customerRepository,
				appointmentRepository,
				appointment_telegram_presenter.NewSuccessRegistrationPresenter(
					servicesPickerPresenter,
				),
				errorPresenter,
			),
			errorSender,
		),
	)
	m.PostStart(startMakeAppointmentDialogController)

	makeAppointmentController := telegram_adapters.NewController(
		"make_appointment_dialog_controller",
		appointment_telegram_controller.NewMakeAppointment(
			bot,
			startMakeAppointmentDialogUseCase,
			appointment_telegram_use_case.NewAppointmentDatePickerUseCase(
				log,
				schedulingService,
				appointment_telegram_presenter.NewDatePickerTextPresenter(
					cfg.WebCalendar.AppUrl,
					webCalendarDatePickerUrl,
					expirableAppointmentStateContainer,
				),
				errorPresenter,
			),
			appointment_telegram_use_case.NewAppointmentTimePickerUseCase(
				log,
				schedulingService,
				appointmentRepository,
				appointment_telegram_presenter.NewTimePickerPresenter(
					expirableAppointmentStateContainer,
				),
				errorPresenter,
			),
			appointment_telegram_use_case.NewAppointmentConfirmationUseCase(
				log,
				appointmentRepository,
				appointment_telegram_presenter.NewConfirmationPresenter(
					expirableAppointmentStateContainer,
				),
				errorPresenter,
			),
			appointment_use_case.NewMakeAppointmentUseCase(
				log,
				schedulingService,
				customerRepository,
				appointmentRepository,
				appointmentInfoPresenter,
				errorPresenter,
			),
			appointment_use_case.NewCancelAppointmentUseCase(
				log,
				schedulingService,
				customerRepository,
				appointment_telegram_presenter.NewAppointmentCancelPresenter(),
				errorCallbackPresenter,
			),
			errorSender,
			expirableServiceIdContainer,
			expirableAppointmentStateContainer,
		),
	)
	m.PostStart(makeAppointmentController)

	return m, nil
}
