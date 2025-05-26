package application

import (
	"context"

	"eda-in-golang/modules/payments/internal/domain"
)

type PaymentRepository interface {
	Save(ctx context.Context, payment *domain.Payment) error
	Find(ctx context.Context, paymentID string) (*domain.Payment, error)
}
