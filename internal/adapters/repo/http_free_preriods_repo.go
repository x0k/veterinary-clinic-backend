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

var weekdayTimePeriod = entity.TimePeriod{
	Start: entity.Time{
		Hours:   9,
		Minutes: 30,
	},
	End: entity.Time{
		Hours:   17,
		Minutes: 0,
	},
}
var saturdayTimePeriod = entity.TimePeriod{
	Start: weekdayTimePeriod.Start,
	End: entity.Time{
		Hours:   13,
		Minutes: 0,
	},
}
var openingHours = entity.OpeningHours{
	1: weekdayTimePeriod,
	2: weekdayTimePeriod,
	3: weekdayTimePeriod,
	4: weekdayTimePeriod,
	5: weekdayTimePeriod,
	6: saturdayTimePeriod,
}

type HttpFreePeriods struct {
	log                *logger.Logger
	calendarUrl        string
	client             *http.Client
	mu                 sync.RWMutex
	productionCalendar entity.ProductionCalendar
}

func NewHttpFreePeriods(
	log *logger.Logger,
	calendarUrl string,
	client *http.Client,
) *HttpFreePeriods {
	return &HttpFreePeriods{
		log:                log.With(slog.String("component", "adapters.repo.FreePeriodsRepo")),
		calendarUrl:        calendarUrl,
		client:             client,
		productionCalendar: entity.ProductionCalendar{},
	}
}

func (s *HttpFreePeriods) Start(ctx context.Context) error {
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

func (r *HttpFreePeriods) FreePeriods(ctx context.Context, t time.Time) ([]entity.TimePeriod, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	calculator := entity.NewFreePeriodsCalculator(
		openingHours,
		r.productionCalendar,
		entity.GoTimeToDateTime(time.Now()),
	)
	return calculator.Calculate(t)
}

func (r *HttpFreePeriods) NextAvailableDay(ctx context.Context, from time.Time) (time.Time, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	calculator := entity.NewNextAvailableDayCalculator(r.productionCalendar)
	return calculator.Calculate(from), nil
}

func (r *HttpFreePeriods) PrevAvailableDay(ctx context.Context, from time.Time) (*time.Time, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	calculator := entity.NewPrevAvailableDayCalculator(r.productionCalendar, time.Now())
	return calculator.Calculate(from), nil
}

func (s *HttpFreePeriods) updateProductionCalendar(calendar entity.ProductionCalendar) {
	s.mu.Lock()
	defer s.mu.Unlock()
	maps.Copy(s.productionCalendar, calendar)
}

func (s *HttpFreePeriods) loadProductionCalendar(ctx context.Context) {
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
