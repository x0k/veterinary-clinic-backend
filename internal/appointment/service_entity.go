package appointment

import (
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type ServiceId string

func NewServiceId(id string) ServiceId {
	return ServiceId(id)
}

func (s ServiceId) String() string {
	return string(s)
}

type ServiceEntity struct {
	Id                ServiceId
	Title             string
	DurationInMinutes entity.DurationInMinutes
	Description       string
	CostDescription   string
}

func NewService(
	id ServiceId,
	title string,
	durationInMinutes entity.DurationInMinutes,
	description string,
	costDescription string,
) ServiceEntity {
	return ServiceEntity{
		Id:                id,
		Title:             title,
		DurationInMinutes: durationInMinutes,
		Description:       description,
		CostDescription:   costDescription,
	}
}
