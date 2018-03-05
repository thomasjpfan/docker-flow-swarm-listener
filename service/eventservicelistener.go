package service

import (
	"context"
	"log"
	"time"

	"../metrics"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// EventService is an event container the serviceType and the
// service ID
type EventService struct {
	Type EventType
	ID   string
}

// EventServiceListening listens for service events
type EventServiceListening interface {
	ListenForServiceEvents(chan<- EventService)
}

// EventServiceListener listens for docker service events
type EventServiceListener struct {
	dockerClient *client.Client
	log          *log.Logger
}

// NewEventServiceListener creates a `EventServiceListener`
func NewEventServiceListener(c *client.Client, logger *log.Logger) *EventServiceListener {
	return &EventServiceListener{dockerClient: c, log: logger}
}

// ListenForServiceEvents listens for events and places them on channels
func (s EventServiceListener) ListenForServiceEvents(eventChan chan<- EventService) {
	go func() {
		filter := filters.NewArgs()
		filter.Add("type", "service")
		msgStream, msgErrs := s.dockerClient.Events(
			context.Background(), types.EventsOptions{Filters: filter})

		for {
			select {
			case msg := <-msgStream:
				eventType := EventTypeCreate
				if msg.Action == "remove" {
					eventType = EventTypeRemove
				}
				eventChan <- EventService{
					Type: eventType,
					ID:   msg.Actor.ID,
				}
			case err := <-msgErrs:
				s.log.Printf("%v, Restarting docker event stream", err)
				metrics.RecordError("ListenForServiceEvents")
				time.Sleep(time.Second)
				// Reopen event stream
				msgStream, msgErrs = s.dockerClient.Events(
					context.Background(), types.EventsOptions{Filters: filter})
			}
		}
	}()
}
