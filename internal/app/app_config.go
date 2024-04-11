package app

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	adapters_telegram "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	appointment_module "github.com/x0k/veterinary-clinic-backend/internal/appointment/module"
	profiler_module "github.com/x0k/veterinary-clinic-backend/internal/profiler"
)

type StorageConfig struct {
	Path                 string `yaml:"path" env:"STORAGE_PATH" env-required:"true"`
	MigrationsPath       string `env:"STORAGE_MIGRATIONS_PATH"`
	RecordsStateFilePath string `yaml:"records_state_file_path" env:"STORAGE_RECORDS_STATE_FILE_PATH" env-required:"true"`
}

type MetricsConfig struct {
	Enabled bool   `yaml:"enabled" env:"METRICS_ENABLED"`
	Address string `yaml:"address" env:"METRICS_ADDRESS"`
}

type NotionConfig struct {
	Token notionapi.Token `yaml:"token" env:"NOTION_TOKEN" env-required:"true"`
}

type TelegramConfig struct {
	Token          adapters_telegram.Token `yaml:"token" env:"TELEGRAM_TOKEN" env-required:"true"`
	PollerTimeout  time.Duration           `yaml:"poller_timeout" env:"TELEGRAM_POLLER_TIMEOUT" env-default:"10s"`
	InitDataExpiry time.Duration           `yaml:"init_data_expiry" env:"TELEGRAM_INIT_DATA_EXPIRY" env-default:"24h"`
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

type Config struct {
	Logger   LoggerConfig   `yaml:"logger"`
	Notion   NotionConfig   `yaml:"notion"`
	Telegram TelegramConfig `yaml:"telegram"`

	Profiler    profiler_module.Config    `yaml:"profiler"`
	Appointment appointment_module.Config `yaml:"appointment"`
}

func MustLoadConfig(configPath string) *Config {
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
