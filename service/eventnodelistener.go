package service

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// EventNodeListening listens to node events
type EventNodeListening interface {
	ListenForEvents(<-chan Event, <-chan error)
}

// EventNodeListener listens for docker node events
type EventNodeListener struct {
	dockerClient *client.Client
}

// NewEventNodeListener creates a `EventNodeListener``
func NewEventNodeListener(c *client.Client) *EventNodeListener {
	return &EventNodeListener{dockerClient: c}
}

// ListenForEvents listens for events and places them on channels
func (s EventNodeListener) ListenForEvents(
	eventChan chan<- Event) {

	go func() {
		filter := filters.NewArgs()
		filter.Add("type", "node")
		msgStream, msgErrs := s.dockerClient.Events(
			context.Background(), types.EventsOptions{Filters: filter})

		for {
			select {
			case msg := <-msgStream:
				eventChan <- Event{
					Action: msg.Action,
					ID:     msg.Actor.ID,
				}
			case err := <-msgErrs:
				logPrintf("%v, Restarting docker event stream", err)
				// Reopen event stream
				msgStream, msgErrs = s.dockerClient.Events(
					context.Background(), types.EventsOptions{Filters: filter})
			}
		}
	}()

}
