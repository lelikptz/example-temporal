package main

import (
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var taskQueueName = "send-order"

type Worker struct {
	temporalClient client.Client
	activity       *SendOrder
}

func NewWorker(temporalClient client.Client, activity *SendOrder) *Worker {
	return &Worker{temporalClient: temporalClient, activity: activity}
}

func (w *Worker) Run() error {
	temporalWorker := worker.New(w.temporalClient, taskQueueName, worker.Options{})
	temporalWorker.RegisterActivityWithOptions(w.activity.Handle, activity.RegisterOptions{
		Name: "SendOrder",
	})

	return temporalWorker.Run(worker.InterruptCh())
}
