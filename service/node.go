package service

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
)

// NodeInspector is able to inspect a swarm node
type NodeInspector interface {
	NodeInspect(ctx context.Context, nodeID string) (swarm.Node, error)
	NodeList(ctx context.Context) ([]swarm.Node, error)
}

// NodeInspect returns swarm.Node from its ID
func (c DockerClient) NodeInspect(ctx context.Context, nodeID string) (swarm.Node, error) {
	node, _, err := c.Client.NodeInspectWithRaw(ctx, nodeID)
	return node, err
}

// NodeList returns a list of all nodes
func (c DockerClient) NodeList(ctx context.Context) ([]swarm.Node, error) {
	return c.Client.NodeList(ctx, types.NodeListOptions{})
}
