package main

import (
	"log"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main() {
	temporalClient, err := client.Dial(client.Options{
		HostPort:  client.DefaultHostPort,
		Namespace: "default",
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer temporalClient.Close()

	orderFlowWorker := worker.New(temporalClient, "order-queue", worker.Options{
		DisableRegistrationAliasing: true,
	})
	orderFlowWorker.RegisterWorkflowWithOptions(NewOrderProcessor().Process, workflow.RegisterOptions{
		Name: "OrderWorkflow",
	})

	err = orderFlowWorker.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Worker", err)
	}
}

func createOrder(ctx workflow.Context) error {
	return workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:              "create-order",
		ScheduleToCloseTimeout: time.Minute,
	}), "CreateOrder").Get(ctx, nil)
}

func sendOrder(ctx workflow.Context) (bool, error) {
	var sendOrderResult = false
	err := workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:              "send-order",
		ScheduleToCloseTimeout: time.Minute,
	}), "SendOrder", workflow.GetInfo(ctx).WorkflowExecution.ID).Get(ctx, &sendOrderResult)

	return sendOrderResult, err
}

func sendCancel(ctx workflow.Context) error {
	return workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:              "cancel-order",
		ScheduleToCloseTimeout: time.Minute,
	}), "CancelOrder").Get(ctx, nil)
}

func statusPolling(ctx workflow.Context) (string, error) {
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
