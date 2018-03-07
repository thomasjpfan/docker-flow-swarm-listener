package service

import (
	"context"

	"github.com/docker/docker/api/types/swarm"
	"github.com/stretchr/testify/mock"
)

type notificationSenderMock struct {
	mock.Mock
}

func (m *notificationSenderMock) Create(params string) error {
	args := m.Called(params)
	return args.Error(0)
}

func (m *notificationSenderMock) Remove(params string) error {
	args := m.Called(params)
	return args.Error(0)
}

func (m *notificationSenderMock) GetCreateAddr() string {
	args := m.Called()
	return args.String(0)
}

func (m *notificationSenderMock) GetRemoveAddr() string {
	args := m.Called()
	return args.String(0)
}

type swarmServiceListeningMock struct {
	mock.Mock
}

func (m *swarmServiceListeningMock) ListenForServiceEvents(eventChan chan<- Event) {
	m.Called(eventChan)
}

type swarmServiceInspector struct {
	mock.Mock
}

func (m *swarmServiceInspector) SwarmServiceInspect(serviceID string, includeNodeIPInfo bool) (*SwarmService, error) {
	args := m.Called(serviceID, includeNodeIPInfo)
	return args.Get(0).(*SwarmService), args.Error(1)
}

func (m *swarmServiceInspector) SwarmServiceList(ctx context.Context, includeNodeIPInfo bool) ([]SwarmService, error) {
	args := m.Called(ctx, includeNodeIPInfo)
	return args.Get(0).([]SwarmService), args.Error(1)
}

type swarmServiceCacherMock struct {
	mock.Mock
}

func (m *swarmServiceCacherMock) InsertAndCheck(ss SwarmServiceMini) bool {
	args := m.Called(ss)
	return args.Bool(0)
}

func (m *swarmServiceCacherMock) Delete(ID string) {
	m.Called(ID)
}

func (m *swarmServiceCacherMock) Get(ID string) (SwarmServiceMini, bool) {
	args := m.Called(ID)
	return args.Get(0).(SwarmServiceMini), args.Bool(1)
}

type nodeListeningMock struct {
	mock.Mock
}

func (m *nodeListeningMock) ListenForNodeEvents(eventChan chan<- Event) {
	m.Called(eventChan)
}

type nodeInspectorMock struct {
	mock.Mock
}

func (m *nodeInspectorMock) NodeInspect(nodeID string) (swarm.Node, error) {
	args := m.Called(nodeID)
	return args.Get(0).(swarm.Node), args.Error(1)
}

func (m *nodeInspectorMock) NodeList(ctx context.Context) ([]swarm.Node, error) {
	args := m.Called(ctx)
	return args.Get(0).([]swarm.Node), args.Error(1)
}

type nodeCacherMock struct {
	mock.Mock
}

func (m *nodeCacherMock) InsertAndCheck(n NodeMini) bool {
	args := m.Called(n)
	return args.Bool(0)
}

func (m *nodeCacherMock) Delete(ID string) {
	m.Called(ID)
}

func (m *nodeCacherMock) Get(ID string) (NodeMini, bool) {
	args := m.Called(ID)
	return args.Get(0).(NodeMini), args.Bool(1)
}

type notifyDistributorMock struct {
	mock.Mock
}

func (m *notifyDistributorMock) Run(serviceChan <-chan Notification, nodeChan <-chan Notification) {
	m.Called(serviceChan, nodeChan)
}

func (m *notifyDistributorMock) HasServiceListeners() bool {
	return m.Called().Bool(0)
}

func (m *notifyDistributorMock) HasNodeListeners() bool {
	return m.Called().Bool(0)
}
