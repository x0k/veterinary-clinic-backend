package appointment_notion_repository

import (
	"context"
	"fmt"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

const customerRepositoryName = "appointment_notion_repository.CustomerRepository"

type CustomerRepository struct {
	client              *notionapi.Client
	customersDatabaseId notionapi.DatabaseID
}

func NewCustomer(
	client *notionapi.Client,
	customersDatabaseId notionapi.DatabaseID,
) *CustomerRepository {
	return &CustomerRepository{
		client:              client,
		customersDatabaseId: customersDatabaseId,
	}
}

func (r *CustomerRepository) Customer(ctx context.Context, id appointment.CustomerId) (appointment.CustomerEntity, error) {
	const op = customerRepositoryName + ".Customer"
	res, err := r.client.Database.Query(ctx, r.customersDatabaseId, &notionapi.DatabaseQueryRequest{
		Filter: notionapi.PropertyFilter{
			Property: CustomerUserId,
			RichText: &notionapi.TextFilterCondition{
				Equals: id.String(),
			},
		},
	})
	if err != nil {
		return appointment.CustomerEntity{}, fmt.Errorf("%s: %w", op, err)
	}
	if res == nil || len(res.Results) == 0 {
		return appointment.CustomerEntity{}, fmt.Errorf("%s: %w", op, entity.ErrNotFound)
	}
	return NotionToCustomer(res.Results[0]), nil
}
