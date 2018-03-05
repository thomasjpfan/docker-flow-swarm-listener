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

// ServicEventType are types of service events
type ServicEventType string

const (
	// ServicEventCreate is for create or update event
	ServicEventCreate ServicEventType = "create"
	// ServicEventRemove is for remove events
	ServicEventRemove ServicEventType = "remove"
)

// ServicEvent is an event container the serviceType and the
// service ID
type ServicEvent struct {
	Type ServicEventType
	ID   string
}

// ServicEventListening listens for service events
type ServicEventListening interface {
	ListenForServiceEvents(chan<- ServicEvent)
}

// ServicEventListener listens for docker service events
type ServicEventListener struct {
	dockerClient *client.Client
	log          *log.Logger
}

// NewServicEventListener creates a `ServicEventListener`
func NewServicEventListener(c *client.Client, logger *log.Logger) *ServicEventListener {
	return &ServicEventListener{dockerClient: c, log: logger}
}

// ListenForServiceEvents listens for events and places them on channels
func (s ServicEventListener) ListenForServiceEvents(eventChan chan<- ServicEvent) {
	go func() {
		filter := filters.NewArgs()
		filter.Add("type", "service")
		msgStream, msgErrs := s.dockerClient.Events(
			context.Background(), types.EventsOptions{Filters: filter})

		for {
			select {
			case msg := <-msgStream:
				eventType := ServicEventCreate
				if msg.Action == "remove" {
					eventType = ServicEventRemove
				}
				eventChan <- ServicEvent{
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
