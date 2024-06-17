package appointment_notion_repository

import (
	"context"
	"fmt"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

type ServicesRepository struct {
	client             *notionapi.Client
	servicesDatabaseId notionapi.DatabaseID
}

func NewServices(
	client *notionapi.Client,
	servicesDatabaseId notionapi.DatabaseID,
) *ServicesRepository {
	return &ServicesRepository{
		client:             client,
		servicesDatabaseId: servicesDatabaseId,
	}
}

func (s *ServicesRepository) Services(ctx context.Context) ([]appointment.ServiceEntity, error) {
	const op = appointmentRepositoryName + ".Services"
	r, err := s.client.Database.Query(ctx, s.servicesDatabaseId, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	services := make([]appointment.ServiceEntity, 0, len(r.Results))
	for _, result := range r.Results {
		services = append(services, NotionToService(result))
	}
	return services, nil
}

func (s *ServicesRepository) Service(ctx context.Context, serviceId appointment.ServiceId) (appointment.ServiceEntity, error) {
	const op = appointmentRepositoryName + ".Service"
	res, err := s.client.Page.Get(ctx, notionapi.PageID(serviceId))
	if err != nil {
		return appointment.ServiceEntity{}, fmt.Errorf("%s: %w", op, err)
	}
	if res == nil {
		return appointment.ServiceEntity{}, fmt.Errorf("%s: %w", op, shared.ErrNotFound)
	}
	return NotionToService(*res), nil
}
