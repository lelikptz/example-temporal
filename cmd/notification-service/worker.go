package main

import (
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var taskQueueName = "notification"

type Worker struct {
	temporalClient client.Client
	activity       *Notification
}

func NewWorker(temporalClient client.Client, activity *Notification) *Worker {
	return &Worker{temporalClient: temporalClient, activity: activity}
}

func (w *Worker) Run() error {
	temporalWorker := worker.New(w.temporalClient, taskQueueName, worker.Options{})
	temporalWorker.RegisterActivityWithOptions(w.activity.Handle, activity.RegisterOptions{
		Name: "SendNotification",
	})

	return temporalWorker.Run(worker.InterruptCh())
}
