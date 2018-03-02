package service

import (
	"os"
	"strings"

	"github.com/docker/docker/api/types/swarm"
)

// EventNodeNotifing notifies on a node event
type EventNodeNotifing interface {
	NotifyCreateNodes(nodes []swarm.Node, retry, retryInterval int) error
	NotifyCreateNode(node swarm.Node, retry, retryInterval int) error
	NotifyUpdateNode(node swarm.Node, retry, retryInterval int) error
	NotifyRemoveNode(node swarm.Node, retry, retryInterval int) error
	HasListeners() bool
}

// EventNodeNotifier sends out node event notifications
type EventNodeNotifier struct {
	CreateAddrs []string
	UpdateAddrs []string
	RemoveAddrs []string
}

func newEventNodeNotifier(
	createAddrs, updateAddrs, removeAddrs []string) *EventNodeNotifier {
	return &EventNodeNotifier{
		CreateAddrs: createAddrs,
		UpdateAddrs: updateAddrs,
		RemoveAddrs: removeAddrs,
	}
}

// NewEventNodeNotifierFromEnv creats a `EventNodeNotifier` from env variables
func NewEventNodeNotifierFromEnv() *EventNodeNotifier {
	createNodeENV := os.Getenv("DF_NOTIFY_CREATE_NODE_URL")
	updateNodeENV := os.Getenv("DF_NOTIFY_UPDATE_NODE_URL")
	removeNodeENV := os.Getenv("DF_NOTIFY_REMOVE_NODE_URL")

	var createAddrs, updateAddrs, removeAddrs []string

	if len(createNodeENV) > 0 {
		createAddrs = strings.Split(createNodeENV, ",")
	}
	if len(updateNodeENV) > 0 {
		updateAddrs = strings.Split(updateNodeENV, ",")
	}
	if len(removeNodeENV) > 0 {
		removeAddrs = strings.Split(removeNodeENV, ",")
	}

	return newEventNodeNotifier(
		createAddrs,
		updateAddrs,
		removeAddrs,
	)

}

// NotifyCreateNodes notifies addresses with create notification
func (n EventNodeNotifier) NotifyCreateNodes(nodes []swarm.Node, retry, retryInterval int) error {
	return nil
}

// NotifyCreateNode notifies addresses with create notification
func (n EventNodeNotifier) NotifyCreateNode(node swarm.Node, retry, retryInterval int) error {
	return nil
}

// NotifyUpdateNode notifies addresses with update notification
func (n EventNodeNotifier) NotifyUpdateNode(node swarm.Node, retry, retryInterval int) error {
	return nil
}

// NotifyRemoveNode notifies addresses with remove notification
func (n EventNodeNotifier) NotifyRemoveNode(node swarm.Node, retry, retryInterval int) error {
	return nil
}

// HasListeners returns true when there are addresses to send
// notifications to
func (n EventNodeNotifier) HasListeners() bool {
	return (len(n.CreateAddrs) + len(n.UpdateAddrs) + len(n.RemoveAddrs)) > 0
}
