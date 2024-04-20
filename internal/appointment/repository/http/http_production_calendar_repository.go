package appointment_http_repository

import (
	"context"
	"encoding/json"
	"log/slog"
	"maps"
	"net/http"
	"sync"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	production_calendar_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/production_calendar"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/mapx"
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
		log:                log.With(slog.String("component", productionCalendarRepositoryName)),
		calendarUrl:        calendarUrl,
		client:             client,
		productionCalendar: appointment.NewProductionCalendar(),
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
	return mapx.Clone(s.productionCalendar), nil
}

func (s *ProductionCalendarRepository) updateProductionCalendar(calendar appointment.ProductionCalendar) {
	s.mu.Lock()
	defer s.mu.Unlock()
	maps.Copy(s.productionCalendar, calendar)
}

func (s *ProductionCalendarRepository) loadProductionCalendar(ctx context.Context) {
	req, err := http.NewRequest("GET", string(s.calendarUrl), nil)
	if err != nil {
		s.log.Error(ctx, "failed to load production calendar", sl.Err(err))
		return
	}
	req = req.WithContext(ctx)
	resp, err := s.client.Do(req)
	if err != nil {
		s.log.Error(ctx, "failed to load production calendar", sl.Err(err))
		return
	}
	tmp := appointment.NewProductionCalendar()
	err = json.NewDecoder(resp.Body).Decode(&tmp)
	if err != nil {
		s.log.Error(ctx, "failed to decode production calendar", sl.Err(err))
		return
	}
	s.updateProductionCalendar(tmp)
}
