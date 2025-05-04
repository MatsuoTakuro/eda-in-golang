package grpc

import (
	"context"

	"google.golang.org/grpc"

	"eda-in-golang/modules/customers/customerspb"
	"eda-in-golang/modules/ordering/internal/domain/infra"
)

type CustomerClient struct {
	client customerspb.CustomersServiceClient
}

var _ infra.CustomerClient = (*CustomerClient)(nil)

func NewCustomerClient(conn *grpc.ClientConn) CustomerClient {
	return CustomerClient{client: customerspb.NewCustomersServiceClient(conn)}
}

func (r CustomerClient) Authorize(ctx context.Context, customerID string) error {
	_, err := r.client.AuthorizeCustomer(ctx, &customerspb.AuthorizeCustomerRequest{Id: customerID})
	return err
}
