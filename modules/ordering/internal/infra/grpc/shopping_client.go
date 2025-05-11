package grpc

import (
	"context"

	"google.golang.org/grpc"

	"eda-in-golang/modules/depot/depotpb"
	"eda-in-golang/modules/ordering/internal/domain"
	"eda-in-golang/modules/ordering/internal/domain/infra"
)

type ShoppingClient struct {
	client depotpb.DepotServiceClient
}

var _ infra.ShoppingClient = (*ShoppingClient)(nil)

func NewShoppingListClient(conn *grpc.ClientConn) ShoppingClient {
	return ShoppingClient{client: depotpb.NewDepotServiceClient(conn)}
}

func (c ShoppingClient) Create(ctx context.Context, orderID string, orderItems []domain.Item) (string, error) {
	items := make([]*depotpb.OrderItem, len(orderItems))
	for i, item := range orderItems {
		items[i] = c.itemFromDomain(item)
	}

	response, err := c.client.CreateShoppingList(ctx, &depotpb.CreateShoppingListRequest{
		OrderId: orderID,
		Items:   items,
	})
	if err != nil {
		return "", err
	}

	return response.GetId(), nil
}

func (c ShoppingClient) Cancel(ctx context.Context, shoppingID string) error {
	_, err := c.client.CancelShoppingList(ctx, &depotpb.CancelShoppingListRequest{Id: shoppingID})
	return err
}

func (c ShoppingClient) itemFromDomain(item domain.Item) *depotpb.OrderItem {
	return &depotpb.OrderItem{
		ProductId: item.ProductID,
		StoreId:   item.StoreID,
		Quantity:  int32(item.Quantity),
	}
}
