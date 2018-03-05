package service

import (
	"context"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"

	"../metrics"
)

// EventNode is an events containing the eventtype and the node ID
type EventNode struct {
	Type EventType
	ID   string
}

// EventNodeListening listens to node events
type EventNodeListening interface {
	ListenForEventNodes(<-chan EventNode)
}

// EventNodeListener listens for docker node events
type EventNodeListener struct {
	dockerClient *client.Client
	log          *log.Logger
}

// NewEventNodeListener creates a `EventNodeListener``
func NewEventNodeListener(c *client.Client, logger *log.Logger) *EventNodeListener {
	return &EventNodeListener{dockerClient: c, log: logger}
}

// ListenForEventNodes listens for events and places them on channels
func (s EventNodeListener) ListenForEventNodes(
	eventChan chan<- EventNode) {

	go func() {
		filter := filters.NewArgs()
		filter.Add("type", "node")
		msgStream, msgErrs := s.dockerClient.Events(
			context.Background(), types.EventsOptions{Filters: filter})

		for {
			select {
			case msg := <-msgStream:
				if !s.validEventNode(msg) {
					continue
				}
				eventType := EventTypeCreate
				if msg.Action == "remove" {
					eventType = EventTypeRemove
				}
				eventChan <- EventNode{
					Type: eventType,
					ID:   msg.Actor.ID,
				}
			case err := <-msgErrs:
				s.log.Printf("%v, Restarting docker event stream", err)
				metrics.RecordError("ListenForEventNodes")
				time.Sleep(time.Second)
				// Reopen event stream
				msgStream, msgErrs = s.dockerClient.Events(
					context.Background(), types.EventsOptions{Filters: filter})
			}
		}
	}()

}

// validEventNode returns false when event is valid (should be passed through)
// this will still allow through 4-5 events from changing a worker node
// to a manager node or vise versa.
func (s EventNodeListener) validEventNode(msg events.Message) bool {
	if msg.Action == "remove" {
		return true
	}
	if name, ok := msg.Actor.Attributes["name"]; !ok || len(name) == 0 {
		return false
	}
	return true
}
