package service

import (
	"os"
	"strings"

	"github.com/docker/docker/api/types/swarm"
)

// NodeEventNotifing notifies on a node event
type NodeEventNotifing interface {
	NotifyCreateNodes(nodes []swarm.Node, retry, retryInterval int)
	NotifyCreateNode(node swarm.Node, retry, retryInterval int)
	NotifyUpdateNode(node swarm.Node, retry, retryInterval int)
	NotifyRemoveNode(node swarm.Node, retry, retryInterval int)
}

// NodeEventNotifier sends out node event notifications
type NodeEventNotifier struct {
	CreateAddrs []string
	UpdateAddrs []string
	RemoveAddrs []string
}

func newNodeEventNotifier(
	createAddrs, updateAddrs, removeAddrs []string) *NodeEventNotifier {
	return &NodeEventNotifier{
		CreateAddrs: createAddrs,
		UpdateAddrs: updateAddrs,
		RemoveAddrs: removeAddrs,
	}
}

// NewNodeEventNotifierFromEnv creats a `NodeEventNotifier` from env variables
func NewNodeEventNotifierFromEnv() *NodeEventNotifier {
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

	return newNodeEventNotifier(
		createAddrs,
		updateAddrs,
		removeAddrs,
	)

}

// NotifyCreateNodes notifies addresses with create notification
func (n NodeEventNotifier) NotifyCreateNodes(nodes []swarm.Node, retry, retryInterval int) {

}

// NotifyCreateNode notifies addresses with create notification
func (n NodeEventNotifier) NotifyCreateNode(node swarm.Node, retry, retryInterval int) {

}

// NotifyUpdateNode notifies addresses with update notification
func (n NodeEventNotifier) NotifyUpdateNode(node swarm.Node, retry, retryInterval int) {

}

// NotifyRemoveNode notifies addresses with remove notification
func (n NodeEventNotifier) NotifyRemoveNode(node swarm.Node, retry, retryInterval int) {

}
