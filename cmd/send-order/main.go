package main

import (
	"context"
	"log"
	"strconv"
	"strings"
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
	firstExampleWorker.RegisterActivityWithOptions(NewSendOrder().Run, activity.RegisterOptions{
		Name: "SendOrder",
	})

	err := firstExampleWorker.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Worker", err)
	}
}

type SendOrder struct{}

func NewSendOrder() *SendOrder {
	return &SendOrder{}
}

func (s *SendOrder) Run(_ context.Context, orderID string) (bool, error) {
	log.Printf("Send order activity started, ID: %s\n", orderID)
	time.Sleep(time.Second)
	log.Println("Send order activity finished")

	intID, _ := strconv.Atoi(strings.TrimLeft(orderID, "OrderWorkflow_"))
	if intID%2 == 0 {
		return true, nil
	}

	return false, nil
}
