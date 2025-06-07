package internal

import (
	"context"
	"fmt"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/sec"
	"eda-in-golang/modules/cosec/internal/models"
	"eda-in-golang/modules/customers/customerspb"
	"eda-in-golang/modules/depot/depotpb"
	"eda-in-golang/modules/ordering/orderingpb"
	"eda-in-golang/modules/payments/paymentspb"
)

const CreateOrderSagaName = "cosec.CreateOrder"
const CreateOrderReplyChannel = "mallbots.cosec.replies.CreateOrder"

type createOrderSaga struct {
	sec.Saga[*models.CreateOrderData]
}

var (
	_ sec.Saga[*models.CreateOrderData] = (*createOrderSaga)(nil)
)

func NewCreateOrderSaga() sec.Saga[*models.CreateOrderData] {
	saga := createOrderSaga{
		Saga: sec.NewSaga[*models.CreateOrderData](CreateOrderSagaName, CreateOrderReplyChannel),
	}

	// 0. RejectOrder
	saga.AddStep().
		Compensation(saga.rejectOrder)

	// 1. AuthorizeCustomer
	saga.AddStep().
		Normal(saga.authorizeCustomer)

	// 2. CreateShoppingList, -CancelShoppingList
	saga.AddStep().
		Normal(saga.createShoppingList).
		OnNormalReply(depotpb.CreatedShoppingListReply, saga.onCreatedShoppingListReply).
		Compensation(saga.cancelShoppingList)

	// 3. ConfirmPayment
	saga.AddStep().
		Normal(saga.confirmPayment)

	// 4. InitiateShopping
	saga.AddStep().
		Normal(saga.initiateShopping)

	// 5. ApproveOrder
	saga.AddStep().
		Normal(saga.approveOrder)

	return saga
}

func (s createOrderSaga) rejectOrder(ctx context.Context, data *models.CreateOrderData) am.Command {
	return am.NewCommand(
		orderingpb.RejectOrderCommand,
		orderingpb.CommandChannel,
		&orderingpb.RejectOrder{Id: data.OrderID},
	)
}

func (s createOrderSaga) authorizeCustomer(ctx context.Context, data *models.CreateOrderData) am.Command {
	return am.NewCommand(
		customerspb.AuthorizeCustomerCommand,
		customerspb.CommandChannel,
		&customerspb.AuthorizeCustomer{Id: data.CustomerID},
	)
}

func (s createOrderSaga) createShoppingList(ctx context.Context, data *models.CreateOrderData) am.Command {
	items := make([]*depotpb.CreateShoppingList_Item, len(data.Items))
	for i, item := range data.Items {
		items[i] = &depotpb.CreateShoppingList_Item{
			ProductId: item.ProductID,
			StoreId:   item.StoreID,
			Quantity:  int32(item.Quantity),
		}
	}

	return am.NewCommand(
		depotpb.CreateShoppingListCommand,
		depotpb.CommandChannel,
		&depotpb.CreateShoppingList{
			OrderId: data.OrderID,
			Items:   items,
		},
	)
}

func (s createOrderSaga) onCreatedShoppingListReply(ctx context.Context, data *models.CreateOrderData, reply ddd.Reply) error {
	payload, ok := reply.Payload().(*depotpb.CreatedShoppingList)
	if !ok {
		return fmt.Errorf(
			"expected reply payload of type depotpb.CreatedShoppingList, got %T",
			reply.Payload(),
		)
	}

	data.ShoppingID = payload.GetId()

	return nil
}

func (s createOrderSaga) cancelShoppingList(ctx context.Context, data *models.CreateOrderData) am.Command {
	return am.NewCommand(depotpb.CancelShoppingListCommand,
		depotpb.CommandChannel,
		&depotpb.CancelShoppingList{Id: data.ShoppingID},
	)
}

func (s createOrderSaga) confirmPayment(ctx context.Context, data *models.CreateOrderData) am.Command {
	return am.NewCommand(
		paymentspb.ConfirmPaymentCommand,
		paymentspb.CommandChannel,
		&paymentspb.ConfirmPayment{
			Id:     data.PaymentID,
			Amount: data.Total,
		},
	)
}

func (s createOrderSaga) initiateShopping(ctx context.Context, data *models.CreateOrderData) am.Command {
	return am.NewCommand(
		depotpb.InitiateShoppingCommand,
		depotpb.CommandChannel,
		&depotpb.InitiateShopping{Id: data.ShoppingID},
	)
}

func (s createOrderSaga) approveOrder(ctx context.Context, data *models.CreateOrderData) am.Command {
	return am.NewCommand(
		orderingpb.ApproveOrderCommand,
		orderingpb.CommandChannel,
		&orderingpb.ApproveOrder{
			Id:         data.OrderID,
			ShoppingId: data.ShoppingID,
		},
	)
}
