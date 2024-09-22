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
	firstExampleWorker := worker.New(temporalClient, "send-order", worker.Options{})
	firstExampleWorker.RegisterActivityWithOptions(SendOrder, activity.RegisterOptions{
		Name: "SendOrder",
	})

	err := firstExampleWorker.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Worker", err)
	}
}

func SendOrder(_ context.Context) error {
	log.Println("Send order activity started")
	time.Sleep(time.Second)
	log.Println("Send order activity finished")

	return nil
}
