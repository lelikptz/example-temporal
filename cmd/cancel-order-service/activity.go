package main

import (
	"context"
	"log"
	"time"
)

type CancelOrder struct{}

func NewCancelOrder() *CancelOrder {
	return &CancelOrder{}
}

func (c *CancelOrder) Handle(_ context.Context) error {
	log.Println("Cancel order activity started")
	for i := 0; i < 5; i++ {
		log.Println("Cancel order activity is running")
		time.Sleep(time.Second)
	}

	log.Println("Cancel order activity finished")
	return nil
}
