package main

import (
	"context"
	"log"
	"time"
)

type StatusPoll struct{}

func NewStatusPoll() *StatusPoll {
	return &StatusPoll{}
}

func (c *StatusPoll) Handle(_ context.Context) (string, error) {
	log.Println("Start status polling")
	for i := 0; i < 15; i++ {
		log.Println("Status polling is running")
		time.Sleep(time.Second)
	}

	log.Println("Status polling finished")

	return "Finished", nil
}
