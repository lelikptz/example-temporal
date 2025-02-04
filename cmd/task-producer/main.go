package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.temporal.io/sdk/client"
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

	wg := sync.WaitGroup{}
	wg.Add(1)
	go runTask(temporalClient, &wg)

	wg.Wait()
}

type Order struct {
	ReservationTime time.Time
}

func runTask(temporalClient client.Client, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := temporalClient.ExecuteWorkflow(
		context.Background(),
		client.StartWorkflowOptions{
			ID:        fmt.Sprintf("%s", time.Now().Format(time.DateTime)),
			TaskQueue: "order-queue",
		},
		"OrderWorkflow",
		Order{
			ReservationTime: time.Now().Add(40 * time.Second),
		},
	)
	if err != nil {
		log.Printf("Unable to execute workflow %s", err)
		return
	}
}
