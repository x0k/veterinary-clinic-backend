package repo

import (
	"context"
	"errors"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
)

var ErrFailedToCreateRecord = errors.New("failed to create record")

type NotionServicesRepo struct {
	servicesDatabaseId notionapi.DatabaseID
	client             *notionapi.Client
}

func NewNotionServices(
	client *notionapi.Client,
	servicesDatabaseId notionapi.DatabaseID,
) *NotionServicesRepo {
	return &NotionServicesRepo{
		client:             client,
		servicesDatabaseId: servicesDatabaseId,
	}
}

func (s *NotionServicesRepo) Services(ctx context.Context) ([]entity.Service, error) {
	r, err := s.client.Database.Query(ctx, s.servicesDatabaseId, nil)
	if err != nil {
		return nil, err
	}
	services := make([]entity.Service, 0, len(r.Results))
	for _, result := range r.Results {
		services = append(services, Service(result))
	}
	return services, nil
}

func (s *NotionServicesRepo) Service(ctx context.Context, serviceId entity.ServiceId) (entity.Service, error) {
	r, err := s.client.Page.Get(ctx, notionapi.PageID(serviceId))
	if err != nil {
		return entity.Service{}, err
	}
	if r == nil {
		return entity.Service{}, usecase.ErrNotFound
	}
	return Service(*r), nil
}
