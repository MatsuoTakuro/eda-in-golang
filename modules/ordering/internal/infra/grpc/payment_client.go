package grpc

import (
	"context"

	"google.golang.org/grpc"

	"eda-in-golang/modules/ordering/internal/domain/infra"
	"eda-in-golang/modules/payments/paymentspb"
)

type PaymentClient struct {
	client paymentspb.PaymentsServiceClient
}

var _ infra.PaymentClient = (*PaymentClient)(nil)

func NewPaymentClient(conn *grpc.ClientConn) PaymentClient {
	return PaymentClient{
		client: paymentspb.NewPaymentsServiceClient(conn),
	}
}

func (r PaymentClient) Confirm(ctx context.Context, paymentID string) error {
	_, err := r.client.ConfirmPayment(ctx, &paymentspb.ConfirmPaymentRequest{Id: paymentID})
	return err
}
