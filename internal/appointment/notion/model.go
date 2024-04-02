package appointment_notion

import (
	"fmt"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/notion"
)

const (
	ServiceTitle             = "Наименование"
	ServiceDurationInMinutes = "Продолжительность в минутах"
	ServiceDescription       = "Описание"
	ServiceCost              = "Стоимость"
)

const (
	RecordTitle          = "ФИО"
	RecordService        = "Услуга"
	RecordPhoneNumber    = "Телефон"
	RecordEmail          = "Почта"
	RecordDateTimePeriod = "Время записи"
	RecordState          = "Статус"
	RecordUserId         = "identity"
)

const (
	BreakTitle  = "Наименование"
	BreakPeriod = "Период"
)

const (
	RecordAwaits            = "Ожидает"
	RecordDone              = "Выполнено"
	RecordNotAppear         = "Не пришел"
	RecordDoneArchived      = "Архив выполнено"
	RecordNotAppearArchived = "Архив не пришел"
)

func RecordStatus(record appointment.RecordEntity) (string, error) {
	if record.IsArchived {
		switch record.Status {
		case appointment.RecordDone:
			return RecordDoneArchived, nil
		case appointment.RecordNotAppear:
			return RecordNotAppearArchived, nil
		default:
			return "", fmt.Errorf("%w: %s", appointment.ErrUnknownRecordStatus, record.Status)
		}
	}
	switch record.Status {
	case appointment.RecordAwaits:
		return RecordAwaits, nil
	case appointment.RecordDone:
		return RecordDone, nil
	case appointment.RecordNotAppear:
		return RecordNotAppear, nil
	default:
		return "", fmt.Errorf("%w: %s", appointment.ErrUnknownRecordStatus, record.Status)
	}
}

func Service(page notionapi.Page) appointment.ServiceEntity {
	return appointment.NewService(
		appointment.NewServiceId(string(page.ID)),
		notion.Title(page.Properties, ServiceTitle),
		entity.DurationInMinutes(
			notion.Number(page.Properties, ServiceDurationInMinutes),
		),
		notion.Text(page.Properties, ServiceDescription),
		notion.Text(page.Properties, ServiceCost),
	)
}
