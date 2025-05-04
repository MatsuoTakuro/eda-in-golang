package eventhandlers

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/ordering/internal/application/eventhandlers"
	"eda-in-golang/modules/ordering/internal/domain"
)

func SubscribeForNotification(notification eventhandlers.Notification, domainSubscriber ddd.EventSubscriber) {
	domainSubscriber.Subscribe(domain.OrderCreated{}, notification.OnOrderCreated)
	domainSubscriber.Subscribe(domain.OrderReadied{}, notification.OnOrderReadied)
	domainSubscriber.Subscribe(domain.OrderCanceled{}, notification.OnOrderCanceled)
}
