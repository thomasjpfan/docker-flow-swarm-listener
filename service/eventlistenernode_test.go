package service

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EventListenerNodeTestSuite struct {
	suite.Suite
}

func TestEventListenerNodeUnitTestSuite(t *testing.T) {
	s := new(EventListenerNodeTestSuite)
	logPrintfOrig := logPrintf
	defer func() {
		logPrintf = logPrintfOrig
		os.Unsetenv("DF_NOTIFY_LABEL")
	}()
	logPrintf = func(format string, v ...interface{}) {}
	os.Setenv("DF_NOTIFY_LABEL", "com.df.notify")

	suite.Run(t, s)
}

func (s *EventListenerNodeTestSuite) Test_ListenForEvents_IncorrectSocket() {
	// Bad eventlistener
	// Errors out
}

func (s *EventListenerNodeTestSuite) Test_ListenForEvents_NodeCreate() {
	// Create one node
	// Listen for events
	// Create second node
	// Get node create event
	// check nodeid points to new node is correct node
	// Destroy all nodes
}

func (s *EventListenerNodeTestSuite) Test_ListenForEvents_NodeRemove() {
	// Create two node
	// Listen for events
	// remove second node
	// Get node remove event
	// check nodeid points to removed node is correct node
	// Destroy all nodes
}

func (s *EventListenerNodeTestSuite) Test_ListenForEvents_NodeUpdate() {
	// Create one node
	// List for events
	// Update node label
	// Get node update event
	// check nodeid points to updated node
	// destroy all nodes
}

func createTestBridgeNetwork(name string) {

}

func removeTestBridgeNetwork(name string) {

}

func createNode(name string, network string, exportPort bool) {

}

func getWorkerToken(nodeName string) string {
	return ""
}

func joinSwarm(nodeName, token string) {

}

func addLabel(nodeName, label string) {

}
