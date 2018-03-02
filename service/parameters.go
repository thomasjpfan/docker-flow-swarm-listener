package service

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/swarm"
)

// GetNodeParameters convert `swarm.Node` metdata into `url.Values``
func GetNodeParameters(node swarm.Node) url.Values {
	params := url.Values{}

	params.Add("nodeID", node.ID)
	params.Add("hostname", node.Description.Hostname)
	params.Add("manager",
		strconv.FormatBool(node.Spec.Role == swarm.NodeRoleManager))

	for k, v := range node.Spec.Annotations.Labels {
		if !strings.HasPrefix(k, "com.df.") {
			continue
		}
		key := strings.TrimPrefix(k, "com.df.")
		params.Add(key, v)
	}

	return params
}
