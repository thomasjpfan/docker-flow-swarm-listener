package service

import (
	"github.com/docker/docker/api/types/swarm"
)

func getNode(
	hostname string, nodeID string,
	role swarm.NodeRole, labels map[string]string) swarm.Node {

	annotations := swarm.Annotations{
		Labels: labels,
	}
	nodeSpec := swarm.NodeSpec{
		Annotations: annotations,
		Role:        role,
	}
	nodeDescription := swarm.NodeDescription{
		Hostname: hostname,
	}
	return swarm.Node{
		ID:          nodeID,
		Description: nodeDescription,
		Spec:        nodeSpec,
	}
}
