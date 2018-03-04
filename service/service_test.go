package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ServicClientTestSuite struct {
	suite.Suite
	SClient *ServicClient
	Util1ID string
	Util2ID string
	Util3ID string
	Util4ID string
}

func TestServicClientUnitTestSuite(t *testing.T) {
	suite.Run(t, new(ServicClientTestSuite))
}

func (s *ServicClientTestSuite) SetupSuite() {
	c, err := NewDockerClientFromEnv()
	s.Require().NoError(err)
	s.SClient = NewServicClient(c, "com.df.notify=true", "com.df.scrapeNetwork")

	createTestOverlayNetwork("util-network")
	createTestService("util-1", []string{"com.df.notify=true", "com.df.scrapeNetwork=util-network"}, false, "", "util-network")
	createTestService("util-2", []string{}, false, "", "util-network")
	createTestService("util-3", []string{"com.df.notify=true"}, true, "", "util-network")
	createTestService("util-4", []string{"com.df.notify=true", "com.df.scrapeNetwork=util-network"}, false, "2", "util-network")

	time.Sleep(time.Second)
	ID1, err := getServiceID("util-1")
	s.Require().NoError(err)
	s.Util1ID = ID1

	ID2, err := getServiceID("util-2")
	s.Require().NoError(err)
	s.Util2ID = ID2

	ID3, err := getServiceID("util-3")
	s.Require().NoError(err)
	s.Util3ID = ID3

	ID4, err := getServiceID("util-4")
	s.Require().NoError(err)
	s.Util4ID = ID4
}

func (s *ServicClientTestSuite) TearDownSuite() {
	removeTestService("util-1")
	removeTestService("util-2")
	removeTestService("util-3")
	removeTestService("util-4")
	removeTestNetwork("util-network")
}

func (s *ServicClientTestSuite) Test_ServicInspect_NodeInfo_UndefinedScrapeNetwork() {

	util3Service, err := s.SClient.ServicInspect(s.Util3ID, true)
	s.Require().NoError(err)
	s.Require().NotNil(util3Service)

	s.Equal(s.Util3ID, util3Service.ID)
	s.Nil(util3Service.NodeInfo)
}

func (s *ServicClientTestSuite) Test_ServiceList_Filtered() {

	util2Service, err := s.SClient.ServicInspect(s.Util2ID, false)
	s.Require().NoError(err)
	s.Nil(util2Service)

}

func (s *ServicClientTestSuite) Test_ServicInspect_NodeInfo_OneReplica() {
	util1Service, err := s.SClient.ServicInspect(s.Util1ID, true)
	s.Require().NoError(err)
	s.Require().NotNil(util1Service)

	s.Equal(s.Util1ID, util1Service.ID)
	s.Require().NotNil(util1Service.NodeInfo)

	nodeInfo := *util1Service.NodeInfo
	s.Require().Len(nodeInfo, 1)
}

func (s *ServicClientTestSuite) Test_ServicInspect_NodeInfo_TwoReplica() {

	util4Service, err := s.SClient.ServicInspect(s.Util4ID, true)
	s.Require().NoError(err)
	s.Require().NotNil(util4Service)

	s.Equal(s.Util4ID, util4Service.ID)
	s.Require().NotNil(util4Service.NodeInfo)

	nodeInfo := *util4Service.NodeInfo
	s.Require().Len(nodeInfo, 2)
}

func (s *ServicClientTestSuite) Test_ServicInspect_IncorrectName() {
	_, err := s.SClient.ServicInspect("cowsfly", true)
	s.Error(err)
}

func (s *ServicClientTestSuite) Test_ServicList_NodeInfo() {
	services, err := s.SClient.ServicList(context.Background(), true)
	s.Require().NoError(err)
	s.Len(services, 3)

	for _, ss := range services {
		if ss.Spec.Name == "util-1" || ss.Spec.Name == "util-4" {
			s.NotNil(ss.NodeInfo)
		} else {
			s.Nil(ss.NodeInfo)
		}
	}
}

func (s *ServicClientTestSuite) Test_ServicList_NoNodeInfo() {

	services, err := s.SClient.ServicList(context.Background(), false)
	s.Require().NoError(err)
	s.Len(services, 3)

	for _, ss := range services {
		s.Nil(ss.NodeInfo)
	}
}
