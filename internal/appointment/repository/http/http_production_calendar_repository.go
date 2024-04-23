package appointment_http_repository

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_production_calendar_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/production_calendar"
	production_calendar_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/production_calendar"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

const productionCalendarRepositoryName = "appointment_http_repository.ProductionCalendarRepository"

type ProductionCalendarRepository struct {
	log                *logger.Logger
	calendarUrl        production_calendar_adapters.Url
	client             *http.Client
	mu                 sync.RWMutex
	productionCalendar appointment.ProductionCalendar
}

func NewProductionCalendar(
	log *logger.Logger,
	calendarUrl production_calendar_adapters.Url,
	client *http.Client,
) *ProductionCalendarRepository {
	return &ProductionCalendarRepository{
		log:         log.With(slog.String("component", productionCalendarRepositoryName)),
		calendarUrl: calendarUrl,
		client:      client,
		productionCalendar: appointment.NewProductionCalendar(
			make(appointment.ProductionCalendarData),
		),
	}
}

func (s *ProductionCalendarRepository) Name() string {
	return productionCalendarRepositoryName
}

func (s *ProductionCalendarRepository) Start(ctx context.Context) error {
	ticker := time.NewTicker(time.Hour * 24)
	defer ticker.Stop()
	s.loadProductionCalendar(ctx)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			s.loadProductionCalendar(ctx)
		}
	}
}

func (s *ProductionCalendarRepository) ProductionCalendar(ctx context.Context) (appointment.ProductionCalendar, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.productionCalendar.Clone(), nil
}

func (s *ProductionCalendarRepository) updateProductionCalendar(data appointment.ProductionCalendarData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.productionCalendar.Update(data)
}

func (s *ProductionCalendarRepository) loadProductionCalendar(ctx context.Context) {
	req, err := http.NewRequest("GET", string(s.calendarUrl), nil)
	if err != nil {
		s.log.Error(ctx, "failed to create production calendar data request", sl.Err(err))
		return
	}
	req = req.WithContext(ctx)
	resp, err := s.client.Do(req)
	if err != nil {
		s.log.Error(ctx, "failed to load production calendar data", sl.Err(err))
		return
	}
	tmp := make(appointment_production_calendar_adapters.ProductionCalendarDataDTO)
	err = json.NewDecoder(resp.Body).Decode(&tmp)
	if err != nil {
		s.log.Error(ctx, "failed to decode production calendar data", sl.Err(err))
		return
	}
	productionCalendarData, err := appointment.NewProductionCalendarData(tmp)
	if err != nil {
		s.log.Error(ctx, "failed to validate production calendar data", sl.Err(err))
		return
	}
	s.updateProductionCalendar(productionCalendarData)
}
