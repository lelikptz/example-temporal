package main

import (
	"log"

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

	log.Fatalln(
		NewWorker(temporalClient, NewSendOrder()).Run(),
	)
}
