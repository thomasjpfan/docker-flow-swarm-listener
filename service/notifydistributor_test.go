package service

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type NotifyDistributorTestSuite struct {
	suite.Suite
	log *log.Logger
}

func TestNotifyDistributorUnitTestSuite(t *testing.T) {
	suite.Run(t, new(NotifyDistributorTestSuite))
}

func (s *NotifyDistributorTestSuite) SetupSuite() {
	s.log = log.New(os.Stdout, "", 0)
}

func (s *NotifyDistributorTestSuite) Test_NewNotifyDistributorFromStrings() {
	notifyD := newNotifyDistributorfromStrings(
		"http://host1:8080/recofigureservice,http://host2:8080/recofigureservice",
		"http://host1:8080/removeservice,http://host2:8080/removeservice",
		"http://host1:8080/reconfigurenode",
		"http://host2:8080/removenode",
		5, 10, s.log)

	s.Len(notifyD.NotifyEndpoints, 2)
	host1EP, ok := notifyD.NotifyEndpoints["host1:8080"]

	s.Require().True(ok)

	s.AssertEndpoints(
		host1EP,
		"http://host1:8080/recofigureservice",
		"http://host1:8080/removeservice",
		"http://host1:8080/reconfigurenode",
		"",
	)

	host2EP, ok := notifyD.NotifyEndpoints["host2:8080"]
	s.Require().True(ok)
	s.AssertEndpoints(
		host2EP,
		"http://host2:8080/recofigureservice",
		"http://host2:8080/removeservice",
		"",
		"http://host2:8080/removenode",
	)

}
func (s *NotifyDistributorTestSuite) Test_NewNotifyDistributorFromStrings_SeparateListeners() {
	notifyD := newNotifyDistributorfromStrings(
		"http://host1:8080/recofigure1",
		"http://host1:8080/removeservice",
		"http://host2:8080/reconfigurenode",
		"http://host2/removenode1,http://host2:8080/removenode2",
		5, 10, s.log)

	s.Len(notifyD.NotifyEndpoints, 3)
	host1EP, ok := notifyD.NotifyEndpoints["host1:8080"]
	s.Require().True(ok)

	s.AssertEndpoints(
		host1EP,
		"http://host1:8080/recofigure1",
		"http://host1:8080/removeservice",
		"",
		"",
	)

	host28080EP, ok := notifyD.NotifyEndpoints["host2:8080"]
	s.Require().True(ok)
	s.AssertEndpoints(
		host28080EP,
		"",
		"",
		"http://host2:8080/reconfigurenode",
		"http://host2:8080/removenode2",
	)

	host2EP, ok := notifyD.NotifyEndpoints["host2"]
	s.Require().True(ok)
	s.AssertEndpoints(
		host2EP,
		"",
		"",
		"",
		"http://host2/removenode1",
	)
}

func (s *NotifyDistributorTestSuite) Test_NewNotifyDistributorFromStrings_JustSwarmListeners() {
	notifyD := newNotifyDistributorfromStrings(
		"http://host1:8080/recofigure1",
		"http://host1:8080/removeservice", "", "",
		5, 10, s.log)

	s.Len(notifyD.NotifyEndpoints, 1)
	host1EP, ok := notifyD.NotifyEndpoints["host1:8080"]
	s.Require().True(ok)

	s.AssertEndpoints(
		host1EP,
		"http://host1:8080/recofigure1",
		"http://host1:8080/removeservice",
		"",
		"",
	)
}
func (s *NotifyDistributorTestSuite) Test_NewNotifyDistributorFromStrings_JustNodeListeners() {
	notifyD := newNotifyDistributorfromStrings(
		"", "",
		"http://host2:8080/reconfigurenode",
		"http://host2:8080/removenode1,http://host2/removenode2",
		5, 10, s.log)

	s.Len(notifyD.NotifyEndpoints, 2)
	host28080EP, ok := notifyD.NotifyEndpoints["host2:8080"]
	s.Require().True(ok)
	s.AssertEndpoints(
		host28080EP,
		"",
		"",
		"http://host2:8080/reconfigurenode",
		"http://host2:8080/removenode1",
	)

	host2EP, ok := notifyD.NotifyEndpoints["host2"]
	s.Require().True(ok)
	s.AssertEndpoints(
		host2EP,
		"",
		"",
		"",
		"http://host2/removenode2",
	)
}

func (s *NotifyDistributorTestSuite) Test_RunDistributesNotificationsToEndpoints() {
}

func (s *NotifyDistributorTestSuite) AssertEndpoints(endpoint NotifyEndpoint, serviceCreateAddr, serviceRemoveAddr, nodeCreateAddr, nodeRemoveAddr string) {
	if len(serviceCreateAddr) == 0 && len(serviceRemoveAddr) == 0 {
		s.Nil(endpoint.ServiceNotifier)
	} else {
		s.Require().NotNil(endpoint.ServiceNotifier)
		s.Equal(serviceCreateAddr, endpoint.ServiceNotifier.GetCreateAddr())
		s.Equal(serviceRemoveAddr, endpoint.ServiceNotifier.GetRemoveAddr())
	}
	if len(nodeCreateAddr) == 0 && len(nodeRemoveAddr) == 0 {
		s.Nil(endpoint.NodeNotifier)
	} else {
		s.Require().NotNil(endpoint.NodeNotifier)
		s.Equal(nodeCreateAddr, endpoint.NodeNotifier.GetCreateAddr())
		s.Equal(nodeRemoveAddr, endpoint.NodeNotifier.GetRemoveAddr())
	}

	if len(serviceCreateAddr) > 0 || len(serviceRemoveAddr) > 0 {
		s.NotNil(endpoint.ServiceChan)
	}
	if len(nodeCreateAddr) > 0 || len(nodeRemoveAddr) > 0 {
		s.NotNil(endpoint.NodeChan)
	}

}
