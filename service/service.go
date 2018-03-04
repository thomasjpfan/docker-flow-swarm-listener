package service

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// ServicInspector is able to inspect services
type ServicInspector interface {
	ServicInspect(serviceID string, includeNodeIPInfo bool) (*SwarmService, error)
	ServicList(ctx context.Context, includeNodeIPInfo bool) ([]SwarmService, error)
}

// ServicClient implements `ServicInspector` for docker
type ServicClient struct {
	DockerClient   *client.Client
	FilterLabel    string
	ScrapeNetLabel string
}

// NewServicClient creats a `ServicClient`
func NewServicClient(c *client.Client, filterLabel, scrapNetLabel string) *ServicClient {
	return &ServicClient{DockerClient: c, FilterLabel: filterLabel,
		ScrapeNetLabel: scrapNetLabel}
}

// ServicInspect returns `SwarmService` from its ID
// Returns nil when service doesnt not have the `FilterLabel`
// When `includeNodeIPInfo` is true, return node info as well
func (c ServicClient) ServicInspect(serviceID string, includeNodeIPInfo bool) (*SwarmService, error) {
	service, _, err := c.DockerClient.ServiceInspectWithRaw(context.Background(), serviceID, types.ServiceInspectOptions{})
	if err != nil {
		return nil, err
	}

	// Check if service has label
	key := strings.SplitN(c.FilterLabel, "=", 2)[0]
	if _, ok := service.Spec.Labels[key]; !ok {
		return nil, nil
	}

	ss := SwarmService{service, nil}
	if includeNodeIPInfo {
		ss.NodeInfo = c.getNodeInfo(service)
	}
	return &ss, nil
}

// ServicList returns a list of services
// When `includeNodeIPInfo` is true, return node info as well
func (c ServicClient) ServicList(ctx context.Context, includeNodeIPInfo bool) ([]SwarmService, error) {
	filter := filters.NewArgs()
	filter.Add("label", c.FilterLabel)
	services, err := c.DockerClient.ServiceList(ctx, types.ServiceListOptions{Filters: filter})
	if err != nil {
		return nil, err
	}
	swarmServices := []SwarmService{}
	for _, s := range services {
		ss := SwarmService{s, nil}
		if includeNodeIPInfo {
			ss.NodeInfo = c.getNodeInfo(s)
		}
		swarmServices = append(swarmServices, ss)
	}
	return swarmServices, nil
}

func (c ServicClient) getNodeInfo(ss swarm.Service) *NodeIPSet {

	nodeInfo := NodeIPSet{}
	filter := filters.NewArgs()
	filter.Add("desired-state", "running")
	filter.Add("service", ss.Spec.Name)
	taskList, err := c.DockerClient.TaskList(
		context.Background(), types.TaskListOptions{Filters: filter})
	if err != nil {
		return nil
	}

	networkName, ok := ss.Spec.Labels[c.ScrapeNetLabel]
	if !ok {
		return nil
	}

	nodeCache := map[string]string{}
	for _, task := range taskList {
		if len(task.NetworksAttachments) == 0 || len(task.NetworksAttachments[0].Addresses) == 0 {
			continue
		}
		var address string
		for _, networkAttach := range task.NetworksAttachments {
			if networkAttach.Network.Spec.Name == networkName && len(networkAttach.Addresses) > 0 {
				address = strings.Split(networkAttach.Addresses[0], "/")[0]
			}
		}

		if len(address) == 0 {
			continue
		}

		if nodeName, ok := nodeCache[task.NodeID]; ok {
			nodeInfo.Add(nodeName, address)
		} else {
			node, _, err := c.DockerClient.NodeInspectWithRaw(context.Background(), task.NodeID)
			if err != nil {
				continue
			}
			nodeInfo.Add(node.Description.Hostname, address)
			nodeCache[task.NodeID] = node.Description.Hostname
		}
	}

	if nodeInfo.Cardinality() == 0 {
		return nil
	}
	return &nodeInfo
}
