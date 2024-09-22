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

func runTask(temporalClient client.Client, wg *sync.WaitGroup) {
	defer wg.Done()
	counter := 1
	for {
		workflowRun, err := temporalClient.ExecuteWorkflow(
			context.Background(),
			client.StartWorkflowOptions{
				ID:        fmt.Sprintf("OrderWorkflow_%d", counter),
				TaskQueue: "my-task-queue",
			},
			"OrderWorkflow",
		)
		if err != nil {
			log.Printf("Unable to execute workflow %s", err)
			continue
		}

		err = workflowRun.Get(context.Background(), nil)
		if err != nil {
			log.Printf("Unable to get workflow %s", err)
			continue
		}

		log.Printf("Result %s", "success")

		counter++
		time.Sleep(5 * time.Second)
	}
}
