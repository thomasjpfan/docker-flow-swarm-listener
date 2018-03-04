package service

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type NodeCacheTestSuite struct {
	suite.Suite
}

func TestNodeCacheUnitTestSuite(t *testing.T) {
	suite.Run(t, new(NodeCacheTestSuite))
}

func (s *NodeCacheTestSuite) Test_InsertAndCheck_NewNode() {

}

func (s *NodeCacheTestSuite) Test_InsertAndCheck_RemovedNode() {

}

func (s *NodeCacheTestSuite) Test_InsertAndCheck_UpdateNode_SameLabel() {

}

func (s *NodeCacheTestSuite) Test_InsertAndCheck_UpdateNode_NewLabel() {

}

func (s *NodeCacheTestSuite) Test_InsertAndCheck_UpdateNode_UpdateLabel_WithPrefix() {

}

func (s *NodeCacheTestSuite) Test_InsertAndCheck_UpdateNode_UpdateLabel_WithOutPrefix() {

}

func (s *NodeCacheTestSuite) Test_InsertAndCheck_UpdateNode_ChangeRole() {

}

func (s *NodeCacheTestSuite) Test_InsertAndCheck_UpdateNode_ChangeStatus() {

}

func (s *NodeCacheTestSuite) Test_InsertAndCheck_UpdateNode_ChangeAvailability() {

}

func (s *NodeCacheTestSuite) Test_InsertAndCheck_UpdateNode_ChangeVersion() {

}
