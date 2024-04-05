package appointment_module

import (
	"github.com/jomei/notionapi"
	adapters_telegram "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	appointment_telegram_controller "github.com/x0k/veterinary-clinic-backend/internal/appointment/controller/telegram"
	appointment_telegram_presenter "github.com/x0k/veterinary-clinic-backend/internal/appointment/presenter/telegram"
	appointment_notion_repository "github.com/x0k/veterinary-clinic-backend/internal/appointment/repository/notion"
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

	return m, nil
}
