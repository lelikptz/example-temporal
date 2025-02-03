package main

import (
	"go.temporal.io/sdk/workflow"
)

const GetOrderStatusQuery = "get_order_status"

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "New"
	OrderStatusCreated    OrderStatus = "Created"
	OrderStatusSend       OrderStatus = "Send"
	OrderStatusFailedSend OrderStatus = "FailedSend"
	OrderStatusCanceled   OrderStatus = "Canceled"
)

type OrderProcessor struct{}

func NewOrderProcessor() *OrderProcessor {
	return &OrderProcessor{}
}

func (o *OrderProcessor) Process(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)

	var orderStatus = OrderStatusNew
	logger.Info("Order process started")

	_ = workflow.SetQueryHandler(ctx, GetOrderStatusQuery, func() (OrderStatus, error) {
		logger.Info("GetOrderStatusQuery")

		return orderStatus, nil
	})

	createOrder(ctx)
	orderStatus = OrderStatusCreated

	isSuccess := sendOrder(ctx)
	if isSuccess {
		orderStatus = OrderStatusSend
		sendNotify(ctx)
	} else {
		orderStatus = OrderStatusFailedSend
		sendCancel(ctx)
		orderStatus = OrderStatusCanceled
	}

	logger.Info("Order process finished")
	return nil
}
