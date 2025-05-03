package application

import (
	"context"

	"eda-in-golang/modules/notifications/internal/models"
)

type CustomerRepository interface {
	Find(ctx context.Context, customerID string) (*models.Customer, error)
}
