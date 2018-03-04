package service

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServicCacheTestSuite struct {
	suite.Suite
}

func TestServicCacheUnitTestSuite(t *testing.T) {
	suite.Run(t, new(ServicCacheTestSuite))
}

func (s *ServicCacheTestSuite) Test_InsertAndCheck_NewService() {

}

func (s *ServicCacheTestSuite) Test_InsertAndCheck_NewServiceGlobal() {

}

func (s *ServicCacheTestSuite) Test_InsertAndCheck_UpdatedService_SameLabel() {

}

func (s *ServicCacheTestSuite) Test_InsertAndCheck_UpdatedService_NewLabel() {

}

func (s *ServicCacheTestSuite) Test_InsertAndCheck_UpdatedService_UpdateLabel_WithPrefix() {

}

func (s *ServicCacheTestSuite) Test_InsertAndCheck_UpdatedService_UpdateLabel_WithoutPrefix() {

}

func (s *ServicCacheTestSuite) Test_InsertAndCheck_UpdatedService_SameReplicas() {

}

func (s *ServicCacheTestSuite) Test_InsertAndCheck_UpdatedService_IncreasedReplicas() {

}

func (s *ServicCacheTestSuite) Test_InsertAndCheck_UpdatedService_ReplicasSetToZero() {

}

func (s *ServicCacheTestSuite) Test_InsertAndCheck_UpdatedService_NewNodeInfo() {

}

func (s *ServicCacheTestSuite) Test_InsertAndCheck_RemovedService() {

}
