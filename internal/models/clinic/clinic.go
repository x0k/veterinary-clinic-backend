package clinic

import (
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/user"
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
	Awaits RecordStatus = "awaits"
	InWork RecordStatus = "inWork"
)

type Record struct {
	Id     RecordId
	UserId user.Id
	Status RecordStatus
}
