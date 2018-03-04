package service

import (
	"os"
	"strings"

	"github.com/docker/docker/api/types/swarm"
)

// EventNodeNotifing notifies on a node event
type EventNodeNotifing interface {
	NotifyCreateNode(node swarm.Node, retry, retryInterval int) error
	NotifyRemoveNode(node swarm.Node, retry, retryInterval int) error
	HasListeners() bool
}

// EventNodeNotifier sends out node event notifications
type EventNodeNotifier struct {
	CreateAddrs []string
	RemoveAddrs []string
}

func newEventNodeNotifier(
	createAddrs, removeAddrs []string) *EventNodeNotifier {
	return &EventNodeNotifier{
		CreateAddrs: createAddrs,
		RemoveAddrs: removeAddrs,
	}
}

// NewEventNodeNotifierFromEnv creats a `EventNodeNotifier` from env variables
func NewEventNodeNotifierFromEnv() *EventNodeNotifier {
	createNodeENV := os.Getenv("DF_NOTIFY_CREATE_NODE_URL")
	removeNodeENV := os.Getenv("DF_NOTIFY_REMOVE_NODE_URL")

	var createAddrs, removeAddrs []string

	if len(createNodeENV) > 0 {
		createAddrs = strings.Split(createNodeENV, ",")
	}
	if len(removeNodeENV) > 0 {
		removeAddrs = strings.Split(removeNodeENV, ",")
	}

	return newEventNodeNotifier(
		createAddrs,
		removeAddrs,
	)

}

// NotifyCreateNode notifies addresses with create notification
func (n EventNodeNotifier) NotifyCreateNode(node swarm.Node, retry, retryInterval int) error {
	return nil
}

// NotifyRemoveNode notifies addresses with remove notification
func (n EventNodeNotifier) NotifyRemoveNode(node swarm.Node, retry, retryInterval int) error {
	return nil
}

// HasListeners returns true when there are addresses to send
// notifications to
func (n EventNodeNotifier) HasListeners() bool {
	return (len(n.CreateAddrs) + len(n.RemoveAddrs)) > 0
}
