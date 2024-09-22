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

	err = firstExampleWorker.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Worker", err)
	}
}

func WorkflowDefinition(ctx workflow.Context) error {
	log.Print("Workflow started")

	err := workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:              "create-order",
		ScheduleToCloseTimeout: time.Second * 10,
	}), "CreateOrder").Get(ctx, nil)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:              "send-order",
		ScheduleToCloseTimeout: time.Second * 10,
	}), "SendOrder").Get(ctx, nil)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Second * 10,
	}), SendNotification).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func SendNotification(_ context.Context) error {
	log.Println("Send notification activity started")
	time.Sleep(time.Second)
	log.Println("Send notification activity finished")

	return nil
}
