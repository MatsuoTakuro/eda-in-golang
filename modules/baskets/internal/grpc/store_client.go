package grpc

import (
	"context"

	"google.golang.org/grpc"

	"eda-in-golang/modules/baskets/internal/domain"
	"eda-in-golang/modules/stores/storespb"
)

type StoreClient struct {
	client storespb.StoresServiceClient
}

var _ domain.StoreClient = (*StoreClient)(nil)

func NewStoreClient(conn *grpc.ClientConn) StoreClient {
	return StoreClient{client: storespb.NewStoresServiceClient(conn)}
}

func (r StoreClient) Find(ctx context.Context, storeID string) (*domain.Store, error) {
	resp, err := r.client.GetStore(ctx, &storespb.GetStoreRequest{
		Id: storeID,
	})
	if err != nil {
		return nil, err
	}

	return r.storeToDomain(resp.Store), nil
}

func (r StoreClient) storeToDomain(store *storespb.Store) *domain.Store {
	return &domain.Store{
		ID:   store.GetId(),
		Name: store.GetName(),
	}
}
