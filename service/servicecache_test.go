package service

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type SwarmServiceCacheTestSuite struct {
	suite.Suite
}

func TestSwarmServiceCacheUnitTestSuite(t *testing.T) {
	suite.Run(t, new(SwarmServiceCacheTestSuite))
}

func (s *SwarmServiceCacheTestSuite) Test_InsertAndCheck_NewService() {

}

func (s *SwarmServiceCacheTestSuite) Test_InsertAndCheck_NewServiceGlobal() {

}

func (s *SwarmServiceCacheTestSuite) Test_InsertAndCheck_UpdatedService_SameLabel() {

}

func (s *SwarmServiceCacheTestSuite) Test_InsertAndCheck_UpdatedService_NewLabel() {

}

func (s *SwarmServiceCacheTestSuite) Test_InsertAndCheck_UpdatedService_UpdateLabel_WithPrefix() {

}

func (s *SwarmServiceCacheTestSuite) Test_InsertAndCheck_UpdatedService_UpdateLabel_WithoutPrefix() {

}

func (s *SwarmServiceCacheTestSuite) Test_InsertAndCheck_UpdatedService_SameReplicas() {

}

func (s *SwarmServiceCacheTestSuite) Test_InsertAndCheck_UpdatedService_IncreasedReplicas() {

}

func (s *SwarmServiceCacheTestSuite) Test_InsertAndCheck_UpdatedService_ReplicasSetToZero() {

}

func (s *SwarmServiceCacheTestSuite) Test_InsertAndCheck_UpdatedService_NewNodeInfo() {

}

func (s *SwarmServiceCacheTestSuite) Test_InsertAndCheck_RemovedService() {

}
