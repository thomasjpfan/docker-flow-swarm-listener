package service

import (
	"bytes"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type NotifyDistributorTestSuite struct {
	suite.Suite
	log      *log.Logger
	logBytes *bytes.Buffer
}

func TestNotifyDistributorUnitTestSuite(t *testing.T) {
	suite.Run(t, new(NotifyDistributorTestSuite))
}

func (s *NotifyDistributorTestSuite) SetupSuite() {
	s.logBytes = new(bytes.Buffer)
	s.log = log.New(s.logBytes, "", 0)
}
func (s *NotifyDistributorTestSuite) SetupTest() {
	s.logBytes.Reset()
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

	s.True(notifyD.HasServiceListeners())
	s.True(notifyD.HasNodeListeners())

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

	s.True(notifyD.HasServiceListeners())
	s.True(notifyD.HasNodeListeners())
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

	s.True(notifyD.HasServiceListeners())
	s.False(notifyD.HasNodeListeners())
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

	s.False(notifyD.HasServiceListeners())
	s.True(notifyD.HasNodeListeners())
}

func (s *NotifyDistributorTestSuite) Test_RunDistributesNotificationsToEndpoints_Servies() {
	mock1Create := make(chan struct{})
	mock1Remove := make(chan struct{})
	mock2Create := make(chan struct{})
	mock2Remove := make(chan struct{})

	serviceNotifyMock1 := notificationSenderMock{}
	serviceNotifyMock1.On("Create", "hello=world").Return(nil).Run(func(args mock.Arguments) {
		mock1Create <- struct{}{}
	})
	serviceNotifyMock1.On("Remove", "hello=world2").Return(nil).Run(func(args mock.Arguments) {
		mock1Remove <- struct{}{}
	})

	serviceNotifyMock2 := notificationSenderMock{}
	serviceNotifyMock2.On("Create", "hello=world").Return(nil).Run(func(args mock.Arguments) {
		mock2Create <- struct{}{}
	})
	serviceNotifyMock2.On("Remove", "hello=world2").Return(nil).Run(func(args mock.Arguments) {
		mock2Remove <- struct{}{}
	})

	endpoints := map[string]NotifyEndpoint{
		"host1": NotifyEndpoint{
			ServiceChan:     make(chan Notification),
			ServiceNotifier: &serviceNotifyMock1,
			NodeChan:        nil,
			NodeNotifier:    nil,
		},
		"host2": NotifyEndpoint{
			ServiceChan:     make(chan Notification),
			ServiceNotifier: &serviceNotifyMock2,
			NodeChan:        nil,
			NodeNotifier:    nil,
		},
	}

	notifyD := newNotifyDistributor(endpoints, 1, s.log)
	serviceChan := make(chan Notification)

	notifyD.Run(serviceChan, nil)

	go func() {
		serviceChan <- Notification{
			EventType:  EventTypeCreate,
			Parameters: "hello=world",
		}
	}()
	go func() {
		serviceChan <- Notification{
			EventType:  EventTypeRemove,
			Parameters: "hello=world2",
		}
	}()

	timer := time.NewTimer(time.Second * 5).C

	for {
		if mock1Create == nil && mock1Remove == nil &&
			mock2Create == nil && mock2Remove == nil {
			break
		}
		select {
		case <-mock1Create:
			mock1Create = nil
		case <-mock1Remove:
			mock1Remove = nil
		case <-mock2Create:
			mock2Create = nil
		case <-mock2Remove:
			mock2Remove = nil
		case <-timer:
			s.Fail("Timeout")
			return
		}
	}

	serviceNotifyMock1.AssertExpectations(s.T())
	serviceNotifyMock2.AssertExpectations(s.T())
}

func (s *NotifyDistributorTestSuite) Test_RunDistributesNotificationsToEndpoints_Nodes1() {
	mock1Create := make(chan struct{})
	mock1Remove := make(chan struct{})
	mock2Create := make(chan struct{})
	mock2Remove := make(chan struct{})

	nodesNotifyMock1 := notificationSenderMock{}
	nodesNotifyMock1.On("Create", "hello=world").Return(nil).Run(func(args mock.Arguments) {
		mock1Create <- struct{}{}
	})
	nodesNotifyMock1.On("Remove", "hello=world2").Return(nil).Run(func(args mock.Arguments) {
		mock1Remove <- struct{}{}
	})

	nodesNotifyMock2 := notificationSenderMock{}
	nodesNotifyMock2.On("Create", "hello=world").Return(nil).Run(func(args mock.Arguments) {
		mock2Create <- struct{}{}
	})
	nodesNotifyMock2.On("Remove", "hello=world2").Return(nil).Run(func(args mock.Arguments) {
		mock2Remove <- struct{}{}
	})

	endpoints := map[string]NotifyEndpoint{
		"host1": NotifyEndpoint{
			ServiceChan:     nil,
			ServiceNotifier: nil,
			NodeChan:        make(chan Notification),
			NodeNotifier:    &nodesNotifyMock1,
		},
		"host2": NotifyEndpoint{
			ServiceChan:     nil,
			ServiceNotifier: nil,
			NodeChan:        make(chan Notification),
			NodeNotifier:    &nodesNotifyMock2,
		},
	}

	notifyD := newNotifyDistributor(endpoints, 1, s.log)
	nodeChan := make(chan Notification)

	notifyD.Run(nil, nodeChan)

	go func() {
		nodeChan <- Notification{
			EventType:  EventTypeCreate,
			Parameters: "hello=world",
		}
	}()
	go func() {
		nodeChan <- Notification{
			EventType:  EventTypeRemove,
			Parameters: "hello=world2",
		}
	}()

	timer := time.NewTimer(time.Second * 5).C

	for {
		if mock1Create == nil && mock1Remove == nil &&
			mock2Create == nil && mock2Remove == nil {
			break
		}
		select {
		case <-mock1Create:
			mock1Create = nil
		case <-mock1Remove:
			mock1Remove = nil
		case <-mock2Create:
			mock2Create = nil
		case <-mock2Remove:
			mock2Remove = nil
		case <-timer:
			s.Fail("Timeout")
			return
		}
	}

	nodesNotifyMock1.AssertExpectations(s.T())
	nodesNotifyMock2.AssertExpectations(s.T())
}

func (s *NotifyDistributorTestSuite) Test_RunDistributesNotificationsToEndpoints_Nodes_OnError_PutsbackNotificationOnQueue() {
	mockCreate := make(chan struct{})
	numOfCalls := 0
	err := fmt.Errorf("Nope")

	nodesNotifyMock := notificationSenderMock{}
	nodesNotifyMock.On("Create", "hello=world2").Return(err).Run(func(args mock.Arguments) {
		numOfCalls++
		if numOfCalls == 2 {
			mockCreate <- struct{}{}
		}
	})
	nodesNotifyMock.On("GetCreateAddr").Return("http://host1/create")

	endpoints := map[string]NotifyEndpoint{
		"host1": NotifyEndpoint{
			ServiceChan:     nil,
			ServiceNotifier: nil,
			NodeChan:        make(chan Notification),
			NodeNotifier:    &nodesNotifyMock,
		},
	}
	notifyD := newNotifyDistributor(endpoints, 1, s.log)
	nodeChan := make(chan Notification)

	notifyD.Run(nil, nodeChan)

	go func() {
		nodeChan <- Notification{
			EventType:  EventTypeCreate,
			Parameters: "hello=world2",
		}
	}()

	timer := time.NewTimer(time.Second * 5).C

	for {
		if mockCreate == nil {
			break
		}
		select {
		case <-mockCreate:
			mockCreate = nil
		case <-timer:
			s.Fail("Timeout")
			return
		}
	}

	s.Equal(2, numOfCalls)
}

func (s *NotifyDistributorTestSuite) Test_RunDistributesNotificationsToEndpoints_Service_OnError_PutsbackNotificationOnQueue() {
	mockRemove := make(chan struct{})
	numOfCalls := 0
	err := fmt.Errorf("Nope")

	serviceNotifyMock := notificationSenderMock{}
	serviceNotifyMock.On("Remove", "hello=world2").Return(err).Run(func(args mock.Arguments) {
		numOfCalls++
		if numOfCalls == 2 {
			mockRemove <- struct{}{}
		}
	})
	serviceNotifyMock.On("GetRemoveAddr").Return("http://host1/remove")

	endpoints := map[string]NotifyEndpoint{
		"host1": NotifyEndpoint{
			ServiceChan:     make(chan Notification),
			ServiceNotifier: &serviceNotifyMock,
			NodeChan:        nil,
			NodeNotifier:    nil,
		},
	}
	notifyD := newNotifyDistributor(endpoints, 1, s.log)
	serviceChan := make(chan Notification)

	notifyD.Run(serviceChan, nil)

	go func() {
		serviceChan <- Notification{
			EventType:  EventTypeRemove,
			Parameters: "hello=world2",
		}
	}()

	timer := time.NewTimer(time.Second * 5).C

	for {
		if mockRemove == nil {
			break
		}
		select {
		case <-mockRemove:
			mockRemove = nil
		case <-timer:
			s.Fail("Timeout")
			return
		}
	}

	s.Equal(2, numOfCalls)
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
