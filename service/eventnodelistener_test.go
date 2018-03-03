package service

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/suite"
)

type EventListenerNodeTestSuite struct {
	suite.Suite
	DockerClient   *client.Client
	NetworkName    string
	Node0          string
	Node0JoinToken string
}

func TestEventListenerNodeUnitTestSuite(t *testing.T) {
	s := new(EventListenerNodeTestSuite)
	logPrintfOrig := logPrintf
	defer func() {
		logPrintf = logPrintfOrig
	}()
	logPrintf = func(format string, v ...interface{}) {}

	suite.Run(t, s)
}

func (s *EventListenerNodeTestSuite) SetupSuite() {
	s.NetworkName = "dockerflowswarmlistener_dfsl_network"
	s.Node0 = "node0"

	createNode(s.Node0, s.NetworkName)
	initSwarm(s.Node0)

	s.Node0JoinToken = getWorkerToken(s.Node0)

	host := fmt.Sprintf("tcp://%s:2375", s.Node0)
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	client, _ := client.NewClient(host, dockerApiVersion, nil, defaultHeaders)
	s.DockerClient = client

}

func (s *EventListenerNodeTestSuite) TearDownSuite() {
	destroyNode(s.Node0)
}

func (s *EventListenerNodeTestSuite) Test_ListenForEvents_NodeCreate() {

	enl := NewEventNodeListener(s.DockerClient)

	// Listen for events
	eventChan := make(chan Event)
	enl.ListenForEvents(eventChan)

	// Create node1
	createNode("node1", s.NetworkName)
	defer func() {
		destroyNode("node1")
	}()

	time.Sleep(time.Second)
	joinSwarm("node1", s.Node0, s.Node0JoinToken)

	// Wait for events
	event, err := s.waitForNodeEvent(eventChan)
	s.Require().NoError(err)

	node1ID := getNodeID("node1", "node0")
	s.Equal(node1ID, event.ID)
	s.Equal("create", event.Action)
}

func (s *EventListenerNodeTestSuite) Test_ListenForEvents_NodeRemove() {

	enl := NewEventNodeListener(s.DockerClient)

	// Create node1 and joing swarm
	createNode("node1", s.NetworkName)
	defer func() {
		destroyNode("node1")
	}()
	time.Sleep(time.Second)
	joinSwarm("node1", s.Node0, s.Node0JoinToken)

	time.Sleep(time.Second)
	node1ID := getNodeID("node1", "node0")

	// Listen for events
	eventChan := make(chan Event)
	enl.ListenForEvents(eventChan)

	//Remove node1
	time.Sleep(4 * time.Second)
	destroyNode("node1")

	// Wait for events
	event, err := s.waitForNodeEvent(eventChan)
	s.Require().NoError(err)

	s.Equal(node1ID, event.ID)
	s.Equal("remove", event.Action)
}

func (s *EventListenerNodeTestSuite) Test_ListenForEvents_NodeUpdate() {
	// Create one node
	// List for events
	// Update node label
	// Get node update event
	// check nodeid points to updated node
	// destroy all nodes
}

func (s *EventListenerNodeTestSuite) waitForNodeEvent(events <-chan Event) (*Event, error) {
	timeOut := time.NewTimer(time.Second * 10).C
	for {
		select {
		case event := <-events:
			return &event, nil
		case <-timeOut:
			return nil, fmt.Errorf("Timeout")
		}
	}
}
func createNode(name string, network string) {
	exec.Command("docker", "container", "run", "-d", "--rm",
		"--privileged", "--network", network, "--name", name,
		"--hostname", name, "docker:17.12.1-ce-dind").Output()
}

func destroyNode(name string) {
	exec.Command("docker", "container", "stop", name).Output()
}

func getWorkerToken(nodeName string) string {
	host := fmt.Sprintf("tcp://%s:2375", nodeName)
	token, _ := exec.Command("docker", "-H", host, "swarm", "join-token", "worker", "-q").Output()
	return strings.TrimRight(string(token), "\n")
}
func initSwarm(nodeName string) {
	host := fmt.Sprintf("tcp://%s:2375", nodeName)
	exec.Command("docker", "-H", host, "swarm", "init").Output()
}

func joinSwarm(nodeName, rootNodeName, token string) {
	host := fmt.Sprintf("tcp://%s:2375", nodeName)
	rootHost := fmt.Sprintf("%s:2377", rootNodeName)
	exec.Command("docker", "-H", host, "swarm", "join", "--token", token, rootHost).Output()
}

func getNodeID(nodeName, rootNodeName string) string {
	host := fmt.Sprintf("tcp://%s:2375", rootNodeName)
	ID, _ := exec.Command("docker", "-H", host, "node", "inspect", nodeName, "-f", `{{ .ID }}`).Output()
	return strings.TrimRight(string(ID), "\n")
}

func addLabel(nodeName, label string) {

}
