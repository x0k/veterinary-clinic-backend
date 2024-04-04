package app

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	appointment_module "github.com/x0k/veterinary-clinic-backend/internal/appointment/module"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/app_logger"
)

type StorageConfig struct {
	Path                 string `yaml:"path" env:"STORAGE_PATH" env-required:"true"`
	MigrationsPath       string `env:"STORAGE_MIGRATIONS_PATH"`
	RecordsStateFilePath string `yaml:"records_state_file_path" env:"STORAGE_RECORDS_STATE_FILE_PATH" env-required:"true"`
}

type ProfilerConfig struct {
	Enabled bool   `yaml:"enabled" env:"PROFILER_ENABLED"`
	Address string `yaml:"address" env:"PROFILER_ADDRESS"`
}

type MetricsConfig struct {
	Enabled bool   `yaml:"enabled" env:"METRICS_ENABLED"`
	Address string `yaml:"address" env:"METRICS_ADDRESS"`
}

type NotionConfig struct {
	Token              notionapi.Token      `yaml:"token" env:"NOTION_TOKEN" env-required:"true"`
	ServicesDatabaseId notionapi.DatabaseID `yaml:"services_database_id" env:"NOTION_SERVICES_DATABASE_ID" env-required:"true"`
	RecordsDatabaseId  notionapi.DatabaseID `yaml:"records_database_id" env:"NOTION_RECORDS_DATABASE_ID" env-required:"true"`
	BreaksDatabaseId   notionapi.DatabaseID `yaml:"breaks_database_id" env:"NOTION_BREAKS_DATABASE_ID" env-required:"true"`
}

type TelegramConfig struct {
	Token             adapters.TelegramToken     `yaml:"token" env:"TELEGRAM_TOKEN" env-required:"true"`
	PollerTimeout     time.Duration              `yaml:"poller_timeout" env:"TELEGRAM_POLLER_TIMEOUT" env-default:"10s"`
	WebHandlerAddress string                     `yaml:"web_handler_address" env:"TELEGRAM_WEB_HANDLER_ADDRESS" env-required:"true"`
	WebHandlerUrl     string                     `yaml:"web_handler_url" env:"TELEGRAM_WEB_HANDLER_URL" env-required:"true"`
	CalendarWebAppUrl adapters.CalendarWebAppUrl `yaml:"calendar_web_app_url" env:"TELEGRAM_CALENDAR_WEB_APP_URL" env-required:"true"`
}

type ProductionCalendarConfig struct {
	Url adapters.ProductionCalendarUrl `yaml:"url" env:"PRODUCTION_CALENDAR_URL" env-required:"true"`
}

type AppointmentChangeDetectorConfig struct {
	CheckInterval time.Duration `yaml:"check_interval" env:"APPOINTMENT_CHANGE_DETECTOR_CHECK_INTERVAL" env-default:"1m"`
}

type AppointmentAutoArchiverConfig struct {
	ArchiveInterval time.Duration `yaml:"archive_interval" env:"APPOINTMENT_AUTO_ARCHIVER_ARCHIVE_INTERVAL" env-default:"24h"`
	ArchiveTime     time.Time     `yaml:"archive_time" env:"APPOINTMENT_AUTO_ARCHIVER_ARCHIVE_TIME" env-default:"2000-01-01T23:30:00+03:00"`
}

type NotificationsConfig struct {
	AdminId entity.UserId `yaml:"admin_id" env:"NOTIFICATIONS_ADMIN_ID" env-required:"true"`
}

type Config struct {
	Profiler ProfilerConfig `yaml:"profiler"`
	Notion   NotionConfig   `yaml:"notion"`
	Telegram TelegramConfig `yaml:"telegram"`

	Logger      app_logger.LoggerConfig              `yaml:"logger"`
	Appointment appointment_module.AppointmentConfig `yaml:"appointment"`
}

func MustLoad(configPath string) *Config {
	cfg := &Config{}
	var cfgErr error
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfgErr = cleanenv.ReadEnv(cfg)
	} else if err == nil {
		cfgErr = cleanenv.ReadConfig(configPath, cfg)
	} else {
		cfgErr = err
	}
	if cfgErr != nil {
		log.Fatalf("cannot read config: %s", cfgErr)
	}
	return cfg
}
