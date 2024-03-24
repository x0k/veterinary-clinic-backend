package entity

import "errors"

var ErrInvalidRecordStatus = errors.New("invalid record status")
var ErrInvalidDate = errors.New("invalid date")

type RecordId string

type RecordStatus string

const (
	RecordAwaits            RecordStatus = "awaits"
	RecordDone              RecordStatus = "done"
	RecordNotAppear         RecordStatus = "failed"
	RecordDoneArchived      RecordStatus = "done_archived"
	RecordNotAppearArchived RecordStatus = "failed_archived"
)

type Record struct {
	Id             RecordId
	Status         RecordStatus
	DateTimePeriod DateTimePeriod
	User           User
	Service        Service
}

func RecordStatusName(status RecordStatus) (string, error) {
	switch status {
	case RecordAwaits:
		return "ожидает", nil
	case RecordDone:
		return "выполнено", nil
	case RecordNotAppear:
		return "не пришел", nil
	case RecordDoneArchived:
		return "архив (выполнено)", nil
	case RecordNotAppearArchived:
		return "архив (не пришел)", nil
	default:
		return "", ErrInvalidRecordStatus
	}
}
