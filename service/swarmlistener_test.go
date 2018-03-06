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
}

func TestWatcherUnitTestSuite(t *testing.T) {
	s := new(WatcherTestSuite)

	suite.Run(t, s)
}

func (s *WatcherTestSuite) Test_WatchServices() {

}
