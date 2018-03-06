package service

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ParametersTestSuite struct {
	suite.Suite
}

func TestParametersTestSuite(t *testing.T) {
	suite.Run(t, new(ParametersTestSuite))
}

func (s *ParametersTestSuite) GetNodeMiniParameters() {
}

func (s *ParametersTestSuite) GetSwarmServiceMiniParameters() {
}
