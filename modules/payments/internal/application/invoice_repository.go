package application

import (
	"context"

	"eda-in-golang/modules/payments/internal/domain"
)

type InvoiceRepository interface {
	Find(ctx context.Context, invoiceID string) (*domain.Invoice, error)
	Save(ctx context.Context, invoice *domain.Invoice) error
	Update(ctx context.Context, invoice *domain.Invoice) error
}
