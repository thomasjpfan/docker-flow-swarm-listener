package main

import (
	"log"
	"os"

	"./service"
)

func main() {
	l := log.New(os.Stdout, "", log.LstdFlags)

	l.Printf("Starting Docker Flow: Swarm Listener")
	args := getArgs()
	swarmListener, err := service.NewSwarmListenerFromEnv(args.Retry, args.RetryInterval, l)
	if err != nil {
		l.Printf("Failed to initialize Docker Flow: Swarm Listener")
		l.Printf("ERROR: %v", err)
		return
	}

	swarmListener.Run()
	l.Printf("Sending notifications for running services")
	swarmListener.NotifyServices()

	serve := NewServe(swarmListener, l)
	go serve.Run()
}
