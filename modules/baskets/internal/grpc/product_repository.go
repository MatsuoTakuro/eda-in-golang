package grpc

import (
	"context"

	"github.com/stackus/errors"
	"google.golang.org/grpc"

	"eda-in-golang/modules/stores/storespb"

	"eda-in-golang/modules/baskets/internal/domain"
)

type ProductRepository struct {
	client storespb.StoresServiceClient
}

var _ domain.ProductRepository = (*ProductRepository)(nil)

func NewProductRepository(conn *grpc.ClientConn) ProductRepository {
	return ProductRepository{client: storespb.NewStoresServiceClient(conn)}
}

func (r ProductRepository) Find(ctx context.Context, productID string) (*domain.Product, error) {
	resp, err := r.client.GetProduct(ctx, &storespb.GetProductRequest{
		Id: productID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "requesting product")
	}

	return r.productToDomain(resp.Product), nil
}

func (r ProductRepository) productToDomain(product *storespb.Product) *domain.Product {
	return &domain.Product{
		ID:      product.GetId(),
		StoreID: product.GetStoreId(),
		Name:    product.GetName(),
		Price:   product.GetPrice(),
	}
}
