package appointment_module

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	adapters_cron "github.com/x0k/veterinary-clinic-backend/internal/adapters/cron"
	http_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/http"
	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/telegram"
	web_calendar_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/web_calendar"
	appointment_http_controller "github.com/x0k/veterinary-clinic-backend/internal/appointment/controller/http"
	appointment_pubsub_controller "github.com/x0k/veterinary-clinic-backend/internal/appointment/controller/pubsub"
	appointment_telegram_controller "github.com/x0k/veterinary-clinic-backend/internal/appointment/controller/telegram"
	appointment_telegram_presenter "github.com/x0k/veterinary-clinic-backend/internal/appointment/presenter/telegram"
	appointment_fs_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/fs"
	appointment_http_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/http"
	appointment_in_memory_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/memory"
	appointment_notion_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/notion"
	appointment_static_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/static"
	appointment_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case"
	appointment_telegram_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/module"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/pubsub"
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

	greetController := appointment_telegram_controller.NewGreet(
		bot,
		appointment_telegram_use_case.NewGreetUseCase(
			appointment_telegram_presenter.RenderGreeting,
		),
	)
	m.PostStart(greetController)

	publisher := pubsub.New[appointment.EventType]()

	appointmentRepository := appointment_notion_repository.NewAppointment(
		log,
		notion,
		cfg.Notion.RecordsDatabaseId,
		cfg.Notion.ServicesDatabaseId,
		cfg.Notion.CustomersDatabaseId,
	)
	m.Append(appointmentRepository)

	servicesController := appointment_telegram_controller.NewServices(
		bot,
		appointment_use_case.NewServicesUseCase(
			log,
			appointmentRepository.Services,
			appointment_telegram_presenter.ServicesPresenter,
			appointment_telegram_presenter.TextErrorPresenter,
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

	dateTimerPeriodLockRepository := appointment_in_memory_repository.NewDateTimePeriodLocksRepository()

	schedulingService := appointment.NewSchedulingService(
		log,
		cfg.SchedulingService.SampleRateInMinutes,
		dateTimerPeriodLockRepository.Lock,
		dateTimerPeriodLockRepository.UnLock,
		appointmentRepository.CreateAppointment,
		productionCalendarRepository.ProductionCalendar,
		workingHoursRepository.WorkingHours,
		appointmentRepository.BusyPeriods,
		workBreaksRepository.WorkBreaks,
		appointmentRepository.CustomerActiveAppointment,
		appointmentRepository.RemoveAppointment,
	)

	webCalendarHandlerUrl := web_calendar_adapters.NewHandlerUrl(cfg.WebCalendar.HandlerUrlRoot)

	scheduleTextPresenter := appointment_telegram_presenter.NewScheduleTextPresenter(
		cfg.WebCalendar.AppUrl,
		webCalendarHandlerUrl,
	)
	scheduleController := appointment_telegram_controller.NewSchedule(
		bot,
		appointment_use_case.NewScheduleUseCase(
			log,
			schedulingService,
			scheduleTextPresenter.RenderSchedule,
			appointment_telegram_presenter.TextErrorPresenter,
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
		"appointment_module.expirable_appointment_state_container",
		uint64(time.Now().UnixNano()),
		5*time.Minute,
	)
	m.Append(expirableAppointmentStateContainer)

	scheduleQueryPresenter := appointment_telegram_presenter.NewScheduleQueryPresenter(
		cfg.WebCalendar.AppUrl,
		webCalendarHandlerUrl,
	)
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
			scheduleQueryPresenter.RenderSchedule,
			appointment_telegram_presenter.QueryErrorPresenter,
		),
	); err != nil {
		return nil, err
	}
	datePickerQueryPresenter := appointment_telegram_presenter.NewDatePickerQueryPresenter(
		cfg.WebCalendar.AppUrl,
		webCalendarDatePickerUrl,
		expirableAppointmentStateContainer.Save,
	)
	if err := appointment_http_controller.UseDatePickerRouter(
		webCalendarServerMux,
		log,
		bot,
		webCalendarAppOrigin,
		telegramInitDataParser,
		appointment_telegram_use_case.NewAppointmentDatePickerUseCase(
			log,
			schedulingService,
			datePickerQueryPresenter.RenderDatePicker,
			appointment_telegram_presenter.QueryErrorPresenter,
		),
	); err != nil {
		return nil, err
	}

	webCalendarService := http_adapters.NewService(
		"appointment_module.web_calendar_server",
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
		"appointment_module.expirable_service_id_container",
		uint64(time.Now().UnixNano()),
		10*time.Minute,
	)
	m.Append(expirableServiceIdContainer)

	expirableTelegramUserIdContainer := adapters.NewExpirableStateContainer[shared.TelegramUserId](
		"appointment_module.expirable_telegram_user_id_container",
		uint64(time.Now().UnixNano()),
		3*time.Minute,
	)
	m.Append(expirableTelegramUserIdContainer)

	servicesPickerPresenter := appointment_telegram_presenter.NewServicesPickerPresenter(
		expirableServiceIdContainer.Save,
	)

	registrationPresenter := appointment_telegram_presenter.NewRegistrationPresenter(
		expirableTelegramUserIdContainer.SaveByKey,
	)
	startMakeAppointmentDialogUseCase := appointment_telegram_use_case.NewStartMakeAppointmentDialogUseCase(
		log,
		customerRepository.CustomerByIdentity,
		appointmentRepository.CustomerActiveAppointment,
		appointmentRepository.Services,
		appointmentRepository.Service,
		appointment_telegram_presenter.RenderAppointmentInfo,
		servicesPickerPresenter.RenderServicesList,
		registrationPresenter.RenderRegistration,
		appointment_telegram_presenter.TextErrorPresenter,
	)

	errorSender := appointment_telegram_adapters.NewErrorSender(
		appointment_telegram_presenter.TextErrorPresenter,
	)

	successRegistrationPresenter := appointment_telegram_presenter.NewSuccessRegistrationPresenter(
		servicesPickerPresenter,
	)
	startMakeAppointmentDialogController := appointment_telegram_controller.NewStartMakeAppointmentDialog(
		bot,
		expirableTelegramUserIdContainer.Pop,
		startMakeAppointmentDialogUseCase,
		appointment_telegram_use_case.NewRegisterCustomerUseCase(
			log,
			customerRepository.CreateCustomer,
			appointmentRepository.Services,
			successRegistrationPresenter.RenderSuccessRegistration,
			appointment_telegram_presenter.TextErrorPresenter,
		),
		errorSender,
	)
	m.PostStart(startMakeAppointmentDialogController)

	textDatePickerPresenter := appointment_telegram_presenter.NewDatePickerTextPresenter(
		cfg.WebCalendar.AppUrl,
		webCalendarDatePickerUrl,
		expirableAppointmentStateContainer.Save,
	)
	timePickerPresenter := appointment_telegram_presenter.NewTimePickerPresenter(
		expirableAppointmentStateContainer.Save,
	)
	confirmationPresenter := appointment_telegram_presenter.NewConfirmationPresenter(
		expirableAppointmentStateContainer.Save,
	)
	makeAppointmentController := appointment_telegram_controller.NewMakeAppointment(
		bot,
		startMakeAppointmentDialogUseCase,
		appointment_telegram_use_case.NewAppointmentDatePickerUseCase(
			log,
			schedulingService,
			textDatePickerPresenter.RenderDatePicker,
			appointment_telegram_presenter.TextErrorPresenter,
		),
		appointment_telegram_use_case.NewAppointmentTimePickerUseCase(
			log,
			schedulingService,
			appointmentRepository.Service,
			timePickerPresenter.RenderTimePicker,
			appointment_telegram_presenter.TextErrorPresenter,
		),
		appointment_telegram_use_case.NewAppointmentConfirmationUseCase(
			log,
			appointmentRepository.Service,
			confirmationPresenter.RenderConfirmation,
			appointment_telegram_presenter.TextErrorPresenter,
		),
		appointment_use_case.NewMakeAppointmentUseCase(
			log,
			schedulingService,
			customerRepository.CustomerByIdentity,
			appointmentRepository.Service,
			appointment_telegram_presenter.RenderAppointmentInfo,
			appointment_telegram_presenter.TextErrorPresenter,
			publisher,
		),
		appointment_use_case.NewCancelAppointmentUseCase(
			log,
			schedulingService,
			customerRepository.CustomerByIdentity,
			appointmentRepository.Service,
			appointment_telegram_presenter.RenderAppointmentCancel,
			appointment_telegram_presenter.CallbackErrorPresenter,
			publisher,
		),
		errorSender,
		expirableServiceIdContainer.Load,
		expirableAppointmentStateContainer.Load,
	)
	m.PostStart(makeAppointmentController)

	adminTgId, err := cfg.Notifications.AdminIdentity.ToTelegramUserId()
	if err != nil {
		return nil, err
	}
	admin := &telebot.User{
		ID: adminTgId.Int(),
	}
	telegramSender := telegram_adapters.NewSender(bot)
	appointmentCreatedEventPresenter := appointment_telegram_presenter.NewAppointmentCreatedEventPresenter(
		admin,
	)
	appointmentCanceledEventPresenter := appointment_telegram_presenter.NewAppointmentCanceledEventPresenter(
		admin,
	)
	appointmentsStateRepository := appointment_fs_repository.NewAppointmentsStateRepository(
		"appointment_module.appointments_state_repository",
		cfg.TrackingService.StatePath,
	)
	m.Append(appointmentsStateRepository)
	trackingService := appointment.NewTracking(
		appointmentRepository.ActualAppointments,
		appointmentsStateRepository.AppointmentsState,
		appointmentsStateRepository.SaveAppointmentsState,
	)
	appointmentEventsController := appointment_pubsub_controller.NewAppointmentEvents(
		publisher,
		appointment_use_case.NewSendAdminNotificationUseCase(
			log,
			telegramSender.Send,
			appointmentCreatedEventPresenter.Present,
			appointmentCanceledEventPresenter.Present,
		),
		appointment_use_case.NewSendCustomerNotificationUseCase(
			log,
			customerRepository.CustomerById,
			appointmentRepository.Service,
			telegramSender.Send,
			appointment_telegram_presenter.AppointmentChangedEventPresenter,
		),
		appointment_use_case.NewUpdateAppointmentsStateUseCase(
			log,
			trackingService,
		),
		m,
	)
	m.Append(appointmentEventsController)

	detectChangesUseCase := appointment_use_case.NewDetectChangesUseCase(
		log,
		trackingService,
		publisher,
	)
	detectChangesCronTask := adapters_cron.NewTask(
		"appointment_module.detect_changes_cron_task",
		cfg.TrackingService.TrackingInterval,
		detectChangesUseCase.DetectChanges,
	)
	m.Append(detectChangesCronTask)

	archiveAppointmentUseCase := appointment_use_case.NewArchiveAppointmentsUseCase(
		log,
		cfg.ArchivingService.ArchivingHour,
		cfg.ArchivingService.ArchivingMinute,
		appointmentRepository.ArchiveRecords,
	)
	archiveAppointmentsCronTask := adapters_cron.NewTask(
		"appointment_module.archive_appointments_cron_task",
		cfg.ArchivingService.ArchivingInterval,
		archiveAppointmentUseCase.ArchiveRecords,
	)
	m.Append(archiveAppointmentsCronTask)

	return m, nil
}
