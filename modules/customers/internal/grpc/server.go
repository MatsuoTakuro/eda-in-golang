package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"eda-in-golang/modules/customers/customerspb"
	"eda-in-golang/modules/customers/internal/application"
	"eda-in-golang/modules/customers/internal/domain"
)

type server struct {
	app application.App
	customerspb.UnimplementedCustomersServiceServer
}

var _ customerspb.CustomersServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	customerspb.RegisterCustomersServiceServer(registrar, server{
		app: app,
	})
	return nil
}

func (s server) RegisterCustomer(ctx context.Context, request *customerspb.RegisterCustomerRequest) (resp *customerspb.RegisterCustomerResponse, err error) {
	id := uuid.New().String()
	err = s.app.RegisterCustomer(ctx, application.RegisterCustomer{
		ID:        id,
		Name:      request.GetName(),
		SmsNumber: request.GetSmsNumber(),
	})

	return &customerspb.RegisterCustomerResponse{Id: id}, err
}

func (s server) AuthorizeCustomer(ctx context.Context, request *customerspb.AuthorizeCustomerRequest) (resp *customerspb.AuthorizeCustomerResponse, err error) {
	err = s.app.AuthorizeCustomer(ctx, application.AuthorizeCustomer{
		ID: request.GetId(),
	})

	return &customerspb.AuthorizeCustomerResponse{}, err
}

func (s server) GetCustomer(ctx context.Context, request *customerspb.GetCustomerRequest) (resp *customerspb.GetCustomerResponse, err error) {
	customer, err := s.app.GetCustomer(ctx, application.GetCustomer{
		ID: request.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &customerspb.GetCustomerResponse{
		Customer: s.customerFromDomain(customer),
	}, nil
}

func (s server) EnableCustomer(ctx context.Context, request *customerspb.EnableCustomerRequest) (resp *customerspb.EnableCustomerResponse, err error) {
	err = s.app.EnableCustomer(ctx, application.EnableCustomer{ID: request.GetId()})
	return &customerspb.EnableCustomerResponse{}, err
}

func (s server) DisableCustomer(ctx context.Context, request *customerspb.DisableCustomerRequest) (resp *customerspb.DisableCustomerResponse, err error) {
	err = s.app.DisableCustomer(ctx, application.DisableCustomer{ID: request.GetId()})
	return &customerspb.DisableCustomerResponse{}, err
}

func (s server) customerFromDomain(customer *domain.Customer) *customerspb.Customer {
	return &customerspb.Customer{
		Id:        customer.ID(),
		Name:      customer.Name,
		SmsNumber: customer.SmsNumber,
		Enabled:   customer.Enabled,
	}
}
