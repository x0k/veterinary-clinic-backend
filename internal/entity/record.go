package entity

import "errors"

var ErrInvalidRecordStatus = errors.New("invalid record status")

type RecordId string

type RecordStatus string

const (
	RecordAwaits RecordStatus = "awaits"
	RecordInWork RecordStatus = "inWork"
)

type Record struct {
	Id             RecordId
	Status         RecordStatus
	DateTimePeriod DateTimePeriod
	UserId         *UserId
	Service        Service
}

func RecordStatusName(status RecordStatus) (string, error) {
	switch status {
	case RecordAwaits:
		return "ожидает", nil
	case RecordInWork:
		return "в работе", nil
	default:
		return "", ErrInvalidRecordStatus
	}
}
