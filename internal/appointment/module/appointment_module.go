package appointment_module

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	http_adpters "github.com/x0k/veterinary-clinic-backend/internal/adapters/http"
	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	web_calendar_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/web_calendar"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/telegram"
	appointment_http_controller "github.com/x0k/veterinary-clinic-backend/internal/appointment/controller/http"
	appointment_telegram_controller "github.com/x0k/veterinary-clinic-backend/internal/appointment/controller/telegram"
	appointment_telegram_presenter "github.com/x0k/veterinary-clinic-backend/internal/appointment/presenter/telegram"
	appointment_http_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/http"
	appointment_notion_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/notion"
	appointment_static_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/static"
	appointment_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case"
	appointment_telegram_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/module"
	"gopkg.in/telebot.v3"
)

func New(
	cfg *Config,
	log *logger.Logger,
	bot *telebot.Bot,
	notion *notionapi.Client,
	telegramInitDataParser *telegram_adapters.InitDataParser,
) (*module.Module, error) {
	m := module.New(log.Logger, "appointment")

	servicesRepository := appointment_notion_repository.NewService(
		notion,
		cfg.Notion.ServicesDatabaseId,
	)
	m.Append(servicesRepository)

	errorPresenter := appointment_telegram_presenter.NewErrorTextPresenter()

	servicesController := telegram_adapters.NewController(
		"services_controller",
		appointment_telegram_controller.NewServices(
			bot,
			appointment_use_case.NewServicesUseCase(
				log,
				servicesRepository,
				appointment_telegram_presenter.NewServices(),
				errorPresenter,
			),
		),
	)
	m.PostStart(servicesController)

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

	webCalendarService := http_adpters.NewService(
		"web_calendar_server",
		&http.Server{
			Addr: cfg.WebCalendar.HandlerAddress.String(),
			Handler: http_adpters.Logging(
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

	expirableTelegramUserIdContainer := adapters.NewExpirableStateContainer[entity.TelegramUserId](
		"expirable_telegram_user_id_container",
		uint64(time.Now().UnixNano()),
		3*time.Minute,
	)
	m.Append(expirableTelegramUserIdContainer)

	servicesPickerPresenter := appointment_telegram_presenter.NewServicesPickerPresenter(
		expirableServiceIdContainer,
	)

	startMakeAppointmentDialogController := telegram_adapters.NewController(
		"start_make_appointment_dialog_controller",
		appointment_telegram_controller.NewStartMakeAppointmentDialog(
			bot,
			expirableTelegramUserIdContainer,
			appointment_telegram_use_case.NewStartMakeAppointmentDialogUseCase(
				log,
				customerRepository,
				servicesRepository,
				servicesPickerPresenter,
				appointment_telegram_presenter.NewRegistrationPresenter(
					expirableTelegramUserIdContainer,
				),
				errorPresenter,
			),
			appointment_telegram_use_case.NewRegisterCustomerUseCase(
				log,
				customerRepository,
				servicesRepository,
				appointment_telegram_presenter.NewSuccessRegistrationPresenter(
					servicesPickerPresenter,
				),
				errorPresenter,
			),
			errorPresenter,
		),
	)
	m.PostStart(startMakeAppointmentDialogController)

	webCalendarDatePickerUrl := web_calendar_adapters.NewDatePickerUrl(
		cfg.WebCalendar.HandlerUrlRoot,
	)

	expirableAppointmentStateContainer := adapters.NewExpirableStateContainer[appointment_telegram_adapters.AppointmentSate](
		"expirable_appointment_state_container",
		uint64(time.Now().UnixNano()),
		5*time.Minute,
	)
	m.Append(expirableAppointmentStateContainer)

	makeAppointmentController := telegram_adapters.NewController(
		"make_appointment_dialog_controller",
		appointment_telegram_controller.NewMakeAppointment(
			bot,
			appointment_telegram_use_case.NewAppointmentDatePickerUseCase(
				schedulingService,
				appointment_telegram_presenter.NewDatePickerTextPresenter(
					cfg.WebCalendar.AppUrl,
					webCalendarDatePickerUrl,
					expirableAppointmentStateContainer,
				),
				errorPresenter,
			),
			errorPresenter,
			expirableServiceIdContainer,
		),
	)
	m.PostStart(makeAppointmentController)

	return m, nil
}
