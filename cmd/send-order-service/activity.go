package main

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"
)

type SendOrder struct{}

func NewSendOrder() *SendOrder {
	return &SendOrder{}
}

func (s *SendOrder) Handle(_ context.Context, orderID string) (bool, error) {
	log.Printf("Send order activity started, ID: %s\n", orderID)
	for i := 0; i < 5; i++ {
		log.Println("Send order activity is running")
		time.Sleep(time.Second)
	}
	log.Println("Send order activity finished")

	intID, _ := strconv.Atoi(strings.TrimLeft(orderID, "OrderWorkflow_"))
	if intID%2 == 0 {
		return true, nil
	}

	return false, nil
}
