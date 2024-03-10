package models

import (
	"time"
)

type ServiceId string

type Service struct {
	Id              ServiceId
	Title           string
	Duration        time.Duration
	Description     string
	CostDescription string
}

type RecordId string

type RecordStatus string

const (
	RecordAwaits RecordStatus = "awaits"
	RecordInWork RecordStatus = "inWork"
)

type Record struct {
	Id     RecordId
	UserId UserId
	Status RecordStatus
}
