package service

import (
	"strings"

	"github.com/docker/docker/api/types/swarm"
)

// MinifyNode minifies `swarm.Node`
// only labels prefixed with `com.df.` will be used
func MinifyNode(n swarm.Node) NodeMini {
	engineLabels := map[string]string{}
	for k, v := range n.Description.Engine.Labels {
		if strings.HasPrefix(k, "com.df.") {
			engineLabels[k] = v
		}
	}
	nodeLabels := map[string]string{}
	for k, v := range n.Spec.Labels {
		if strings.HasPrefix(k, "com.df.") {
			nodeLabels[k] = v
		}
	}

	return NodeMini{
		ID:           n.ID,
		Hostname:     n.Description.Hostname,
		VersionIndex: n.Meta.Version.Index,
		State:        n.Status.State,
		Addr:         n.Status.Addr,
		NodeLabels:   nodeLabels,
		EngineLabels: engineLabels,
		Role:         n.Spec.Role,
		Availability: n.Spec.Availability,
	}
}

// MinifySwarmService minifies `SwarmService`
// only labels prefixed with `com.df.` will be used
// `ignoreKey` wll be ignored from labels
func MinifySwarmService(ss SwarmService, ignoreKey string) SwarmServiceMini {
	filterLabels := map[string]string{}
	for k, v := range ss.Spec.Labels {
		if k != ignoreKey && strings.HasPrefix(k, "com.df.") {
			filterLabels[k] = v
		}
	}
	return SwarmServiceMini{
		ID:       ss.ID,
		Name:     ss.Spec.Name,
		Labels:   filterLabels,
		Mode:     ss.Spec.Mode,
		NodeInfo: ss.NodeInfo,
	}
}
