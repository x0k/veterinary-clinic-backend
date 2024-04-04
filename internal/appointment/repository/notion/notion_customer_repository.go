package appointment_notion_repository

import (
	"context"
	"fmt"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

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
	const op = "appointment_notion.CustomerRepository.Customer"
	res, err := r.client.Page.Get(ctx, notionapi.PageID(id))
	if err != nil {
		return appointment.CustomerEntity{}, fmt.Errorf("%s: %w", op, err)
	}
	if res == nil {
		return appointment.CustomerEntity{}, fmt.Errorf("%s: %w", op, entity.ErrNotFound)
	}
	return NotionToCustomer(*res), nil
}
