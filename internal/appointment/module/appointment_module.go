package appointment_module

import (
	"context"

	"github.com/jomei/notionapi"
	appointment_telegram_controller "github.com/x0k/veterinary-clinic-backend/internal/appointment/controller/telegram"
	appointment_telegram_presenter "github.com/x0k/veterinary-clinic-backend/internal/appointment/presenter/telegram"
	appointment_notion_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/notion"
	appointment_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case"
	"github.com/x0k/veterinary-clinic-backend/internal/infra"
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
	module := module.New(log.Logger, "appointment")

	servicesRepository := appointment_notion_repository.NewService(
		notion,
		cfg.Notion.ServicesDatabaseId,
	)

	module.Append(
		servicesRepository,
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
