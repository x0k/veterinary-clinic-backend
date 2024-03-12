package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	TextHandler   = "text"
	JSONHandler   = "json"
	PrettyHandler = "pretty"
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
)

type LoggerConfig struct {
	Level       string `yaml:"level" env:"LOGGER_LEVEL" env-default:"info"`
	HandlerType string `yaml:"handler_type" env:"LOGGER_HANDLER_TYPE" env-default:"text"`
}

type StorageConfig struct {
	Path string `yaml:"path" env:"STORAGE_PATH" env-required:"true"`
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
	Token              string `yaml:"token" env:"NOTION_TOKEN" env-required:"true"`
	ServicesDatabaseId string `yaml:"services_database_id" env:"NOTION_SERVICES_DATABASE_ID" env-required:"true"`
	RecordsDatabaseId  string `yaml:"records_database_id" env:"NOTION_RECORDS_DATABASE_ID" env-required:"true"`
}

type TelegramConfig struct {
	Token             string        `yaml:"token" env:"TELEGRAM_TOKEN" env-required:"true"`
	PollerTimeout     time.Duration `yaml:"poller_timeout" env:"TELEGRAM_POLLER_TIMEOUT" env-default:"10s"`
	WebHandlerAddress string        `yaml:"web_handler_address" env:"TELEGRAM_WEB_HANDLER_ADDRESS" env-required:"true"`
	WebHandlerUrl     string        `yaml:"web_handler_url" env:"TELEGRAM_WEB_HANDLER_URL" env-required:"true"`
}

type Config struct {
	Logger   LoggerConfig   `yaml:"logger"`
	Storage  StorageConfig  `yaml:"storage"`
	Profiler ProfilerConfig `yaml:"profiler"`
	Metrics  MetricsConfig  `yaml:"metrics"`
	Notion   NotionConfig   `yaml:"notion"`
	Telegram TelegramConfig `yaml:"telegram"`
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
