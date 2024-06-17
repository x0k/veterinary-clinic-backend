package appointment_js_adapters

import (
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	shared_js_adapters "github.com/x0k/veterinary-clinic-backend/internal/shared/adapters/js"
)

type WorkBreakDTO struct {
	Id              string                           `js:"id"`
	Title           string                           `js:"title"`
	MatchExpression string                           `js:"matchExpression"`
	Period          shared_js_adapters.TimePeriodDTO `js:"period"`
}

func WorkBreakFromDTO(dto WorkBreakDTO) (appointment.WorkBreak, error) {
	return appointment.NewWorkBreak(
		appointment.NewWorkBreakId(dto.Id),
		dto.Title,
		dto.MatchExpression,
		shared_js_adapters.TimePeriodFromDTO(dto.Period),
	), nil
}

func WorkBreakToDTO(workBreak appointment.WorkBreak) (WorkBreakDTO, error) {
	return WorkBreakDTO{
		Id:              workBreak.Id.String(),
		Title:           workBreak.Title,
		MatchExpression: workBreak.MatchExpression,
		Period:          shared_js_adapters.TimePeriodToDTO(workBreak.Period),
	}, nil
}
