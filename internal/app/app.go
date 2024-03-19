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
	"github.com/x0k/veterinary-clinic-backend/internal/adapters/presenter"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters/presenter/telegram_clinic_make_appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters/repo"
	"github.com/x0k/veterinary-clinic-backend/internal/config"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/infra"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/app_logger"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/boot"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/profiler_http_server"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/telegram_bot"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/telegram_http_server"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase/clinic_make_appointment"
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

	calendarWebHandlerUrl := adapters.CalendarWebHandlerUrl(fmt.Sprintf("%s%s", cfg.Telegram.WebHandlerOrigin, adapters.CalendarWebHandlerPath))

	makeAppointmentDatePickerHandlerUrl := adapters.MakeAppointmentDatePickerHandlerUrl(fmt.Sprintf("%s%s", cfg.Telegram.WebHandlerOrigin, adapters.MakeAppointmentDatePickerHandlerPath))

	notionClient := notionapi.NewClient(cfg.Notion.Token)

	productionCalendarRepo := repo.NewHttpProductionCalendar(log, cfg.ProductionCalendar.Url, &http.Client{})
	openingHoursRepo := repo.NewStaticOpeningHoursRepo()
	busyPeriodsRepo := repo.NewBusyPeriods(log, notionClient, cfg.Notion.RecordsDatabaseId)
	workBreaksRepo := repo.NewStaticWorkBreaks()
	clinicServicesRepo := repo.NewNotionClinicServices(
		notionClient,
		cfg.Notion.ServicesDatabaseId,
		cfg.Notion.RecordsDatabaseId,
	)

	query := make(chan entity.DialogMessage[adapters.TelegramQueryResponse])

	seed := uint64(time.Now().UnixNano())
	clinicServiceIdContainer := infra.NewMemoryExpirableStateContainer[entity.ServiceId](
		seed,
		10*time.Minute,
	)
	clinicDatePickerStateContainer := infra.NewMemoryExpirableStateContainer[telegram_clinic_make_appointment.TelegramDatePickerState](
		seed,
		10*time.Minute,
	)

	b.Append(
		productionCalendarRepo,
		clinicServiceIdContainer,
		clinicDatePickerStateContainer,
		telegram_http_server.New(
			log,
			query,
			cfg.Telegram.WebHandlerAddress,
			calendarWebAppOrigin,
			usecase.NewClinicScheduleUseCase(
				productionCalendarRepo,
				openingHoursRepo,
				busyPeriodsRepo,
				workBreaksRepo,
				presenter.NewTelegramClinicScheduleQueryPresenter(
					cfg.Telegram.CalendarWebAppUrl,
					calendarWebHandlerUrl,
				),
			),
			infra.NewTelegramInitData(
				cfg.Telegram.Token,
				24*time.Hour,
			),
			clinic_make_appointment.NewDatePickerUseCase(
				productionCalendarRepo,
				openingHoursRepo,
				busyPeriodsRepo,
				workBreaksRepo,
				telegram_clinic_make_appointment.NewTelegramDatePickerQueryPresenter(
					cfg.Telegram.CalendarWebAppUrl,
					makeAppointmentDatePickerHandlerUrl,
					clinicDatePickerStateContainer,
				),
			),
		),
		telegram_bot.New(
			log,
			cfg.Telegram.Token,
			cfg.Telegram.PollerTimeout,
			query,
			usecase.NewClinicGreetUseCase(
				presenter.NewTelegramClinicGreet(),
			),
			usecase.NewClinicServicesUseCase(
				clinicServicesRepo,
				presenter.NewTelegramClinicServices(),
			),
			usecase.NewClinicScheduleUseCase(
				productionCalendarRepo,
				openingHoursRepo,
				busyPeriodsRepo,
				workBreaksRepo,
				presenter.NewTelegramClinicScheduleTextPresenter(
					cfg.Telegram.CalendarWebAppUrl,
					calendarWebHandlerUrl,
				),
			),
			clinic_make_appointment.NewServicePickerUseCase(
				clinicServicesRepo,
				telegram_clinic_make_appointment.NewTelegramServicePickerPresenter(
					clinicServiceIdContainer,
				),
			),
			clinicServiceIdContainer,
			clinic_make_appointment.NewDatePickerUseCase(
				productionCalendarRepo,
				openingHoursRepo,
				busyPeriodsRepo,
				workBreaksRepo,
				telegram_clinic_make_appointment.NewTelegramDatePickerTextPresenter(
					cfg.Telegram.CalendarWebAppUrl,
					makeAppointmentDatePickerHandlerUrl,
					clinicDatePickerStateContainer,
				),
			),
			clinicDatePickerStateContainer,
		),
	)

	if cfg.Profiler.Enabled {
		b.Append(profiler_http_server.New(log, &cfg.Profiler))
	}

	b.Start(ctx)

	log.Info(ctx, "application stopped")
	return nil
}
