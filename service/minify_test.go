package service

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type MinifyUnitTestSuite struct {
	suite.Suite
}

func TestMinifyUnitTest(t *testing.T) {
	suite.Run(t, new(MinifyUnitTestSuite))
}

func (t *MinifyUnitTestSuite) Test_MinifyNode() {

}

func (t *MinifyUnitTestSuite) Test_MinifySwarmService() {

}
