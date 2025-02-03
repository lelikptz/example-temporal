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

func createOrder(ctx workflow.Context) {
	_ = workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:              "create-order",
		ScheduleToCloseTimeout: time.Minute,
	}), "CreateOrder").Get(ctx, nil)
}

func sendOrder(ctx workflow.Context) bool {
	var sendOrderResult bool
	_ = workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:              "send-order",
		ScheduleToCloseTimeout: time.Minute,
	}), "SendOrder", workflow.GetInfo(ctx).WorkflowExecution.ID).Get(ctx, &sendOrderResult)

	return sendOrderResult
}

func sendNotify(ctx workflow.Context) {
	_ = workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:              "notification",
		ScheduleToCloseTimeout: time.Minute,
	}), "SendNotification").Get(ctx, nil)
}

func sendCancel(ctx workflow.Context) {
	_ = workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:              "cancel-order",
		ScheduleToCloseTimeout: time.Minute,
	}), "CancelOrder").Get(ctx, nil)
}
