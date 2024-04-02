package appointment_notion

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

func (s *ServiceRepository) Services(ctx context.Context) ([]appointment.ServiceEntity, error) {
	return s.servicesCache.Load(func() ([]appointment.ServiceEntity, error) {
		r, err := s.client.Database.Query(ctx, s.servicesDatabaseId, nil)
		if err != nil {
			return nil, err
		}
		services := make([]appointment.ServiceEntity, 0, len(r.Results))
		for _, result := range r.Results {
			services = append(services, Service(result))
		}
		return services, nil
	})
}

func (s *ServiceRepository) Service(ctx context.Context, serviceId entity.ServiceId) (appointment.ServiceEntity, error) {
	res, err := s.client.Page.Get(ctx, notionapi.PageID(serviceId))
	if err != nil {
		return appointment.ServiceEntity{}, fmt.Errorf("%w: %s", appointment.ErrServiceLoadFailed, err.Error())
	}
	if res == nil {
		return appointment.ServiceEntity{}, appointment.ErrServiceNotFound
	}
	return Service(*res), nil
}
