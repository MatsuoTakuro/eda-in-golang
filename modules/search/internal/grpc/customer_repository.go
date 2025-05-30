package grpc

import (
	"context"

	"google.golang.org/grpc"

	"eda-in-golang/modules/customers/customerspb"
	"eda-in-golang/modules/search/internal/application"
	"eda-in-golang/modules/search/internal/models"
)

type CustomerRepository struct {
	client customerspb.CustomersServiceClient
}

var _ application.CustomerRepository = (*CustomerRepository)(nil)

func NewCustomerRepository(conn *grpc.ClientConn) CustomerRepository {
	return CustomerRepository{
		client: customerspb.NewCustomersServiceClient(conn),
	}
}

func (r CustomerRepository) Find(ctx context.Context, customerID string) (*models.Customer, error) {
	resp, err := r.client.GetCustomer(ctx, &customerspb.GetCustomerRequest{Id: customerID})
	if err != nil {
		return nil, err
	}

	return &models.Customer{
		ID:   resp.GetCustomer().GetId(),
		Name: resp.GetCustomer().GetName(),
	}, nil
}
