package repo

import (
	"context"
	"encoding/json"
	"log/slog"
	"maps"
	"net/http"
	"sync"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

type HttpProductionCalendar struct {
	log                *logger.Logger
	calendarUrl        string
	client             *http.Client
	mu                 sync.RWMutex
	productionCalendar entity.ProductionCalendar
}

func NewHttpProductionCalendar(
	log *logger.Logger,
	calendarUrl string,
	client *http.Client,
) *HttpProductionCalendar {
	return &HttpProductionCalendar{
		log:                log.With(slog.String("component", "adapters.repo.HttpProductionCalendar")),
		calendarUrl:        calendarUrl,
		client:             client,
		productionCalendar: entity.ProductionCalendar{},
	}
}

func (s *HttpProductionCalendar) updateProductionCalendar(calendar entity.ProductionCalendar) {
	s.mu.Lock()
	defer s.mu.Unlock()
	maps.Copy(s.productionCalendar, calendar)
}

func (s *HttpProductionCalendar) loadProductionCalendar(ctx context.Context) {
	req, err := http.NewRequest("GET", s.calendarUrl, nil)
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
	tmp := entity.ProductionCalendar{}
	err = json.NewDecoder(resp.Body).Decode(&tmp)
	if err != nil {
		s.log.Error(ctx, "failed to decode production calendar", sl.Err(err))
		return
	}
	s.updateProductionCalendar(tmp)
}

func (s *HttpProductionCalendar) Start(ctx context.Context) error {
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

func (s *HttpProductionCalendar) ProductionCalendar(ctx context.Context) (entity.ProductionCalendar, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return maps.Clone(s.productionCalendar), nil
}
