package main

import (
	"context"
	"log"
	"time"

	"go.temporal.io/sdk/activity"
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

	firstExampleWorker := worker.New(temporalClient, "my-task-queue", worker.Options{
		DisableRegistrationAliasing: true,
	})
	firstExampleWorker.RegisterWorkflowWithOptions(WorkflowDefinition, workflow.RegisterOptions{
		Name: "OrderWorkflow",
	})
	firstExampleWorker.RegisterActivityWithOptions(SendNotification, activity.RegisterOptions{
		Name: "SendNotification",
	})
	firstExampleWorker.RegisterActivityWithOptions(CancelOrder, activity.RegisterOptions{
		Name: "CancelOrder",
	})

	err = firstExampleWorker.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Worker", err)
	}
}

func WorkflowDefinition(ctx workflow.Context) error {
	log.Print("Workflow started")

	createOrder(ctx)
	isSuccess := sendOrder(ctx)
	if isSuccess {
		sendNotify(ctx)
	} else {
		sendCancel(ctx)
	}

	log.Print("Workflow finished")

	return nil
}

func createOrder(ctx workflow.Context) {
	_ = workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:              "create-order",
		ScheduleToCloseTimeout: time.Second * 10,
	}), "CreateOrder").Get(ctx, nil)
}

func sendOrder(ctx workflow.Context) bool {
	var sendOrderResult bool
	_ = workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:              "send-order",
		ScheduleToCloseTimeout: time.Second * 10,
	}), "SendOrder", workflow.GetInfo(ctx).WorkflowExecution.ID).Get(ctx, &sendOrderResult)

	return sendOrderResult
}

func sendNotify(ctx workflow.Context) {
	_ = workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Second * 10,
	}), SendNotification).Get(ctx, nil)
}

func sendCancel(ctx workflow.Context) {
	_ = workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Second * 10,
	}), CancelOrder).Get(ctx, nil)
}

func SendNotification(_ context.Context) (bool, error) {
	log.Println("Send notification activity started")
	time.Sleep(time.Second)
	log.Println("Send notification activity finished")

	return true, nil
}

func CancelOrder(_ context.Context) error {
	log.Println("Cancel order activity started")
	time.Sleep(time.Second)
	log.Println("Cancel order activity finished")

	return nil
}
