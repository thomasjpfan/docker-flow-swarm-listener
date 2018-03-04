package service

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/suite"
)

type ServicEventListenerTestSuite struct {
	suite.Suite
	ServiceName  string
	DockerClient *client.Client
	Logger       *log.Logger
}

func TestServicEventListenerTestSuite(t *testing.T) {
	suite.Run(t, new(ServicEventListenerTestSuite))
}

func (s *ServicEventListenerTestSuite) SetupSuite() {
	s.ServiceName = "my-service"
	client, err := NewDockerClientFromEnv()
	s.Require().NoError(err)
	s.DockerClient = client.Client
	s.Logger = log.New(os.Stdout, "", 0)
}

func (s *ServicEventListenerTestSuite) Test_ListenForServiceEvents_CreateService() {
	snl := NewServicEventListener(s.DockerClient, s.Logger)

	// Listen for events
	eventChan := make(chan ServicEvent)
	snl.ListenForServiceEvents(eventChan)

	createTestService("util-1", []string{}, "", "")
	defer func() {
		removeTestService("util-1")
	}()

	time.Sleep(time.Second)
	utilID := getServiceID("util-1")

	event, err := s.waitForServiceEvent(eventChan)
	s.Require().NoError(err)

	s.Equal(ServicEventCreate, event.Type)
	s.Equal(utilID, event.ID)
}

func (s *ServicEventListenerTestSuite) Test_ListenForServiceEvents_UpdateService() {
	snl := NewServicEventListener(s.DockerClient, s.Logger)

	createTestService("util-1", []string{}, "", "")
	defer func() {
		removeTestService("util-1")
	}()

	time.Sleep(time.Second)
	utilID := getServiceID("util-1")

	// Listen for events
	eventChan := make(chan ServicEvent)
	snl.ListenForServiceEvents(eventChan)

	// Update label
	addLabelToService("util-1", "hello=world")

	event, err := s.waitForServiceEvent(eventChan)
	s.Require().NoError(err)

	s.Equal(ServicEventCreate, event.Type)
	s.Equal(utilID, event.ID)

	// Remove label
	removeLabelFromService("util-1", "hello")

	event, err = s.waitForServiceEvent(eventChan)
	s.Require().NoError(err)

	s.Equal(ServicEventCreate, event.Type)
	s.Equal(utilID, event.ID)
}

func (s *ServicEventListenerTestSuite) Test_ListenForServiceEvents_RemoveService() {
	snl := NewServicEventListener(s.DockerClient, s.Logger)

	createTestService("util-1", []string{}, "", "")
	defer func() {
		removeTestService("util-1")
	}()

	time.Sleep(time.Second)
	utilID := getServiceID("util-1")

	// Listen for events
	eventChan := make(chan ServicEvent)
	snl.ListenForServiceEvents(eventChan)

	// Remove service
	removeTestService("util-1")

	event, err := s.waitForServiceEvent(eventChan)
	s.Require().NoError(err)

	s.Equal(ServicEventRemove, event.Type)
	s.Equal(utilID, event.ID)
}

func (s *ServicEventListenerTestSuite) waitForServiceEvent(events <-chan ServicEvent) (*ServicEvent, error) {
	timeOut := time.NewTimer(time.Second * 5).C
	for {
		select {
		case event := <-events:
			return &event, nil
		case <-timeOut:
			return nil, fmt.Errorf("Timeout")
		}
	}
}

func addLabelToService(name, label string) {
	args := []string{"service", "update", "--label-add", label, name}
	runDockerCommandOnSocket(args)
}

func removeLabelFromService(name, label string) {
	args := []string{"service", "update", "--label-rm", label, name}
	runDockerCommandOnSocket(args)
}

func runDockerCommandOnSocket(args []string) (string, error) {
	output, err := exec.Command("docker", args...).Output()
	return string(output), err
}
