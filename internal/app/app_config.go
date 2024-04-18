package app

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jomei/notionapi"
	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	appointment_module "github.com/x0k/veterinary-clinic-backend/internal/appointment/module"
	profiler_module "github.com/x0k/veterinary-clinic-backend/internal/profiler"
)

type MetricsConfig struct {
	Enabled bool   `yaml:"enabled" env:"METRICS_ENABLED"`
	Address string `yaml:"address" env:"METRICS_ADDRESS"`
}

type NotionConfig struct {
	Token notionapi.Token `yaml:"token" env:"NOTION_TOKEN" env-required:"true"`
}

type TelegramConfig struct {
	Token          telegram_adapters.Token `yaml:"token" env:"TELEGRAM_TOKEN" env-required:"true"`
	PollerTimeout  time.Duration           `yaml:"poller_timeout" env:"TELEGRAM_POLLER_TIMEOUT" env-default:"10s"`
	InitDataExpiry time.Duration           `yaml:"init_data_expiry" env:"TELEGRAM_INIT_DATA_EXPIRY" env-default:"24h"`
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
