package main

import (
	"context"
	"log"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
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

	runWorker(temporalClient)
}

func runWorker(temporalClient client.Client) {
	firstExampleWorker := worker.New(temporalClient, "create-order", worker.Options{})
	firstExampleWorker.RegisterActivityWithOptions(CreateOrder, activity.RegisterOptions{
		Name: "CreateOrder",
	})

	err := firstExampleWorker.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Worker", err)
	}
}

func CreateOrder(_ context.Context) error {
	log.Println("Create order activity started")
	time.Sleep(time.Second)
	log.Println("Create order activity finished")

	return nil
}
