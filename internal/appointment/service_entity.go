package appointment

import (
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type ServiceId string

func (s ServiceId) String() string {
	return string(s)
}

type Service struct {
	Id                ServiceId
	Title             string
	DurationInMinutes entity.DurationInMinutes
	Description       string
	CostDescription   string
}
