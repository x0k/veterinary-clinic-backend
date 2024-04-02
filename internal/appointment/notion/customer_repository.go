package appointment_notion

import "github.com/jomei/notionapi"

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
