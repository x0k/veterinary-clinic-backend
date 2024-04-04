package appointment_notion

import (
	"errors"
	"fmt"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/notion"
)

var ErrUnknownRecordStatus = errors.New("unknown record status")

const (
	RecordAwaits            = "Ожидает"
	RecordDone              = "Выполнено"
	RecordNotAppear         = "Не пришел"
	RecordDoneArchived      = "Архив выполнено"
	RecordNotAppearArchived = "Архив не пришел"
)

func RecordStatusToNotion(status appointment.RecordStatus, isArchived bool) (string, error) {
	if isArchived {
		switch status {
		case appointment.RecordDone:
			return RecordDoneArchived, nil
		case appointment.RecordNotAppear:
			return RecordNotAppearArchived, nil
		default:
			return "", fmt.Errorf("%w: %s", ErrUnknownRecordStatus, status)
		}
	}
	switch status {
	case appointment.RecordAwaits:
		return RecordAwaits, nil
	case appointment.RecordDone:
		return RecordDone, nil
	case appointment.RecordNotAppear:
		return RecordNotAppear, nil
	default:
		return "", fmt.Errorf("%w: %s", ErrUnknownRecordStatus, status)
	}
}

const (
	ServiceTitle             = "Наименование"
	ServiceDurationInMinutes = "Продолжительность в минутах"
	ServiceDescription       = "Описание"
	ServiceCost              = "Стоимость"
)

func NotionToService(page notionapi.Page) appointment.ServiceEntity {
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

const (
	CustomerTitle       = "ФИО"
	CustomerEmail       = "Почта"
	CustomerPhoneNumber = "Телефон"
	CustomerUserId      = "identity"
)

func NotionToCustomer(page notionapi.Page) appointment.CustomerEntity {
	return appointment.NewCustomer(
		appointment.NewCustomerId(string(page.ID)),
		notion.Title(page.Properties, CustomerTitle),
		notion.Phone(page.Properties, CustomerPhoneNumber),
		notion.Email(page.Properties, CustomerEmail),
	)
}

const (
	RecordTitle          = "Сводка"
	RecordDateTimePeriod = "Время записи"
	RecordState          = "Статус"
	RecordCustomer       = "Клиент"
	RecordService        = "Услуга"
	RecordCreatedAt      = "Дата записи"
)

const (
	BreakTitle  = "Наименование"
	BreakPeriod = "Период"
)
