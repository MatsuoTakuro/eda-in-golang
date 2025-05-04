package grpc

import (
	"context"

	"google.golang.org/grpc"

	"eda-in-golang/modules/notifications/notificationspb"
	"eda-in-golang/modules/ordering/internal/domain/infra"
)

type NotificationClient struct {
	client notificationspb.NotificationsServiceClient
}

var _ infra.NotificationClient = (*NotificationClient)(nil)

func NewNotificationClient(conn *grpc.ClientConn) NotificationClient {
	return NotificationClient{client: notificationspb.NewNotificationsServiceClient(conn)}
}

func (r NotificationClient) NotifyOrderCreated(ctx context.Context, orderID, customerID string) error {
	_, err := r.client.NotifyOrderCreated(ctx, &notificationspb.NotifyOrderCreatedRequest{
		OrderId:    orderID,
		CustomerId: customerID,
	})
	return err
}

func (r NotificationClient) NotifyOrderCanceled(ctx context.Context, orderID, customerID string) error {
	_, err := r.client.NotifyOrderCanceled(ctx, &notificationspb.NotifyOrderCanceledRequest{
		OrderId:    orderID,
		CustomerId: customerID,
	})
	return err
}

func (r NotificationClient) NotifyOrderReady(ctx context.Context, orderID, customerID string) error {
	_, err := r.client.NotifyOrderReady(ctx, &notificationspb.NotifyOrderReadyRequest{
		OrderId:    orderID,
		CustomerId: customerID,
	})
	return err
}
