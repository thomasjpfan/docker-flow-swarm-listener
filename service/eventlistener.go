package service

import (
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// eventAdmitter admits docker events
type eventAdmitter interface {
	DockerEvents(eventType string) (<-chan events.Message, <-chan error)
}

// dockerClient implements eventAdmitter interface
type dockerClient struct {
	*client.Client
}

// DockerEvents uses docker client to listen for docker events
func (c dockerClient) DockerEvents(eventType string) (<-chan events.Message, <-chan error) {
	filter := filters.NewArgs()
	filter.Add("type", eventType)
	return c.Events(context.Background(), types.EventsOptions{Filters: filter})
}

// Event contains information about docker events
type Event struct {
	Action string
	ID     string
}

// EventListener listens for docker service events
type EventListener struct {
	eventAdmitter
	eventType string
}

// NewEventListener returns `EventListener` for an eventAdmitter
func NewEventListener(eventAdmit eventAdmitter, eventType string) *EventListener {
	return &EventListener{eventAdmit, eventType}
}

// NewEventListenerForDocker creates a `EventListener` with a docker host
func NewEventListenerForDocker(host, eventType string) *EventListener {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	dc, err := client.NewClient(host, dockerApiVersion, nil, defaultHeaders)
	if err != nil {
		logPrintf(err.Error())
	}
	return NewEventListener(dockerClient{dc}, eventType)
}

// NewEventListenerFromEnv returns a new instance of the `EventListener` structure using environment variable `DF_DOCKER_HOST` for the host
func NewEventListenerFromEnv(eventType string) *EventListener {
	host := "unix:///var/run/docker.sock"
	if len(os.Getenv("DF_DOCKER_HOST")) > 0 {
		host = os.Getenv("DF_DOCKER_HOST")
	}
	return NewEventListenerForDocker(host, eventType)
}

// ListenForEvents returns a stream of Events
func (s *EventListener) ListenForEvents() (<-chan Event, <-chan error) {

	events := make(chan Event)
	errs := make(chan error, 1)
	started := make(chan struct{})

	go func() {
		defer close(errs)
		eventStream, eventErrors := s.DockerEvents(s.eventType)

		close(started)
		for {
			select {
			case msg := <-eventStream:
				events <- Event{
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
