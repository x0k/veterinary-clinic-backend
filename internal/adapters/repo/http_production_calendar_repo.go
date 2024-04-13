package repo

import (
	"context"
	"encoding/json"
	"log/slog"
	"maps"
	"net/http"
	"sync"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

type HttpProductionCalendarRepo struct {
	log                *logger.Logger
	calendarUrl        adapters.ProductionCalendarUrl
	client             *http.Client
	mu                 sync.RWMutex
	productionCalendar shared.ProductionCalendar
}

func NewHttpProductionCalendar(
	log *logger.Logger,
	calendarUrl adapters.ProductionCalendarUrl,
	client *http.Client,
) *HttpProductionCalendarRepo {
	return &HttpProductionCalendarRepo{
		log:                log.With(slog.String("component", "adapters.repo.FreePeriodsRepo")),
		calendarUrl:        calendarUrl,
		client:             client,
		productionCalendar: shared.ProductionCalendar{},
	}
}

func (s *HttpProductionCalendarRepo) Start(ctx context.Context) error {
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

func (s *HttpProductionCalendarRepo) ProductionCalendar(ctx context.Context) (shared.ProductionCalendar, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return maps.Clone(s.productionCalendar), nil
}

func (s *HttpProductionCalendarRepo) updateProductionCalendar(calendar shared.ProductionCalendar) {
	s.mu.Lock()
	defer s.mu.Unlock()
	maps.Copy(s.productionCalendar, calendar)
}

func (s *HttpProductionCalendarRepo) loadProductionCalendar(ctx context.Context) {
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
	tmp := shared.ProductionCalendar{}
	err = json.NewDecoder(resp.Body).Decode(&tmp)
	if err != nil {
		s.log.Error(ctx, "failed to decode production calendar", sl.Err(err))
		return
	}
	s.updateProductionCalendar(tmp)
}
