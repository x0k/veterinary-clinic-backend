package appointment_notion_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/containers"
)

type ServiceRepository struct {
	client             *notionapi.Client
	servicesDatabaseId notionapi.DatabaseID
	servicesCache      *containers.Expiable[[]appointment.ServiceEntity]
}

func NewService(
	client *notionapi.Client,
	servicesDatabaseId notionapi.DatabaseID,
) *ServiceRepository {
	return &ServiceRepository{
		client:             client,
		servicesDatabaseId: servicesDatabaseId,
		servicesCache:      containers.NewExpiable[[]appointment.ServiceEntity](time.Hour),
	}
}

func (s *ServiceRepository) Start(ctx context.Context) error {
	s.servicesCache.Start(ctx)
	return nil
}

func (s *ServiceRepository) Services(ctx context.Context) ([]appointment.ServiceEntity, error) {
	const op = "appointment_notion.ServiceRepository.Services"
	return s.servicesCache.Load(func() ([]appointment.ServiceEntity, error) {
		r, err := s.client.Database.Query(ctx, s.servicesDatabaseId, nil)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		services := make([]appointment.ServiceEntity, 0, len(r.Results))
		for _, result := range r.Results {
			services = append(services, NotionToService(result))
		}
		return services, nil
	})
}

func (s *ServiceRepository) Service(ctx context.Context, serviceId entity.ServiceId) (appointment.ServiceEntity, error) {
	const op = "appointment_notion.ServiceRepository.Service"
	res, err := s.client.Page.Get(ctx, notionapi.PageID(serviceId))
	if err != nil {
		return appointment.ServiceEntity{}, fmt.Errorf("%s: %w", op, err)
	}
	if res == nil {
		return appointment.ServiceEntity{}, fmt.Errorf("%s: %w", op, entity.ErrNotFound)
	}
	return NotionToService(*res), nil
}
