package main

import (
	"time"

	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"
)

const GetOrderStatusQuery = "get_order_status"

type OrderStatus string

const (
	OrderStatusNew      OrderStatus = "New"
	OrderStatusFinished OrderStatus = "Finished"
	OrderStatusCanceled OrderStatus = "Canceled"
)

type Order struct {
	ReservationTime time.Time
}

type OrderProcessor struct {
	logger log.Logger
}

func NewOrderProcessor() *OrderProcessor {
	return &OrderProcessor{}
}

func (o *OrderProcessor) Process(ctx workflow.Context, order Order) error {
	o.logger = workflow.GetLogger(ctx)
	o.logger.Info("ORDER PROCESS STARTED")

	var orderStatus = OrderStatusNew
	_ = workflow.SetQueryHandler(ctx, GetOrderStatusQuery, func() (OrderStatus, error) {
		return orderStatus, nil
	})

	wg := workflow.NewWaitGroup(ctx)
	workflow.Go(ctx, func(ctx workflow.Context) {
		selector := workflow.NewSelector(ctx)
		selector.AddReceive(workflow.GetSignalChannel(ctx, "cancel-order-signal"), func(c workflow.ReceiveChannel, more bool) {
			c.Receive(ctx, nil)
			orderStatus = OrderStatusCanceled
			wg.Add(1)
			_ = o.sendCancel(ctx)
			o.logger.Info("SEND CANCEL FINISHED")
			wg.Done()
			return
		})
		selector.Select(ctx)
	})
	defer wg.Wait(ctx)

	if orderStatus == OrderStatusCanceled {
		return nil
	}
	err := o.createOrder(ctx)
	if err != nil {
		return err
	}

	if orderStatus == OrderStatusCanceled {
		return nil
	}
	_ = workflow.Sleep(ctx, order.ReservationTime.Sub(workflow.Now(ctx))-10*time.Second)

	if orderStatus == OrderStatusCanceled {
		return nil
	}
	isSuccess, err := o.sendOrder(ctx)
	if err != nil {
		return err
	}
	if orderStatus == OrderStatusCanceled {
		return nil
	}
	if !isSuccess {
		err = o.sendCancel(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	for orderStatus != OrderStatusFinished && orderStatus != OrderStatusCanceled {
		newStatus, err := o.statusPolling(ctx)
		if err != nil {
			return err
		}
		orderStatus = OrderStatus(newStatus)
	}

	o.logger.Info("ORDER PROCESS FINISHED")
	return nil
}

func (o *OrderProcessor) createOrder(ctx workflow.Context) error {
	o.logger.Info("ORDER CREATED")
	return workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:              "create-order",
		ScheduleToCloseTimeout: time.Minute,
	}), "CreateOrder").Get(ctx, nil)
}

func (o *OrderProcessor) sendOrder(ctx workflow.Context) (bool, error) {
	var sendOrderResult = false
	err := workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:              "send-order",
		ScheduleToCloseTimeout: time.Minute,
	}), "SendOrder", workflow.GetInfo(ctx).WorkflowExecution.ID).Get(ctx, &sendOrderResult)

	return sendOrderResult, err
}

func (o *OrderProcessor) sendCancel(ctx workflow.Context) error {
	return workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:              "cancel-order",
		ScheduleToCloseTimeout: time.Minute,
	}), "CancelOrder").Get(ctx, nil)
}

func (o *OrderProcessor) statusPolling(ctx workflow.Context) (string, error) {
	var status string
	err := workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:              "status-polling",
		ScheduleToCloseTimeout: time.Minute,
	}), "StatusPoll").Get(ctx, &status)
	if err != nil {
		return "", err
	}

	return status, nil
}
