package appointment_js_adapters

import "github.com/x0k/veterinary-clinic-backend/internal/appointment"

type ProductionCalendarDTO map[string]int

func ProductionCalendarToDTO(productionCalendar appointment.ProductionCalendar) (ProductionCalendarDTO, error) {
	return productionCalendar.ToDTO(), nil
}

func ProductionCalendarFromDTO(dto ProductionCalendarDTO) (appointment.ProductionCalendar, error) {
	return appointment.NewProductionCalendar(dto)
}
