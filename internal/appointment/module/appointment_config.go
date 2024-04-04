package appointment_module

import "github.com/jomei/notionapi"

type NotionConfig struct {
	ServicesDatabaseId notionapi.DatabaseID `yaml:"services_database_id" env:"APPOINTMENT_NOTION_SERVICES_DATABASE_ID" env-required:"true"`
}

type AppointmentConfig struct {
	Notion NotionConfig `yaml:"notion"`
}
