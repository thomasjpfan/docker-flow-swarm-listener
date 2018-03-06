package service

import (
	"encoding/json"
	"testing"

	"github.com/docker/docker/api/types/swarm"
	"github.com/stretchr/testify/suite"
)

type TypesTestSuite struct {
	suite.Suite
}

func TestTypesUnitTestSuite(t *testing.T) {
	suite.Run(t, new(TypesTestSuite))
}

func (s *TypesTestSuite) Test_EqualSwarmServiceMode_Global_ReturnsTrue() {
	a := swarm.ServiceMode{Global: &swarm.GlobalService{}}
	b := swarm.ServiceMode{Global: &swarm.GlobalService{}}

	s.True(EqualSwarmServiceMode(a, b))
	s.True(EqualSwarmServiceMode(b, a))
}

func (s *TypesTestSuite) Test_EqualSwarmServiceMode_Global_Replica_ReturnsFalse() {

	replicasNum := uint64(10)
	a := swarm.ServiceMode{Global: &swarm.GlobalService{}}
	b := swarm.ServiceMode{Replicated: &swarm.ReplicatedService{Replicas: &replicasNum}}

	s.False(EqualSwarmServiceMode(a, b))
	s.False(EqualSwarmServiceMode(b, a))
}
func (s *TypesTestSuite) Test_EqualSwarmServiceMode_Different_Replica_ReturnsFalse() {

	replicasNumA := uint64(4)
	replicasNumB := uint64(10)
	a := swarm.ServiceMode{Replicated: &swarm.ReplicatedService{Replicas: &replicasNumA}}
	b := swarm.ServiceMode{Replicated: &swarm.ReplicatedService{Replicas: &replicasNumB}}

	s.False(EqualSwarmServiceMode(a, b))
	s.False(EqualSwarmServiceMode(b, a))
}

func (s *TypesTestSuite) Test_EqualSwarmServiceMode_Same_Replica_ReturnsTrue() {

	replicasNumA := uint64(4)
	replicasNumB := uint64(4)
	a := swarm.ServiceMode{Replicated: &swarm.ReplicatedService{Replicas: &replicasNumA}}
	b := swarm.ServiceMode{Replicated: &swarm.ReplicatedService{Replicas: &replicasNumB}}

	s.True(EqualSwarmServiceMode(a, b))
	s.True(EqualSwarmServiceMode(b, a))
}

func (s *TypesTestSuite) Test_EqualMapStringString_SameKeys_DifferentValue() {
	a := map[string]string{"k1": "v1", "k2": "v2"}
	b := map[string]string{"k1": "v2", "k2": "v2"}

	s.False(EqualMapStringString(a, b))
	s.False(EqualMapStringString(b, a))
}

func (s *TypesTestSuite) Test_EqualMapStringString_SameKeys_SameValue() {

	a := map[string]string{"k1": "v1", "k2": "v2"}
	b := map[string]string{"k1": "v1", "k2": "v2"}

	s.True(EqualMapStringString(a, b))
	s.True(EqualMapStringString(b, a))
}

func (s *TypesTestSuite) Test_EqualMapStringString_DifferentKeys() {

	a := map[string]string{"k1": "v1", "k2": "v2"}
	b := map[string]string{"k1": "v1", "k3": "v2"}

	s.False(EqualMapStringString(a, b))
	s.False(EqualMapStringString(b, a))
}

func (s *TypesTestSuite) Test_EqualMapStringString_DifferentNumberOfValues() {

	a := map[string]string{"k1": "v1", "k2": "v2", "k3": "v3"}
	b := map[string]string{"k1": "v1", "k2": "v2"}

	s.False(EqualMapStringString(a, b))
	s.False(EqualMapStringString(b, a))
}

func (s *TypesTestSuite) Test_Cardinality_DifferentElms() {
	a := NodeIPSet{}
	a.Add("node-1", "1.0.0.1")
	a.Add("node-2", "1.0.1.1")
	s.Equal(2, a.Cardinality())
}

func (s *TypesTestSuite) Test_Cardinality_RepeatElems() {
	a := NodeIPSet{}
	a.Add("node-1", "1.0.0.1")
	a.Add("node-1", "1.0.0.1")
	s.Equal(1, a.Cardinality())
}

func (s *TypesTestSuite) Test_NodeIPSetEqual_RepeatElems() {
	a := NodeIPSet{}
	a.Add("node-1", "1.0.0.1")
	a.Add("node-1", "1.0.0.1")
	b := NodeIPSet{}
	b.Add("node-1", "1.0.0.1")
	s.True(a.Equal(b))
}

func (s *TypesTestSuite) Test_NodeIPSetEqual_LenUnequal() {
	a := NodeIPSet{}
	a.Add("node-1", "1.0.0.1")
	a.Add("node-2", "1.0.1.1")
	b := NodeIPSet{}
	b.Add("node-1", "1.0.0.1")
	b.Add("node-2", "1.0.1.1")
	b.Add("node-2", "1.0.1.2")
	s.False(a.Equal(b))
}

func (s *TypesTestSuite) Test_NodeIPSetEqual_EqualSets() {
	a := NodeIPSet{}
	a.Add("node-1", "1.0.0.1")
	a.Add("node-2", "1.0.1.1")
	b := NodeIPSet{}
	b.Add("node-1", "1.0.0.1")
	b.Add("node-2", "1.0.1.1")
	s.True(a.Equal(b))
}

func (s *TypesTestSuite) Test_NodeIPSetEqual_AddrNotEqual() {
	a := NodeIPSet{}
	a.Add("node-1", "1.0.0.1")
	a.Add("node-2", "1.0.1.1")
	b := NodeIPSet{}
	b.Add("node-1", "1.0.0.1")
	b.Add("node-2", "1.0.1.2")
	s.False(a.Equal(b))
}

func (s *TypesTestSuite) Test_NodeIPSetEqual_NodeNameNotEqual() {
	a := NodeIPSet{}
	a.Add("node-1", "1.0.0.1")
	a.Add("node-2", "1.0.1.1")
	b := NodeIPSet{}
	b.Add("node-1", "1.0.0.1")
	b.Add("node-1", "1.0.1.1")
	s.False(a.Equal(b))
}

func (s *TypesTestSuite) Test_NodeIPSetEqual_EmptySets() {
	a := NodeIPSet{}
	b := NodeIPSet{}
	s.True(a.Equal(b))
}

func (s *TypesTestSuite) Test_NodeIPSetEqual_OneEmpty() {
	a := NodeIPSet{}
	b := NodeIPSet{}
	b.Add("node-1", "1.0.0.1")
	b.Add("node-1", "1.0.1.1")
	s.False(a.Equal(b))
}

func (s *TypesTestSuite) Test_NodeIPMarshallJSON_EmptySet() {
	a := NodeIPSet{}
	b, err := json.Marshal(a)
	s.Require().NoError(err)

	s.Equal([]byte("[]"), b)
}

func (s *TypesTestSuite) Test_NodeIP_JSONCycle() {
	a := NodeIPSet{}
	a.Add("node-1", "1.0.0.1")
	a.Add("node-2", "1.0.1.1")
	by, err := json.Marshal(a)
	s.Require().NoError(err)

	i := NodeIPSet{}
	err = json.Unmarshal(by, &i)
	s.Require().NoError(err)

	s.True(a.Equal(i))
}

func (s *TypesTestSuite) Test_NodeIPSet_Add() {
	a := NodeIPSet{}
	a.Add("node-1", "1.0.0.1")
	a.Add("node-1", "1.0.1.1")
	b := NodeIPSet{}
	b.Add("node-1", "1.0.0.1")
	b.Add("node-1", "1.0.1.1")

	s.True(a.Equal(b))
}
