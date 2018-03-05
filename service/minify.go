package service

import (
	"github.com/docker/docker/api/types/swarm"
)

// MinifyNode minifies `swarm.Node`
func MinifyNode(n swarm.Node) NodeMini {
	return NodeMini{}
}

// MinifySwarmService minifies `SwarmService`
func MinifySwarmService(ss SwarmService) SwarmServiceMini {
	return SwarmServiceMini{}
}
