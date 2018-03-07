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
	host1EP, ok := notifyD.NotifyEndpoints["host1"]

	s.Require().True(ok)

	s.AssertEndpoints(
		host1EP,
		[]string{"http://host1:8080/recofigureservice"},
		[]string{"http://host1:8080/removeservice"},
		[]string{"http://host1:8080/reconfigurenode"},
		[]string{},
	)

	host2EP, ok := notifyD.NotifyEndpoints["host2"]
	s.Require().True(ok)
	s.AssertEndpoints(
		host2EP,
		[]string{"http://host2:8080/recofigureservice"},
		[]string{"http://host2:8080/removeservice"},
		[]string{},
		[]string{"http://host2:8080/removenode"},
	)

}
func (s *NotifyDistributorTestSuite) Test_NewNotifyDistributorFromStrings_SeparateListeners() {
	notifyD := newNotifyDistributorfromStrings(
		"http://host1:8080/recofigure1,http://host1:8080/recofigure2",
		"http://host1:8080/removeservice",
		"http://host2:8080/reconfigurenode",
		"http://host2/removenode1,http://host2/removenode2",
		5, 10, s.log)

	s.Len(notifyD.NotifyEndpoints, 2)
	host1EP, ok := notifyD.NotifyEndpoints["host1"]
	s.Require().True(ok)

	s.AssertEndpoints(
		host1EP,
		[]string{"http://host1:8080/recofigure1", "http://host1:8080/recofigure2"},
		[]string{"http://host1:8080/removeservice"},
		[]string{},
		[]string{},
	)

	host2EP, ok := notifyD.NotifyEndpoints["host2"]
	s.Require().True(ok)
	s.AssertEndpoints(
		host2EP,
		[]string{},
		[]string{},
		[]string{"http://host2:8080/reconfigurenode"},
		[]string{"http://host2/removenode1", "http://host2/removenode2"},
	)
}

func (s *NotifyDistributorTestSuite) Test_NewNotifyDistributorFromStrings_JustSwarmListeners() {
	notifyD := newNotifyDistributorfromStrings(
		"http://host1:8080/recofigure1,http://host1:8080/recofigure2",
		"http://host1:8080/removeservice", "", "",
		5, 10, s.log)

	s.Len(notifyD.NotifyEndpoints, 1)
	host1EP, ok := notifyD.NotifyEndpoints["host1"]
	s.Require().True(ok)

	s.AssertEndpoints(
		host1EP,
		[]string{"http://host1:8080/recofigure1", "http://host1:8080/recofigure2"},
		[]string{"http://host1:8080/removeservice"},
		[]string{},
		[]string{},
	)
}
func (s *NotifyDistributorTestSuite) Test_NewNotifyDistributorFromStrings_JustNodeListeners() {
	notifyD := newNotifyDistributorfromStrings(
		"", "",
		"http://host2:8080/reconfigurenode",
		"http://host2/removenode1,http://host2/removenode2",
		5, 10, s.log)

	s.Len(notifyD.NotifyEndpoints, 1)
	host2EP, ok := notifyD.NotifyEndpoints["host2"]
	s.Require().True(ok)
	s.AssertEndpoints(
		host2EP,
		[]string{},
		[]string{},
		[]string{"http://host2:8080/reconfigurenode"},
		[]string{"http://host2/removenode1", "http://host2/removenode2"},
	)
}

func (s *NotifyDistributorTestSuite) Test_RunDistributesNotificationsToEndpoints() {
}

func (s *NotifyDistributorTestSuite) AssertEndpoints(endpoint NotifyEndpoint, serviceCreateAddrs, serviceRemoveAddrs, nodeCreateAddrs, nodeRemoveAddrs []string) {
	if len(serviceCreateAddrs) == 0 && len(serviceRemoveAddrs) == 0 {
		s.Nil(endpoint.ServiceNotifier)
	} else {
		s.Require().NotNil(endpoint.ServiceNotifier)

		epServiceCreateAddrs := endpoint.ServiceNotifier.GetCreateAddrs()
		epServiceRemoveAddrs := endpoint.ServiceNotifier.GetRemoveAddrs()
		s.Assert().EqualValues(serviceCreateAddrs, epServiceCreateAddrs)
		s.Assert().EqualValues(serviceRemoveAddrs, epServiceRemoveAddrs)
	}

	if len(nodeCreateAddrs) == 0 && len(nodeRemoveAddrs) == 0 {
		s.Nil(endpoint.NodeNotifier)
	} else {
		s.Require().NotNil(endpoint.NodeNotifier)

		epNodeCreateAddrs := endpoint.NodeNotifier.GetCreateAddrs()
		epNodeRemoveAddrs := endpoint.NodeNotifier.GetRemoveAddrs()
		s.Assert().EqualValues(nodeCreateAddrs, epNodeCreateAddrs)
		s.Assert().EqualValues(nodeRemoveAddrs, epNodeRemoveAddrs)
	}

	if len(serviceCreateAddrs) > 0 || len(serviceRemoveAddrs) > 0 {
		s.NotNil(endpoint.ServiceChan)
	}
	if len(nodeCreateAddrs) > 0 || len(nodeRemoveAddrs) > 0 {
		s.NotNil(endpoint.NodeChan)
	}

}
