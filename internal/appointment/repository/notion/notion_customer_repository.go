package appointment_notion_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/notion"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
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

func (r *CustomerRepository) CustomerByIdentity(ctx context.Context, identity appointment.CustomerIdentity) (appointment.CustomerEntity, error) {
	const op = customerRepositoryName + ".Customer"
	res, err := r.client.Database.Query(ctx, r.customersDatabaseId, &notionapi.DatabaseQueryRequest{
		Filter: notionapi.PropertyFilter{
			Property: CustomerUserId,
			RichText: &notionapi.TextFilterCondition{
				Equals: identity.String(),
			},
		},
	})
	if err != nil {
		return appointment.CustomerEntity{}, fmt.Errorf("%s: %w", op, err)
	}
	if res == nil || len(res.Results) == 0 {
		return appointment.CustomerEntity{}, fmt.Errorf("%s: %w", op, shared.ErrNotFound)
	}
	return NotionToCustomer(res.Results[0]), nil
}

func (s *CustomerRepository) CustomerById(ctx context.Context, customerId appointment.CustomerId) (appointment.CustomerEntity, error) {
	const op = customerRepositoryName + ".CustomerById"
	res, err := s.client.Page.Get(ctx, notionapi.PageID(customerId.String()))
	if err != nil {
		return appointment.CustomerEntity{}, fmt.Errorf("%s: %w", op, err)
	}
	if res == nil {
		return appointment.CustomerEntity{}, fmt.Errorf("%s: %w", op, shared.ErrNotFound)
	}
	return NotionToCustomer(*res), nil
}

func (r *CustomerRepository) CreateCustomer(ctx context.Context, customer *appointment.CustomerEntity) error {
	const op = customerRepositoryName + ".CreateCustomer"
	if _, err := r.CustomerByIdentity(ctx, customer.Identity); !errors.Is(err, shared.ErrNotFound) {
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return fmt.Errorf("%s: %w", op, shared.ErrAlreadyExists)
	}
	properties := notionapi.Properties{
		CustomerTitle: &notionapi.TitleProperty{
			Type:  notionapi.PropertyTypeTitle,
			Title: notion.ToRichText(customer.Name),
		},
		CustomerEmail: &notionapi.EmailProperty{
			Type:  notionapi.PropertyTypeEmail,
			Email: customer.Email,
		},
		CustomerPhoneNumber: &notionapi.PhoneNumberProperty{
			Type:        notionapi.PropertyTypePhoneNumber,
			PhoneNumber: customer.PhoneNumber,
		},
		CustomerUserId: &notionapi.RichTextProperty{
			Type:     notionapi.PropertyTypeRichText,
			RichText: notion.ToRichText(customer.Identity.String()),
		},
	}
	res, err := r.client.Page.Create(ctx, &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: r.customersDatabaseId,
		},
		Properties: properties,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return customer.SetId(appointment.NewCustomerId(res.ID.String()))
}
