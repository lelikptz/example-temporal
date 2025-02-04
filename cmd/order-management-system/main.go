package main

import (
	"log"

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
