package application

import (
	"context"

	"eda-in-golang/modules/search/internal/models"
)

type CustomerRepository interface {
	Find(ctx context.Context, customerID string) (*models.Customer, error)
}
