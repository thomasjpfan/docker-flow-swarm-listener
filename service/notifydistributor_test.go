package service

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type NotifyDistributorTestSuite struct {
	suite.Suite
}

func NotifyDistributorUnitTestSuite(t *testing.T) {
	suite.Run(t, new(NotifyDistributorTestSuite))
}
