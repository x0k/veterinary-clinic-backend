package appointment_module

import (
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	production_calendar_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/production_calendar"
	web_calendar_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/web_calendar"
)

type NotionConfig struct {
	ServicesDatabaseId  notionapi.DatabaseID `yaml:"services_database_id" env:"APPOINTMENT_NOTION_SERVICES_DATABASE_ID" env-required:"true"`
	RecordsDatabaseId   notionapi.DatabaseID `yaml:"records_database_id" env:"APPOINTMENT_NOTION_RECORDS_DATABASE_ID" env-required:"true"`
	BreaksDatabaseId    notionapi.DatabaseID `yaml:"breaks_database_id" env:"APPOINTMENT_NOTION_BREAKS_DATABASE_ID" env-required:"true"`
	CustomersDatabaseId notionapi.DatabaseID `yaml:"customers_database_id" env:"APPOINTMENT_NOTION_CUSTOMERS_DATABASE_ID" env-required:"true"`
}

type ProductionCalendarConfig struct {
	Url                   production_calendar_adapters.Url `yaml:"url" env:"APPOINTMENT_PRODUCTION_CALENDAR_URL" env-required:"true"`
	TLSInsecureSkipVerify bool                             `yaml:"tls_insecure_skip_verify" env:"APPOINTMENT_PRODUCTION_CALENDAR_TLS_INSECURE_SKIP_VERIFY" env-default:"false"`
}

type WebCalendarConfig struct {
	AppUrl         web_calendar_adapters.AppUrl         `yaml:"app_url" env:"APPOINTMENT_WEB_CALENDAR_APP_URL" env-required:"true"`
	HandlerAddress web_calendar_adapters.HandlerAddress `yaml:"handler_address" env:"APPOINTMENT_WEB_CALENDAR_HANDLER_ADDRESS" env-required:"true"`
	HandlerUrlRoot web_calendar_adapters.HandlerUrlRoot `yaml:"handler_url_root" env:"APPOINTMENT_WEB_CALENDAR_HANDLER_URL_ROOT" env-required:"true"`
}

type SchedulingServiceConfig struct {
	SampleRateInMinutes appointment.SampleRateInMinutes `yaml:"sample_rate_in_minutes" env:"APPOINTMENT_SCHEDULING_SERVICE_SAMPLE_RATE_IN_MINUTES" env-default:"30"`
}

type NotificationsConfig struct {
	AdminIdentity appointment.CustomerIdentity `yaml:"admin_identity" env:"APPOINTMENT_NOTIFICATIONS_ADMIN_IDENTITY" env-required:"true"`
}

type TrackingServiceConfig struct {
	StatePath        string        `yaml:"state_path" env:"APPOINTMENT_TRACKING_STATE_PATH" env-required:"true"`
	TrackingInterval time.Duration `yaml:"tracking_interval" env:"APPOINTMENT_TRACKING_TRACKING_INTERVAL" env-default:"1m"`
}

type ArchivingServiceConfig struct {
	ArchivingInterval time.Duration `yaml:"archiving_interval" env:"APPOINTMENT_ARCHIVING_SERVICE_ARCHIVING_INTERVAL" env-default:"24h"`
	ArchivingHour     int           `yaml:"archiving_hour" env:"APPOINTMENT_ARCHIVING_SERVICE_ARCHIVING_HOUR" env-default:"23"`
	ArchivingMinute   int           `yaml:"archiving_minute" env:"APPOINTMENT_ARCHIVING_SERVICE_ARCHIVING_MINUTE" env-default:"0"`
}

type TelegramBotConfig struct {
	CreateAppointment bool `yaml:"create_appointment" env:"APPOINTMENT_TELEGRAM_BOT_CREATE_APPOINTMENT"`
}

type Config struct {
	Notion             NotionConfig             `yaml:"notion"`
	ProductionCalendar ProductionCalendarConfig `yaml:"production_calendar"`
	WebCalendar        WebCalendarConfig        `yaml:"web_calendar"`
	SchedulingService  SchedulingServiceConfig  `yaml:"scheduling_service"`
	Notifications      NotificationsConfig      `yaml:"notifications"`
	TrackingService    TrackingServiceConfig    `yaml:"tracking_service"`
	ArchivingService   ArchivingServiceConfig   `yaml:"archiving_service"`
	TelegramBot        TelegramBotConfig        `yaml:"telegram_bot"`
}
