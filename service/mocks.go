package service

import (
	"github.com/docker/docker/api/types/swarm"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

type eventListeningMock struct {
	mock.Mock
}

func (m *eventListeningMock) ListenForEvents() (<-chan Event, <-chan error) {
	args := m.Called()
	return args.Get(0).(chan Event), args.Get(1).(chan error)
}

type nodeInspectorMock struct {
	mock.Mock
}

func (m *nodeInspectorMock) NodeInspect(
	ctx context.Context, nodeID string) (swarm.Node, error) {
	args := m.Called(ctx, nodeID)
	return args.Get(0).(swarm.Node), args.Error(1)
}

func (m *nodeInspectorMock) NodeList(
	ctx context.Context) ([]swarm.Node, error) {
	args := m.Called(ctx)
	return args.Get(0).([]swarm.Node), args.Error(1)
}

type eventNodeNotifingMock struct {
	mock.Mock
}

// func (m *eventNodeNotifingMock) NotifyCreateNodes(nodes []swarm.Node, retry, retryInterval int) error {
// 	args := m.Called(nodes, retry, retryInterval)
// }

func (m *eventNodeNotifingMock) HasListeners() bool {
	args := m.Called()
	return args.Bool(0)
}
