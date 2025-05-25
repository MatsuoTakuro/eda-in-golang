package grpc

import (
	"context"

	"github.com/stackus/errors"
	"google.golang.org/grpc"

	"eda-in-golang/modules/stores/storespb"

	"eda-in-golang/modules/baskets/internal/domain"
)

type ProductClient struct {
	client storespb.StoresServiceClient
}

var _ domain.ProductClient = (*ProductClient)(nil)

func NewProductClient(conn *grpc.ClientConn) ProductClient {
	return ProductClient{client: storespb.NewStoresServiceClient(conn)}
}

func (r ProductClient) Find(ctx context.Context, productID string) (*domain.Product, error) {
	resp, err := r.client.GetProduct(ctx, &storespb.GetProductRequest{
		Id: productID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "requesting product")
	}

	return r.productToDomain(resp.Product), nil
}

func (r ProductClient) productToDomain(product *storespb.Product) *domain.Product {
	return &domain.Product{
		ID:      product.GetId(),
		StoreID: product.GetStoreId(),
		Name:    product.GetName(),
		Price:   product.GetPrice(),
	}
}
