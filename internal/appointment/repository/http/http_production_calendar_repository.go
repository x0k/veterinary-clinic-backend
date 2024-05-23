package appointment_http_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_production_calendar_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/production_calendar"
	production_calendar_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/production_calendar"
)

const productionCalendarRepositoryName = "appointment_http_repository.ProductionCalendarRepository"

type ProductionCalendarRepository struct {
	calendarUrl production_calendar_adapters.Url
	client      *http.Client
}

func NewProductionCalendar(
	calendarUrl production_calendar_adapters.Url,
	client *http.Client,
) *ProductionCalendarRepository {
	return &ProductionCalendarRepository{
		calendarUrl: calendarUrl,
		client:      client,
	}
}

func (s *ProductionCalendarRepository) ProductionCalendar(ctx context.Context) (appointment.ProductionCalendar, error) {
	const op = productionCalendarRepositoryName + ".ProductionCalendar"
	req, err := http.NewRequest("GET", string(s.calendarUrl), nil)
	if err != nil {
		return appointment.ProductionCalendar{}, fmt.Errorf("%s: %w", op, err)
	}
	req = req.WithContext(ctx)
	resp, err := s.client.Do(req)
	if err != nil {
		return appointment.ProductionCalendar{}, fmt.Errorf("%s: %w", op, err)
	}
	productionCalendarDTO := make(appointment_production_calendar_adapters.ProductionCalendarDataDTO)
	if err := json.NewDecoder(resp.Body).Decode(&productionCalendarDTO); err != nil {
		return appointment.ProductionCalendar{}, fmt.Errorf("%s: %w", op, err)
	}
	return appointment.NewProductionCalendar(
		productionCalendarDTO,
	)
}
