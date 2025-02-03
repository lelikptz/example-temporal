package main

import (
	"context"
	"log"
	"time"
)

type CreateOrder struct{}

func NewCreateOrder() *CreateOrder {
	return &CreateOrder{}
}

func (c *CreateOrder) Handle(_ context.Context) error {
	log.Println("Create order activity started")
	for i := 0; i < 15; i++ {
		log.Println("Create order activity is running")
		time.Sleep(time.Second)
	}
	log.Println("Create order activity finished")

	return nil
}
