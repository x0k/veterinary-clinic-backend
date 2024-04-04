package appointment_module

import (
	"context"

	"github.com/jomei/notionapi"
	appointment_telegram_controller "github.com/x0k/veterinary-clinic-backend/internal/appointment/controller/telegram"
	appointment_telegram_presenter "github.com/x0k/veterinary-clinic-backend/internal/appointment/presenter/telegram"
	appointment_notion_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/notion"
	appointment_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case"
	"github.com/x0k/veterinary-clinic-backend/internal/infra"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/module"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"gopkg.in/telebot.v3"
)

type NotionConfig struct {
	ServicesDatabaseId notionapi.DatabaseID `yaml:"services_database_id" env:"APPOINTMENT_NOTION_SERVICES_DATABASE_ID" env-required:"true"`
}

type AppointmentConfig struct {
	Notion NotionConfig `yaml:"notion"`
}

func New(
	cfg *AppointmentConfig,
	log *logger.Logger,
	bot *telebot.Bot,
	notion *notionapi.Client,
) (*module.Module, error) {
	module := module.New(log, "appointment")

	servicesRepository := appointment_notion_repository.NewService(
		notion,
		cfg.Notion.ServicesDatabaseId,
	)

	module.Append(
		infra.Starter(func(ctx context.Context) error {
			if err := appointment_telegram_controller.UseServices(
				ctx, bot,
				appointment_use_case.NewServicesUseCase(
					servicesRepository,
					appointment_telegram_presenter.NewServices(),
				),
			); err != nil {
				return err
			}
			return nil
		}),
	)

	return module, nil
}
