package service

import (
	"testing"

	"github.com/docker/docker/api/types/swarm"
	"github.com/stretchr/testify/suite"
)

type SwarmServiceCacheTestSuite struct {
	suite.Suite
	Cache  *SwarmServiceCache
	SSMini SwarmServiceMini
}

func TestSwarmServiceCacheUnitTestSuite(t *testing.T) {
	suite.Run(t, new(SwarmServiceCacheTestSuite))
}

func (s *SwarmServiceCacheTestSuite) SetupTest() {
	s.Cache = NewSwarmServiceCache()
	s.SSMini = getNewSwarmServiceMini()
}

func (s *SwarmServiceCacheTestSuite) Test_InsertAndCheck_NewService_ReturnsTrue() {
	isUpdated := s.Cache.InsertAndCheck(s.SSMini)
	s.True(isUpdated)

	s.AssertInCache()
}

func (s *SwarmServiceCacheTestSuite) Test_InsertAndCheck_NewServiceGlobal_ReturnsTrue() {

	s.SSMini.Mode = swarm.ServiceMode{
		Global: &swarm.GlobalService{},
	}

	isUpdated := s.Cache.InsertAndCheck(s.SSMini)
	s.True(isUpdated)
	s.AssertInCache()
}

func (s *SwarmServiceCacheTestSuite) Test_InsertAndCheck_SameService_ReturnsFalse() {

	isUpdated := s.Cache.InsertAndCheck(s.SSMini)
	s.True(isUpdated)
	s.AssertInCache()

	isUpdated = s.Cache.InsertAndCheck(s.SSMini)
	s.False(isUpdated)
	s.AssertInCache()
}

func (s *SwarmServiceCacheTestSuite) Test_InsertAndCheck_NewLabel_ReturnsTrue() {
	isUpdated := s.Cache.InsertAndCheck(s.SSMini)
	s.True(isUpdated)
	s.AssertInCache()

	s.SSMini.Labels["com.df.whatisthis"] = "howareyou"

	isUpdated = s.Cache.InsertAndCheck(s.SSMini)
	s.True(isUpdated)
	s.AssertInCache()
}

func (s *SwarmServiceCacheTestSuite) Test_InsertAndCheck_NewLabel_SameKey_ReturnsTrue() {
	isUpdated := s.Cache.InsertAndCheck(s.SSMini)
	s.True(isUpdated)
	s.AssertInCache()

	s.SSMini.Labels["com.df.hello"] = "sf"

	isUpdated = s.Cache.InsertAndCheck(s.SSMini)
	s.True(isUpdated)
	s.AssertInCache()
}

func (s *SwarmServiceCacheTestSuite) Test_InsertAndCheck_IncreasedReplicas_ReturnsTrue() {

	isUpdated := s.Cache.InsertAndCheck(s.SSMini)
	s.True(isUpdated)
	s.AssertInCache()

	newReplicas := uint64(5)
	s.SSMini.Mode.Replicated.Replicas = &newReplicas

	isUpdated = s.Cache.InsertAndCheck(s.SSMini)
	s.True(isUpdated)
	s.AssertInCache()
}

func (s *SwarmServiceCacheTestSuite) Test_InsertAndCheck_ReplicasDescToZero_ReturnsTrue() {

	isUpdated := s.Cache.InsertAndCheck(s.SSMini)
	s.True(isUpdated)
	s.AssertInCache()

	newReplicas := uint64(0)
	s.SSMini.Mode.Replicated.Replicas = &newReplicas

	isUpdated = s.Cache.InsertAndCheck(s.SSMini)
	s.True(isUpdated)
	s.AssertInCache()
}

func (s *SwarmServiceCacheTestSuite) Test_InsertAndCheck_NewNodeInfo_ReturnsTrue() {

	isUpdated := s.Cache.InsertAndCheck(s.SSMini)
	s.True(isUpdated)
	s.AssertInCache()

	nodeSet := NodeIPSet{}
	nodeSet.Add("node-3", "1.0.2.1")

	isUpdated = s.Cache.InsertAndCheck(s.SSMini)
	s.True(isUpdated)
	s.AssertInCache()
}

func (s *SwarmServiceCacheTestSuite) Test_GetAndRemove_InCache_ReturnsSwarmServiceMini_RemovesFromCache() {

	isUpdated := s.Cache.InsertAndCheck(s.SSMini)
	s.True(isUpdated)
	s.AssertInCache()

	removedSSMini, ok := s.Cache.GetAndRemove(s.SSMini.ID)
	s.True(ok)
	s.AssertNotInCache()
	s.Equal(s.SSMini, removedSSMini)

}

func (s *SwarmServiceCacheTestSuite) Test_GetAndRemove_NotInCache_ReturnsFalse() {

	_, ok := s.Cache.GetAndRemove(s.SSMini.ID)
	s.False(ok)
	s.AssertNotInCache()
}

func (s *SwarmServiceCacheTestSuite) AssertInCache() {
	ss, ok := s.Cache.get(s.SSMini.ID)
	s.True(ok)
	s.Equal(s.SSMini, ss)
}

func (s *SwarmServiceCacheTestSuite) AssertNotInCache() {
	_, ok := s.Cache.get(s.SSMini.ID)
	s.False(ok)
}
