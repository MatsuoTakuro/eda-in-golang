package infra

import (
	"context"
)

type NotificationClient interface {
	NotifyOrderCreated(ctx context.Context, orderID, customerID string) error
	NotifyOrderCanceled(ctx context.Context, orderID, customerID string) error
	NotifyOrderReady(ctx context.Context, orderID, customerID string) error
}
