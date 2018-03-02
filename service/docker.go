package service

import (
	"context"
	"os"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

var dockerApiVersion string = "v1.36"

// NodeInspector is able to inspect a swarm node
type NodeInspector interface {
	NodeInspect(ctx context.Context, nodeID string) (swarm.Node, error)
}

// DockerClient wraps the docker client and implements the interfaces:
// `NodeInspector`
type DockerClient struct {
	*client.Client
}

// NewDockerClientFromEnv returns a `DockerClient` struct using environment variable
// `DF_DOCKER_HOST` for the host
func NewDockerClientFromEnv() (DockerClient, error) {
	host := "unix:///var/run/docker.sock"
	if len(os.Getenv("DF_DOCKER_HOST")) > 0 {
		host = os.Getenv("DF_DOCKER_HOST")
	}
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	dc, err := client.NewClient(host, dockerApiVersion, nil, defaultHeaders)
	return DockerClient{dc}, err
}
