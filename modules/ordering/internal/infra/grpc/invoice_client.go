package grpc

import (
	"context"

	"google.golang.org/grpc"

	"eda-in-golang/modules/ordering/internal/domain/infra"
	"eda-in-golang/modules/payments/paymentspb"
)

type InvoiceClient struct {
	client paymentspb.PaymentsServiceClient
}

var _ infra.InvoiceClient = (*InvoiceClient)(nil)

func NewInvoiceClient(conn *grpc.ClientConn) InvoiceClient {
	return InvoiceClient{client: paymentspb.NewPaymentsServiceClient(conn)}
}

func (r InvoiceClient) Save(ctx context.Context, orderID, paymentID string, amount float64) error {
	_, err := r.client.CreateInvoice(ctx, &paymentspb.CreateInvoiceRequest{
		OrderId:   orderID,
		PaymentId: paymentID,
		Amount:    amount,
	})
	return err
}

func (r InvoiceClient) Delete(ctx context.Context, invoiceID string) error {
	_, err := r.client.CancelInvoice(ctx, &paymentspb.CancelInvoiceRequest{Id: invoiceID})
	return err
}
