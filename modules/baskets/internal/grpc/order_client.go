package grpc

import (
	"context"

	"github.com/stackus/errors"
	"google.golang.org/grpc"

	"eda-in-golang/modules/baskets/internal/domain"
	"eda-in-golang/modules/ordering/orderingpb"
)

type OrderClient struct {
	client orderingpb.OrderingServiceClient
}

var _ domain.OrderClient = (*OrderClient)(nil)

func NewOrderClient(conn *grpc.ClientConn) OrderClient {
	return OrderClient{client: orderingpb.NewOrderingServiceClient(conn)}
}

func (r OrderClient) Save(ctx context.Context, paymentID, customerID string, basketItems map[string]domain.Item) (string, error) {
	items := make([]*orderingpb.Item, 0, len(basketItems))
	for _, item := range basketItems {
		items = append(items, &orderingpb.Item{
			StoreId:     item.StoreID,
			ProductId:   item.ProductID,
			StoreName:   item.StoreName,
			ProductName: item.ProductName,
			Price:       item.ProductPrice,
			Quantity:    int32(item.Quantity),
		})
	}

	resp, err := r.client.CreateOrder(ctx, &orderingpb.CreateOrderRequest{
		Items:      items,
		CustomerId: customerID,
		PaymentId:  paymentID,
	})
	if err != nil {
		return "", errors.Wrap(err, "saving order")
	}

	return resp.GetId(), nil
}
