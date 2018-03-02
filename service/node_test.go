package service

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type NodeInspectorTestSuite struct {
	suite.Suite
}

func TestNodeInspectorUnitTestSuite(t *testing.T) {
	suite.Run(t, new(NodeInspectorTestSuite))
}

func (s *NodeInspectorTestSuite) Test_NodeInspect() {

}

func (s *NodeInspectorTestSuite) Test_NodeInspect_Error() {

}

func (s *NodeInspectorTestSuite) Test_NodeList() {

}
