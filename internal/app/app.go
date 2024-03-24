package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters/controller"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters/presenter"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters/presenter/telegram_make_appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters/repo"
	"github.com/x0k/veterinary-clinic-backend/internal/config"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/infra"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/app_logger"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/boot"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase/make_appointment"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

func Run(cfg *config.Config) {
	ctx := context.Background()
	log := app_logger.MustNew(&cfg.Logger)
	if err := run(ctx, cfg, log); err != nil {
		log.Error(ctx, "failed to run", sl.Err(err))
	}
}

func run(ctx context.Context, cfg *config.Config, log *logger.Logger) error {

	log.Info(ctx, "starting application", slog.String("log_level", cfg.Logger.Level))

	b := boot.New(log)

	calendarWebAppUrl, err := url.Parse(string(cfg.Telegram.CalendarWebAppUrl))
	if err != nil {
		return err
	}
	calendarWebAppOrigin := adapters.CalendarWebAppOrigin(fmt.Sprintf("%s://%s", calendarWebAppUrl.Scheme, calendarWebAppUrl.Host))

	calendarWebHandlerUrl := adapters.CalendarWebHandlerUrl(fmt.Sprintf("%s%s", cfg.Telegram.WebHandlerUrl, adapters.CalendarWebHandlerPath))

	makeAppointmentDatePickerHandlerUrl := adapters.MakeAppointmentDatePickerHandlerUrl(fmt.Sprintf("%s%s", cfg.Telegram.WebHandlerUrl, adapters.MakeAppointmentDatePickerHandlerPath))

	notionClient := notionapi.NewClient(cfg.Notion.Token)

	productionCalendarRepo := repo.NewHttpProductionCalendar(log, cfg.ProductionCalendar.Url, &http.Client{})
	openingHoursRepo := repo.NewStaticOpeningHoursRepo()
	workBreaksRepo := repo.NewNotionWorkBreaks(
		notionClient,
		log,
		cfg.Notion.BreaksDatabaseId,
	)
	recordsRepo := repo.NewNotionRecords(
		notionClient,
		log,
		cfg.Notion.RecordsDatabaseId,
		cfg.Notion.ServicesDatabaseId,
	)
	actualRecordsStateRepo := repo.NewFsActualRecordsStateRepo(
		cfg.Storage.RecordsStateFilePath,
	)

	query := make(chan entity.DialogMessage[adapters.TelegramQueryResponse])

	seed := uint64(time.Now().UnixNano())
	serviceIdContainer := infra.NewMemoryExpirableStateContainer[entity.ServiceId](
		seed,
		10*time.Minute,
	)
	datePickerStateContainer := infra.NewMemoryExpirableStateContainer[adapters.TelegramDatePickerState](
		seed,
		10*time.Minute,
	)

	notification := make(chan entity.NotificationMessage[adapters.TelegramTextResponse])
	appointmentChangeDetector := usecase.NewAppointmentChangeDetectorUseCase(
		entity.TelegramUserIdToUserId(cfg.Telegram.AdminUserId),
		actualRecordsStateRepo,
		recordsRepo,
		notification,
		presenter.NewTelegramChangePresenter(),
	)

	bot, err := telebot.NewBot(telebot.Settings{
		Token: string(cfg.Telegram.Token),
		Poller: &telebot.LongPoller{
			Timeout: cfg.Telegram.PollerTimeout,
		},
	})
	if err != nil {
		return err
	}

	b.Append(
		productionCalendarRepo,
		workBreaksRepo,
		recordsRepo,
		actualRecordsStateRepo,
		serviceIdContainer,
		datePickerStateContainer,
		infra.NewHttpService(
			log,
			&http.Server{
				Addr: cfg.Telegram.WebHandlerAddress,
				Handler: infra.Logging(
					log,
					controller.UseHttpTelegramRouter(
						http.NewServeMux(),
						log,
						query,
						calendarWebAppOrigin,
						infra.NewTelegramInitData(
							cfg.Telegram.Token,
							24*time.Hour,
						),
						usecase.NewScheduleUseCase(
							productionCalendarRepo,
							openingHoursRepo,
							recordsRepo,
							workBreaksRepo,
							presenter.NewTelegramScheduleQueryPresenter(
								cfg.Telegram.CalendarWebAppUrl,
								calendarWebHandlerUrl,
							),
						),
						make_appointment.NewDatePickerUseCase(
							productionCalendarRepo,
							openingHoursRepo,
							recordsRepo,
							workBreaksRepo,
							telegram_make_appointment.NewTelegramDatePickerQueryPresenter(
								cfg.Telegram.CalendarWebAppUrl,
								makeAppointmentDatePickerHandlerUrl,
								datePickerStateContainer,
							),
						),
					),
				),
			},
		),
		infra.Starter(func(ctx context.Context) error {
			bot.Use(
				middleware.Logger(slog.NewLogLogger(log.Logger.Handler(), slog.LevelInfo)),
				middleware.AutoRespond(),
			)
			if err := controller.UseTelegramBotRouter(
				ctx,
				bot,
				usecase.NewGreetUseCase(
					presenter.NewTelegramGreet(),
				),
				usecase.NewServicesUseCase(
					recordsRepo,
					presenter.NewTelegramServices(),
				),
				usecase.NewScheduleUseCase(
					productionCalendarRepo,
					openingHoursRepo,
					recordsRepo,
					workBreaksRepo,
					presenter.NewTelegramScheduleTextPresenter(
						cfg.Telegram.CalendarWebAppUrl,
						calendarWebHandlerUrl,
					),
				),
				make_appointment.NewServicePickerUseCase(
					recordsRepo,
					recordsRepo,
					telegram_make_appointment.NewTelegramServicePickerPresenter(
						serviceIdContainer,
					),
					telegram_make_appointment.NewTelegramAppointmentInfoPresenter(),
				),
				serviceIdContainer,
				make_appointment.NewDatePickerUseCase(
					productionCalendarRepo,
					openingHoursRepo,
					recordsRepo,
					workBreaksRepo,
					telegram_make_appointment.NewTelegramDatePickerTextPresenter(
						cfg.Telegram.CalendarWebAppUrl,
						makeAppointmentDatePickerHandlerUrl,
						datePickerStateContainer,
					),
				),
				datePickerStateContainer,
				make_appointment.NewTimeSlotPickerUseCase(
					entity.SampleRateInMinutes(30),
					productionCalendarRepo,
					openingHoursRepo,
					recordsRepo,
					workBreaksRepo,
					recordsRepo,
					telegram_make_appointment.NewTelegramTimeSlotsPickerPresenter(
						datePickerStateContainer,
					),
				),
				make_appointment.NewAppointmentConfirmationUseCase(
					recordsRepo,
					telegram_make_appointment.NewTelegramConfirmationPresenter(
						datePickerStateContainer,
					),
				),
				make_appointment.NewMakeAppointmentUseCase(
					recordsRepo,
					recordsRepo,
					telegram_make_appointment.NewTelegramAppointmentInfoPresenter(),
				),
				usecase.NewCancelAppointmentUseCase(
					recordsRepo,
					presenter.NewTelegramCancelAppointmentPresenter(),
				),
			); err != nil {
				return err
			}
			context.AfterFunc(ctx, func() {
				bot.Stop()
			})
			bot.Start()
			return nil
		}),
		infra.Starter(func(ctx context.Context) error {
			controller.StartTelegramBotQueryHandler(ctx, log, bot, query)
			return nil
		}),
		infra.Starter(func(ctx context.Context) error {
			controller.StartTelegramBotNotificationHandler(ctx, log, bot, notification)
			return nil
		}),
		controller.NewCron(
			log,
			cfg.AppointmentChangeDetector.CheckInterval,
			appointmentChangeDetector.DetectChanges,
		),
	)

	if cfg.Profiler.Enabled {
		b.Append(infra.NewHttpService(
			log,
			&http.Server{
				Addr:    cfg.Profiler.Address,
				Handler: controller.UseHttpProfilerRouter(http.NewServeMux()),
			},
		))
	}

	b.Start(ctx)

	log.Info(ctx, "application stopped")
	return nil
}
