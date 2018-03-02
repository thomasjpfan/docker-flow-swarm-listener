package service

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type EventListeningMock struct {
	mock.Mock
}

type NodeInspectorMock struct {
	mock.Mock
}

type NodeEventNotifingMock struct {
	mock.Mock
}

type WatcherTestSuite struct {
	suite.Suite
	elMock  EventListeningMock
	niMock  NodeInspectorMock
	nenMock NodeEventNotifingMock
}

func TestWatcherUnitTestSuite(t *testing.T) {
	s := new(WatcherTestSuite)
	s.elMock = EventListeningMock{}
	s.niMock = NodeInspectorMock{}
	s.nenMock = NodeEventNotifingMock{}

	suite.Run(t, s)
}
