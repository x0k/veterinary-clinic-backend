package appointment

import (
	"github.com/google/uuid"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type ServiceId uuid.UUID

func (s ServiceId) String() string {
	return uuid.UUID(s).String()
}

type Service struct {
	Id                ServiceId
	Title             string
	DurationInMinutes entity.DurationInMinutes
	Description       string
	CostDescription   string
}
