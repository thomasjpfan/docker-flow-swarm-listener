package service

import (
	"context"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// EventOld contains information about docker events
type EventOld struct {
	Action string
	ID     string
}

// EventListening object listens for events
type EventListening interface {
	ListenForEvents() (<-chan EventOld, <-chan error)
}

// EventListener listens for docker service events
type EventListener struct {
	*client.Client
	eventType string
}

// NewEventListener creates a `EventListener` with a docker host
func NewEventListener(host, eventType string) *EventListener {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	dc, err := client.NewClient(host, dockerAPIVersion, nil, defaultHeaders)
	if err != nil {
		logPrintf(err.Error())
	}
	return &EventListener{dc, eventType}
}

// NewEventListenerFromEnv returns a new instance of the `EventListener` structure using environment variable `DF_DOCKER_HOST` for the host
func NewEventListenerFromEnv(eventType string) *EventListener {
	host := "unix:///var/run/docker.sock"
	if len(os.Getenv("DF_DOCKER_HOST")) > 0 {
		host = os.Getenv("DF_DOCKER_HOST")
	}
	return NewEventListener(host, eventType)
}

// ListenForEvents returns a stream of Events
func (s *EventListener) ListenForEvents() (<-chan EventOld, <-chan error) {

	events := make(chan EventOld)
	errs := make(chan error, 1)
	started := make(chan struct{})

	go func() {
		defer close(errs)
		filter := filters.NewArgs()
		filter.Add("type", s.eventType)
		eventStream, eventErrors := s.Events(
			context.Background(), types.EventsOptions{Filters: filter})

		close(started)
		for {
			select {
			case msg := <-eventStream:
				events <- EventOld{
					Action: msg.Action,
					ID:     msg.Actor.ID,
				}
			case err := <-eventErrors:
				logPrintf("%v", err)
				errs <- err
				return
			}
		}

	}()
	<-started

	return events, errs
}
