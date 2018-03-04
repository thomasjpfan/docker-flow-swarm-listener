package service

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"

	"../metrics"
)

// NodeEventType are types of node events
type NodeEventType string

const (
	// NodeEventCreate is for create or update events
	NodeEventCreate NodeEventType = "create"
	// NodeEventRemove is for remove events
	NodeEventRemove NodeEventType = "remove"
)

// NodeEvent is an events containing the eventtype and the node ID
type NodeEvent struct {
	Type NodeEventType
	ID   string
}

// NodeEventListening listens to node events
type NodeEventListening interface {
	ListenForNodeEvents(<-chan NodeEvent)
}

// NodeEventListener listens for docker node events
type NodeEventListener struct {
	dockerClient *client.Client
	log          *log.Logger
}

// NewNodeEventListener creates a `NodeEventListener``
func NewNodeEventListener(c *client.Client, logger *log.Logger) *NodeEventListener {
	return &NodeEventListener{dockerClient: c, log: logger}
}

// ListenForNodeEvents listens for events and places them on channels
func (s NodeEventListener) ListenForNodeEvents(
	eventChan chan<- NodeEvent) {

	go func() {
		filter := filters.NewArgs()
		filter.Add("type", "node")
		msgStream, msgErrs := s.dockerClient.Events(
			context.Background(), types.EventsOptions{Filters: filter})

		for {
			select {
			case msg := <-msgStream:
				if !s.validNodeEvent(msg) {
					continue
				}
				eventType := NodeEventCreate
				if msg.Action == "remove" {
					eventType = NodeEventRemove
				}
				eventChan <- NodeEvent{
					Type: eventType,
					ID:   msg.Actor.ID,
				}
			case err := <-msgErrs:
				s.log.Printf("%v, Restarting docker event stream", err)
				metrics.RecordError("ListenForNodeEvents")
				// Reopen event stream
				msgStream, msgErrs = s.dockerClient.Events(
					context.Background(), types.EventsOptions{Filters: filter})
			}
		}
	}()

}

// validNodeEvent returns false when event is valid (should be passed through)
// this will still allow through 4-5 events from changing a worker node
// to a manager node or vise versa.
func (s NodeEventListener) validNodeEvent(msg events.Message) bool {
	if msg.Action == "remove" {
		return true
	}
	if name, ok := msg.Actor.Attributes["name"]; !ok || len(name) == 0 {
		return false
	}
	return true
}
