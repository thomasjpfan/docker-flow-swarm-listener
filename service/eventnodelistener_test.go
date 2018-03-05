package service

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/suite"
)

type EventListenerNodeTestSuite struct {
	suite.Suite
	DockerClient   *client.Client
	Logger         *log.Logger
	NetworkName    string
	Node0          string
	Node0JoinToken string
}

func TestEventListenerNodeUnitTestSuite(t *testing.T) {
	suite.Run(t, new(EventListenerNodeTestSuite))
}

func (s *EventListenerNodeTestSuite) SetupSuite() {
	s.Logger = log.New(os.Stdout, "", 0)

	// Assumes running test with docker-compose.yml
	s.NetworkName = "dockerflowswarmlistener_dfsl_network"
	s.Node0 = "node0"

	createNode(s.Node0, s.NetworkName)
	initSwarm(s.Node0)

	s.Node0JoinToken = getWorkerToken(s.Node0)

	client, err := newTestNodeDockerClient(s.Node0)
	s.Require().NoError(err)
	s.DockerClient = client

}

func (s *EventListenerNodeTestSuite) TearDownSuite() {
	destroyNode(s.Node0)
}

func (s *EventListenerNodeTestSuite) Test_ListenForEventNodes_NodeCreate() {

	enl := NewEventNodeListener(s.DockerClient, s.Logger)

	// Listen for events
	eventChan := make(chan EventNode)
	enl.ListenForEventNodes(eventChan)

	// Create node1
	createNode("node1", s.NetworkName)
	defer func() {
		destroyNode("node1")
	}()

	time.Sleep(time.Second)
	joinSwarm("node1", s.Node0, s.Node0JoinToken)

	// Wait for events
	event, err := s.waitForEventNode(eventChan)
	s.Require().NoError(err)

	node1ID, err := getNodeID("node1", "node0")
	s.Require().NoError(err)

	s.Equal(node1ID, event.ID)
	s.Equal(EventTypeCreate, event.Type)
}

func (s *EventListenerNodeTestSuite) Test_ListenForEventNodes_NodeRemove() {

	enl := NewEventNodeListener(s.DockerClient, s.Logger)

	// Create node1 and joing swarm
	createNode("node1", s.NetworkName)
	defer func() {
		destroyNode("node1")
	}()
	joinSwarm("node1", s.Node0, s.Node0JoinToken)

	time.Sleep(time.Second)
	node1ID, err := getNodeID("node1", "node0")
	s.Require().NoError(err)

	// Listen for events
	eventChan := make(chan EventNode)
	enl.ListenForEventNodes(eventChan)

	//Remove node1
	removeNodeFromSwarm("node1", "node0")

	// Wait for events
	event, err := s.waitForEventNode(eventChan)
	s.Require().NoError(err)

	s.Equal(node1ID, event.ID)
	s.Equal(EventTypeRemove, event.Type)
}

func (s *EventListenerNodeTestSuite) Test_ListenForEventNodes_NodeUpdateLabel() {
	// Create one node
	enl := NewEventNodeListener(s.DockerClient, s.Logger)

	// Listen for events
	eventChan := make(chan EventNode)
	enl.ListenForEventNodes(eventChan)

	// addLabelToNode
	addLabelToNode(s.Node0, "cats=flay", s.Node0)

	// Wait for events
	event, err := s.waitForEventNode(eventChan)
	s.Require().NoError(err)

	s.Equal(s.Node0, event.ID)
	s.Equal(EventTypeCreate, event.Type)

	// removeLabelFromNode
	removeLabelFromNode(s.Node0, "cats", s.Node0)

	// Wait for events
	event, err = s.waitForEventNode(eventChan)
	s.Require().NoError(err)

	s.Equal(s.Node0, event.ID)
	s.Equal(EventTypeCreate, event.Type)

}

func (s *EventListenerNodeTestSuite) waitForEventNode(events <-chan EventNode) (*EventNode, error) {
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
