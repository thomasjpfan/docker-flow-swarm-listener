package main

import (
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ArgsTestSuite struct {
	suite.Suite
	serviceName string
}

func TestArgsUnitTestSuite(t *testing.T) {
	s := new(ArgsTestSuite)

	suite.Run(t, s)
}

// GetArgs

func (s *ArgsTestSuite) Test_GetArgs_ReturnsDefaultValues() {
	args := getArgs()

	s.Equal(5, args.Interval)
	s.Equal(1, args.Retry)
	s.Equal(0, args.RetryInterval)
}

func (s *ArgsTestSuite) Test_GetArgs_ReturnsRetryFromEnv() {
	expected := rand.Int()
	intervalOrig := os.Getenv("DF_RETRY")
	defer func() { os.Setenv("DF_RETRY", intervalOrig) }()
	os.Setenv("DF_RETRY", strconv.Itoa(expected))

	args := getArgs()

	s.Equal(expected, args.Retry)
}

func (s *ArgsTestSuite) Test_GetArgs_ReturnsRetryIntervalFromEnv() {
	expected := rand.Int()
	intervalOrig := os.Getenv("DF_RETRY_INTERVAL")
	defer func() { os.Setenv("DF_RETRY_INTERVAL", intervalOrig) }()
	os.Setenv("DF_RETRY_INTERVAL", strconv.Itoa(expected))

	args := getArgs()

	s.Equal(expected, args.RetryInterval)
}
