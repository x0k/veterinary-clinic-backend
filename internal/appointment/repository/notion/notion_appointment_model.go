package appointment_notion_repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/notion"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
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
		shared.DurationInMinutes(
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
	CustomerRecords     = "Записи"
)

func NotionToCustomer(page notionapi.Page) (appointment.CustomerEntity, error) {
	identity, err := appointment.NewCustomerIdentity(notion.Text(page.Properties, CustomerUserId))
	if err != nil {
		return appointment.CustomerEntity{}, err
	}
	return appointment.NewCustomer(
		appointment.NewCustomerId(string(page.ID)),
		identity,
		notion.Title(page.Properties, CustomerTitle),
		notion.Phone(page.Properties, CustomerPhoneNumber),
		notion.Email(page.Properties, CustomerEmail),
	), nil
}

const (
	RecordTitle          = "Сводка"
	RecordDateTimePeriod = "Время записи"
	RecordState          = "Статус"
	RecordCustomer       = "Клиент"
	RecordService        = "Услуга"
	RecordCreatedAt      = "Дата записи"
)

func NotionToRecordStatus(notionStatus string) (appointment.RecordStatus, bool, error) {
	switch notionStatus {
	case RecordAwaits:
		return appointment.RecordAwaits, false, nil
	case RecordDone:
		return appointment.RecordDone, false, nil
	case RecordNotAppear:
		return appointment.RecordNotAppear, false, nil
	case RecordDoneArchived:
		return appointment.RecordDone, true, nil
	case RecordNotAppearArchived:
		return appointment.RecordNotAppear, true, nil
	default:
		return "", false, fmt.Errorf("%w: %s", ErrUnknownRecordStatus, notionStatus)
	}
}

func NotionToRecord(page notionapi.Page) (appointment.RecordEntity, error) {
	status, isArchived, err := NotionToRecordStatus(notion.Select(page.Properties, RecordState))
	if err != nil {
		return appointment.RecordEntity{}, err
	}
	period, err := notion.DatePeriod(page.Properties, RecordDateTimePeriod)
	if err != nil {
		return appointment.RecordEntity{}, err
	}
	return appointment.NewRecord(
		appointment.NewRecordId(string(page.ID)),
		notion.Title(page.Properties, RecordTitle),
		status,
		isArchived,
		shared.DateTimePeriod{
			Start: shared.UTCTimeToDateTime(shared.NewUTCTime(period.Start)),
			End:   shared.UTCTimeToDateTime(shared.NewUTCTime(period.End)),
		},
		appointment.NewCustomerId(
			notion.Relations(page.Properties, RecordCustomer)[0].ID.String(),
		),
		appointment.NewServiceId(
			notion.Relations(page.Properties, RecordService)[0].ID.String(),
		),
		notion.CreatedTime(page.Properties, RecordCreatedAt),
	)
}

const (
	BreakTitle  = "Наименование"
	BreakPeriod = "Период"
)

func NotionToWorkBreak(page notionapi.Page) (appointment.WorkBreak, error) {
	const op = "appointment_notion_repository.NotionToWorkBreak"
	period, err := notion.DatePeriod(page.Properties, BreakPeriod)
	if err != nil {
		return appointment.WorkBreak{}, fmt.Errorf("%s: %w", op, err)
	}
	// start := entity.GoTimeToDateTime(period.Start)
	dt := time.Date(
		period.Start.Year(),
		period.Start.Month(),
		period.Start.Day(),
		0, 0, 0, 0, time.Local)
	sb := strings.Builder{}
	sb.WriteString("^\\d (")
	for dt.Before(period.End) {
		sb.WriteString(dt.Format(time.DateOnly))
		sb.WriteByte('|')
		dt = dt.AddDate(0, 0, 1)
	}
	sb.WriteString(dt.Format(time.DateOnly))
	sb.WriteByte(')')
	return appointment.NewWorkBreak(
		appointment.NewWorkBreakId(string(page.ID)),
		notion.Title(page.Properties, BreakTitle),
		sb.String(),
		shared.TimePeriod{
			Start: shared.UTCTimeToTime(shared.NewUTCTime(period.Start)),
			End:   shared.UTCTimeToTime(shared.NewUTCTime(period.End)),
		},
	), nil
}
