package appointment_js_adapters

import (
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

type ServiceDTO struct {
	Id                string `js:"id"`
	Title             string `js:"title"`
	DurationInMinutes int    `js:"durationInMinutes"`
	Description       string `js:"description"`
	CostDescription   string `js:"costDescription"`
}

func ServiceIdToDTO(id appointment.ServiceId) (string, error) {
	return id.String(), nil
}

func ServiceToDTO(service appointment.ServiceEntity) (ServiceDTO, error) {
	return ServiceDTO{
		Id:                service.Id.String(),
		Title:             service.Title,
		DurationInMinutes: service.DurationInMinutes.Int(),
		Description:       service.Description,
		CostDescription:   service.CostDescription,
	}, nil
}

func ServiceFromDTO(dto ServiceDTO) (appointment.ServiceEntity, error) {
	return appointment.NewService(
		appointment.NewServiceId(dto.Id),
		dto.Title,
		shared.NewDurationInMinutes(dto.DurationInMinutes),
		dto.Description,
		dto.CostDescription,
	), nil
}
