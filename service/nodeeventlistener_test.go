package service

import (
	"fmt"
	"log"
	"os"
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
	s.NetworkName = "dockerflowswarmlistener_dfsl_network"
	s.Node0 = "node0"

	createNode(s.Node0, s.NetworkName)
	initSwarm(s.Node0)

	s.Node0JoinToken = getWorkerToken(s.Node0)

	host := fmt.Sprintf("tcp://%s:2375", s.Node0)
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	client, err := client.NewClient(host, dockerApiVersion, nil, defaultHeaders)
	s.Require().NoError(err)
	s.DockerClient = client

}

func (s *EventListenerNodeTestSuite) TearDownSuite() {
	destroyNode(s.Node0)
}

func (s *EventListenerNodeTestSuite) Test_ListenForNodeEvents_NodeCreate() {

	enl := NewNodeEventListener(s.DockerClient, s.Logger)

	// Listen for events
	eventChan := make(chan NodeEvent)
	enl.ListenForNodeEvents(eventChan)

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

	node1ID, err := getNodeID("node1", "node0")
	s.Require().NoError(err)

	s.Equal(node1ID, event.ID)
	s.Equal(NodeEventCreate, event.Type)
}

func (s *EventListenerNodeTestSuite) Test_ListenForNodeEvents_NodeRemove() {

	enl := NewNodeEventListener(s.DockerClient, s.Logger)

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
	eventChan := make(chan NodeEvent)
	enl.ListenForNodeEvents(eventChan)

	//Remove node1
	removeNodeFromSwarm("node1", "node0")

	// Wait for events
	event, err := s.waitForNodeEvent(eventChan)
	s.Require().NoError(err)

	s.Equal(node1ID, event.ID)
	s.Equal(NodeEventRemove, event.Type)
}

func (s *EventListenerNodeTestSuite) Test_ListenForNodeEvents_NodeUpdateLabel() {
	// Create one node
	enl := NewNodeEventListener(s.DockerClient, s.Logger)

	// Listen for events
	eventChan := make(chan NodeEvent)
	enl.ListenForNodeEvents(eventChan)

	// addLabelToNode
	addLabelToNode(s.Node0, "cats=flay", s.Node0)

	// Wait for events
	event, err := s.waitForNodeEvent(eventChan)
	s.Require().NoError(err)

	s.Equal(s.Node0, event.ID)
	s.Equal(NodeEventCreate, event.Type)

	// removeLabelFromNode
	removeLabelFromNode(s.Node0, "cats", s.Node0)

	// Wait for events
	event, err = s.waitForNodeEvent(eventChan)
	s.Require().NoError(err)

	s.Equal(s.Node0, event.ID)
	s.Equal(NodeEventCreate, event.Type)

}

func (s *EventListenerNodeTestSuite) waitForNodeEvent(events <-chan NodeEvent) (*NodeEvent, error) {
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

func createNode(name string, network string) {
	exec.Command("docker", "container", "run", "-d", "--rm",
		"--privileged", "--network", network, "--name", name,
		"--hostname", name, "docker:17.12.1-ce-dind").Output()
}

func destroyNode(name string) {
	exec.Command("docker", "container", "stop", name).Output()
}

func getWorkerToken(nodeName string) string {
	args := []string{"swarm", "join-token", "worker", "-q"}
	token, _ := runDockerCommandOnNode(args, nodeName)
	return strings.TrimRight(string(token), "\n")
}
func initSwarm(nodeName string) {
	args := []string{"swarm", "init"}
	runDockerCommandOnNode(args, nodeName)
}

func joinSwarm(nodeName, rootNodeName, token string) {
	rootHost := fmt.Sprintf("%s:2377", rootNodeName)
	args := []string{"swarm", "join", "--token", token, rootHost}
	runDockerCommandOnNode(args, nodeName)
}

func getNodeID(nodeName, rootNodeName string) (string, error) {
	args := []string{"node", "inspect", nodeName, "-f", "{{ .ID }}"}
	ID, err := runDockerCommandOnNode(args, rootNodeName)
	return strings.TrimRight(string(ID), "\n"), err
}

func removeNodeFromSwarm(nodeName, rootNodeName string) {
	args := []string{"node", "rm", "--force", nodeName}
	runDockerCommandOnNode(args, rootNodeName)
}

func addLabelToNode(nodeName, label, rootNodeName string) {
	args := []string{"node", "update", "--label-add", label, nodeName}
	runDockerCommandOnNode(args, nodeName)
}

func removeLabelFromNode(nodeName, label, rootNodeName string) {
	args := []string{"node", "update", "--label-rm", label, nodeName}
	runDockerCommandOnNode(args, nodeName)
}

func runDockerCommandOnNode(args []string, nodeName string) (string, error) {
	host := fmt.Sprintf("tcp://%s:2375", nodeName)
	dockerCmd := []string{"-H", host}
	fullCmd := append(dockerCmd, args...)
	output, err := exec.Command("docker", fullCmd...).Output()
	return string(output), err
}
