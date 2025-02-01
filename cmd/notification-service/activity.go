package main

import (
	"context"
	"log"
	"time"
)

type Notification struct{}

func NewNotification() *Notification {
	return &Notification{}
}

func (c *Notification) Handle(_ context.Context) error {
	log.Println("Notification activity started")
	for i := 0; i < 5; i++ {
		log.Println("Notification activity is running")
		time.Sleep(time.Second)
	}

	log.Println("Notification activity finished")
	return nil
}
