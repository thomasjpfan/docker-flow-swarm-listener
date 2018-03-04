package service

import (
	"testing"

	"github.com/docker/docker/api/types/swarm"
	"github.com/stretchr/testify/suite"
)

type ParametersTestSuite struct {
	suite.Suite
}

func TestParametersTestSuite(t *testing.T) {
	suite.Run(t, new(ParametersTestSuite))
}

func (s *ParametersTestSuite) Test_GetNodeParameters() {
	node := getNode(
		"hostname", "node123", swarm.NodeRoleManager,
		map[string]string{
			"com.df.wow":    "cats",
			"com.df.cows":   "fly",
			"com.df2.hello": "word"})

	params := GetNodeParameters(node)

	s.Equal("cats", params.Get("wow"))
	s.Equal("fly", params.Get("cows"))

	s.Equal("hostname", params.Get("hostname"))
	s.Equal("node123", params.Get("nodeID"))
	s.Equal("true", params.Get("manager"))
}