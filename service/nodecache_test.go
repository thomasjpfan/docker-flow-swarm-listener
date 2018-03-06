package service

import (
	"testing"

	"github.com/docker/docker/api/types/swarm"
	"github.com/stretchr/testify/suite"
)

type NodeCacheTestSuite struct {
	suite.Suite
	Cache *NodeCache
	NMini NodeMini
}

func TestNodeCacheUnitTestSuite(t *testing.T) {
	suite.Run(t, new(NodeCacheTestSuite))
}

func (s *NodeCacheTestSuite) SetupTest() {
	s.Cache = NewNodeCache()
	s.NMini = getNewNodeMini()
}

func (s *NodeCacheTestSuite) Test_InsertAndCheck_NewNode_ReturnsTrue() {
	isUpdated := s.Cache.InsertAndCheck(s.NMini)
	s.True(isUpdated)
	s.AssertInCache()
}

func (s *NodeCacheTestSuite) Test_InsertAndCheck_SameLabel_ReturnsFalse() {
	isUpdated := s.Cache.InsertAndCheck(s.NMini)
	s.True(isUpdated)
	s.AssertInCache()

	isUpdated = s.Cache.InsertAndCheck(s.NMini)
	s.False(isUpdated)
	s.AssertInCache()
}

func (s *NodeCacheTestSuite) Test_InsertAndCheck_NewNodeLabel_ReturnsTrue() {
	isUpdated := s.Cache.InsertAndCheck(s.NMini)
	s.True(isUpdated)
	s.AssertInCache()

	s.NMini.NodeLabels["com.df.wow2"] = "yup2"

	isUpdated = s.Cache.InsertAndCheck(s.NMini)
	s.True(isUpdated)
	s.AssertInCache()
}

func (s *NodeCacheTestSuite) Test_InsertAndCheck_UpdateNodeLabel_ReturnsTrue() {
	isUpdated := s.Cache.InsertAndCheck(s.NMini)
	s.True(isUpdated)
	s.AssertInCache()

	s.NMini.NodeLabels["com.df.wow"] = "yup2"

	isUpdated = s.Cache.InsertAndCheck(s.NMini)
	s.True(isUpdated)
	s.AssertInCache()
}

func (s *NodeCacheTestSuite) Test_InsertAndCheck_NewEngineLabel_ReturnsTrue() {
	isUpdated := s.Cache.InsertAndCheck(s.NMini)
	s.True(isUpdated)
	s.AssertInCache()

	s.NMini.NodeLabels["com.df.mars"] = "far"

	isUpdated = s.Cache.InsertAndCheck(s.NMini)
	s.True(isUpdated)
	s.AssertInCache()
}

func (s *NodeCacheTestSuite) Test_InsertAndCheck_UpdateEngineLabel_ReturnsTrue() {
	isUpdated := s.Cache.InsertAndCheck(s.NMini)
	s.True(isUpdated)
	s.AssertInCache()

	s.NMini.NodeLabels["com.df.world"] = "flat"

	isUpdated = s.Cache.InsertAndCheck(s.NMini)
	s.True(isUpdated)
	s.AssertInCache()
}

func (s *NodeCacheTestSuite) Test_InsertAndCheck_ChangeRole_ReturnsTrue() {
	isUpdated := s.Cache.InsertAndCheck(s.NMini)
	s.True(isUpdated)
	s.AssertInCache()

	s.NMini.Role = swarm.NodeRoleManager

	isUpdated = s.Cache.InsertAndCheck(s.NMini)
	s.True(isUpdated)
	s.AssertInCache()
}

func (s *NodeCacheTestSuite) Test_InsertAndCheck_ChangeState_ReturnsTrue() {

	isUpdated := s.Cache.InsertAndCheck(s.NMini)
	s.True(isUpdated)
	s.AssertInCache()

	s.NMini.State = swarm.NodeStateDown

	isUpdated = s.Cache.InsertAndCheck(s.NMini)
	s.True(isUpdated)
	s.AssertInCache()
}

func (s *NodeCacheTestSuite) Test_InsertAndCheck_ChangeAvailability_ReturnsTrue() {

	isUpdated := s.Cache.InsertAndCheck(s.NMini)
	s.True(isUpdated)
	s.AssertInCache()

	s.NMini.Availability = swarm.NodeAvailabilityPause

	isUpdated = s.Cache.InsertAndCheck(s.NMini)
	s.True(isUpdated)
	s.AssertInCache()
}

func (s *NodeCacheTestSuite) Test_InsertAndCheck_ChangeIndexVersion_ReturnsTrue() {

	isUpdated := s.Cache.InsertAndCheck(s.NMini)
	s.True(isUpdated)
	s.AssertInCache()

	s.NMini.VersionIndex = uint64(4)

	isUpdated = s.Cache.InsertAndCheck(s.NMini)
	s.True(isUpdated)
	s.AssertInCache()
}

func (s *NodeCacheTestSuite) Test_GetAndRemove_InCache_ReturnsNodeMini_RemovesFromCache() {

	isUpdated := s.Cache.InsertAndCheck(s.NMini)
	s.True(isUpdated)
	s.AssertInCache()

	removedNMini, ok := s.Cache.GetAndRemove(s.NMini.ID)
	s.True(ok)
	s.AssertNotInCache()
	s.Equal(s.NMini, removedNMini)
}

func (s *NodeCacheTestSuite) Test_GetAndRemove_NotInCache_ReturnsFalse() {
	_, ok := s.Cache.GetAndRemove(s.NMini.ID)
	s.False(ok)
	s.AssertNotInCache()
}

func (s *NodeCacheTestSuite) AssertInCache() {
	ss, ok := s.Cache.get(s.NMini.ID)
	s.True(ok)
	s.Equal(s.NMini, ss)
}

func (s *NodeCacheTestSuite) AssertNotInCache() {
	_, ok := s.Cache.get(s.NMini.ID)
	s.False(ok)
}
