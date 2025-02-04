package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"go.temporal.io/api/common/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	ctx := context.Background()
	ID := fmt.Sprintf("%s", time.Now().Format(time.DateTime))
	_, err := temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:        ID,
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

	time.Sleep(time.Duration(rand.Intn(60)) * time.Second)

	conn, err := grpc.NewClient(client.DefaultHostPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer conn.Close()

	grpcClient := workflowservice.NewWorkflowServiceClient(conn)
	req := &workflowservice.DescribeWorkflowExecutionRequest{
		Namespace: "default",
		Execution: &common.WorkflowExecution{
			WorkflowId: ID,
		},
	}
	resp, _ := grpcClient.DescribeWorkflowExecution(ctx, req)

	_ = temporalClient.SignalWorkflow(
		ctx,
		resp.WorkflowExecutionInfo.GetExecution().GetWorkflowId(),
		resp.WorkflowExecutionInfo.GetExecution().GetRunId(),
		"cancel-order-signal",
		nil,
	)
}
