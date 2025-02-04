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
	logger.Info("ORDER PROCESS STARTED")

	_ = workflow.SetQueryHandler(ctx, GetOrderStatusQuery, func() (OrderStatus, error) {
		logger.Info("GetOrderStatusQuery")

		return orderStatus, nil
	})

	err := createOrder(ctx)
	logger.Info("ORDER CREATED")
	if err != nil {
		logger.Error("Create order failed", "error", err)
		return err
	}

	orderStatus = OrderStatusCreated
	isSuccess, err := sendOrder(ctx)
	logger.Info("ORDER SEND")
	if err != nil {
		logger.Error("Send order failed", "error", err)
		return err
	}
	if isSuccess {
		orderStatus = OrderStatusSend
		err = sendNotify(ctx)
		if err != nil {
			logger.Error("Send notification failed", "error", err)
			return err
		}
		logger.Info("NOTIFY SEND")

	} else {
		orderStatus = OrderStatusFailedSend
		err = sendCancel(ctx)
		if err != nil {
			logger.Error("Send cancel failed", "error", err)
			return err
		}
		orderStatus = OrderStatusCanceled
		logger.Info("ORDER CANCELED")
	}

	logger.Info("ORDER PROCESS FINISHED")
	return nil
}
