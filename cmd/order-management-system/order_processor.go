package main

import (
	"log"

	"go.temporal.io/sdk/workflow"
)

type OrderProcessor struct{}

func NewOrderProcessor() *OrderProcessor {
	return &OrderProcessor{}
}

func (o *OrderProcessor) Process(ctx workflow.Context) error {
	log.Print("Order process started")

	createOrder(ctx)
	isSuccess := sendOrder(ctx)
	if isSuccess {
		sendNotify(ctx)
	} else {
		sendCancel(ctx)
	}

	log.Print("Order process finished")
	return nil
}
