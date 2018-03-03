package service

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type EventNodeNotifingMock struct {
	mock.Mock
}

type WatcherTestSuite struct {
	suite.Suite
	elMock  eventListeningMock
	niMock  nodeInspectorMock
	nenMock EventNodeNotifingMock
}

func TestWatcherUnitTestSuite(t *testing.T) {
	s := new(WatcherTestSuite)
	s.elMock = eventListeningMock{}
	s.niMock = nodeInspectorMock{}
	s.nenMock = EventNodeNotifingMock{}

	suite.Run(t, s)
}

func (s *WatcherTestSuite) Test_WatchServices() {

}
