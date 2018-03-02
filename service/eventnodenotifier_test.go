package service

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types/swarm"
	"github.com/stretchr/testify/suite"
)

type EventNodeNotifierTestSuite struct {
	suite.Suite
}

func TestEventNodeNotifierUnitTestSuite(t *testing.T) {
	logPrintfOrig := logPrintf
	defer func() {
		logPrintf = logPrintfOrig
	}()
	logPrintf = func(format string, v ...interface{}) {}

	s := new(EventNodeNotifierTestSuite)
	suite.Run(t, s)
}

func (s *EventNodeNotifierTestSuite) Test_NewNotificationFromEnv_ParseENV() {
	defer func() {
		os.Unsetenv("DF_NOTIFY_CREATE_NODE_URL")
		os.Unsetenv("DF_NOTIFY_UPDATE_NODE_URL")
	}()
	os.Setenv("DF_NOTIFY_CREATE_NODE_URL", "create_url1,create_url2")
	os.Setenv("DF_NOTIFY_UPDATE_NODE_URL", "update_url1")

	n := NewEventNodeNotifierFromEnv()
	s.Require().NotNil(n)

	s.Require().Len(n.CreateAddrs, 2)
	s.Equal("create_url1", n.CreateAddrs[0])
	s.Equal("create_url2", n.CreateAddrs[1])

	s.Require().Len(n.UpdateAddrs, 1)
	s.Equal("update_url1", n.UpdateAddrs[0])

	s.Len(n.RemoveAddrs, 0)

	s.True(n.HasListeners())
}

func (s *EventNodeNotifierTestSuite) Test_NewNotification_NoListeners() {
	n := NewEventNodeNotifierFromEnv()
	s.False(n.HasListeners())
}

func (s *EventNodeNotifierTestSuite) Test_CreateNodes_SendRequests() {
	queryChan1 := make(chan url.Values, 1)
	queryChan2 := make(chan url.Values, 1)

	httpSrv := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" &&
			r.URL.Path == "/v1/docker-flow-proxy/reconfigure" {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")

			if r.URL.Query().Get("hostname") == "managerHostname" {
				queryChan1 <- r.URL.Query()
			} else if r.URL.Query().Get("hostname") == "workerNodeID" {
				queryChan2 <- r.URL.Query()
			}
		}
	}))

	defer httpSrv.Close()

	nodeHN1 := "managerHostname"
	label1 := map[string]string{
		"com.df.cows":   "grass",
		"com.df2.bears": "fight",
	}
	node1 := getNode(nodeHN1, "managerNodeID", swarm.NodeRoleManager, label1)

	nodeHN2 := "workerHostname"
	label2 := map[string]string{
		"com.df.birds":  "fly",
		"com.df.zebras": "run",
	}
	node2 := getNode(nodeHN2, "workerNodeID", swarm.NodeRoleWorker, label2)

	url1 := fmt.Sprintf("%s/v1/docker-flow-proxy/reconfigure", httpSrv.URL)
	n := newEventNodeNotifier(
		[]string{url1}, []string{}, []string{})

	err := n.NotifyCreateNodes([]swarm.Node{node1, node2}, 50, 5)
	s.Require().NoError(err)

	timeoutChan := time.NewTimer(5 * time.Second).C

	var query1 url.Values
	var query2 url.Values

	for {
		if queryChan1 == nil && queryChan2 == nil {
			break
		}
		select {
		case q := <-queryChan1:
			query1 = q
			queryChan1 = nil
		case q := <-queryChan2:
			query2 = q
			queryChan2 = nil
		case <-timeoutChan:
			s.Fail("Timeout")
			return
		}
	}

	params1 := GetNodeParameters(node1)
	params2 := GetNodeParameters(node2)

	s.EqualURLValues(params1, query1)
	s.EqualURLValues(params2, query2)
}

func (s *EventNodeNotifierTestSuite) Test_CreateNode_SendRequests_TwoURLs() {

	queryChan1 := make(chan url.Values, 1)
	queryChan2 := make(chan url.Values, 1)
	httpSrv := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			switch r.URL.Path {
			case "/v1/docker-flow-proxy/reconfigure":
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				queryChan1 <- r.URL.Query()
			case "/something/else":
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				queryChan2 <- r.URL.Query()
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}
	}))
	defer httpSrv.Close()

	url1 := fmt.Sprintf("%s/v1/docker-flow-proxy/reconfigure", httpSrv.URL)
	url2 := fmt.Sprintf("%s/something/else", httpSrv.URL)

	// Setup mock
	labels := map[string]string{
		"com.df.hello":     "world",
		"com.df2.areducks": "real",
	}
	hostname := "managerHostname"
	managerNode := getNode(hostname, "managerNodeID", swarm.NodeRoleManager, labels)

	n := newEventNodeNotifier(
		[]string{url1, url2}, []string{}, []string{})
	err := n.NotifyCreateNode(managerNode, 50, 5)
	s.Require().NoError(err)

	timeoutChan := time.NewTimer(5 * time.Second).C

	var query1 url.Values
	var query2 url.Values

	for {
		if queryChan1 == nil && queryChan2 == nil {
			break
		}
		select {
		case q := <-queryChan1:
			query1 = q
			queryChan1 = nil
		case q := <-queryChan2:
			query2 = q
			queryChan2 = nil
		case <-timeoutChan:
			s.Fail("Timeout")
			return
		}
	}

	params := GetNodeParameters(managerNode)

	s.EqualURLValues(params, query1)
	s.EqualURLValues(params, query2)
}

func (s *EventNodeNotifierTestSuite) Test_CreateNode_RetriesRequests() {

	count := 0
	done := make(chan struct{})
	httpSrv := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		count++
		if count == 3 {
			done <- struct{}{}
		}
	}))

	node := getNode(
		"hostname", "node123", swarm.NodeRoleManager,
		map[string]string{})
	n := newEventNodeNotifier([]string{httpSrv.URL}, []string{}, []string{})

	err := n.NotifyCreateNode(node, 3, 1)
	s.Require().NoError(err)

	timerChan := time.NewTicker(5 * time.Second).C

	for {
		select {
		case <-done:
			s.True(true, "Retried three times")
		case <-timerChan:
			s.Fail("Timeout")
			return
		}
	}

}

func (s *EventNodeNotifierTestSuite) Test_UpdateNode_SendRequests_OneURL() {

	queryChan := make(chan url.Values, 1)
	httpSrv := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" &&
			r.URL.Path == "/v1/docker-flow-proxy/reconfigure" {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			queryChan <- r.URL.Query()
		}
	}))
	defer httpSrv.Close()

	labels := map[string]string{
		"com.df.hello": "world",
	}
	hostname := "workerHostname"
	workerNode := getNode(hostname, "workerNodeID", swarm.NodeRoleWorker, labels)

	url1 := fmt.Sprintf("%s/v1/docker-flow-proxy/reconfigure", httpSrv.URL)
	n := newEventNodeNotifier(
		[]string{}, []string{url1}, []string{})
	err := n.NotifyUpdateNode(workerNode, 50, 5)
	s.Require().NoError(err)

	timeoutChan := time.NewTimer(5 * time.Second).C

	var query url.Values
	for {
		if queryChan == nil {
			break
		}
		select {
		case q := <-queryChan:
			query = q
			queryChan = nil
		case <-timeoutChan:
			s.Fail("Timeout")
			return
		}
	}

	params := GetNodeParameters(workerNode)

	s.EqualURLValues(params, query)
}

func (s *EventNodeNotifierTestSuite) Test_RemoveNode_LogsError_WhenHTTPStatusIsNot200() {

	// Mock logger
	logPrintfOrig := logPrintf
	msg := ""
	defer func() { logPrintf = logPrintfOrig }()
	logPrintf = func(format string, v ...interface{}) {
		msg = format
	}

	called := make(chan struct{})
	httpSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		called <- struct{}{}
	}))

	hostname := "workerHostname"
	workerNode := getNode(
		hostname, "workerNodeID", swarm.NodeRoleWorker, map[string]string{})

	n := newEventNodeNotifier(
		[]string{}, []string{}, []string{httpSrv.URL})
	err := n.NotifyRemoveNode(workerNode, 50, 5)
	s.Require().NoError(err)

	timeoutChan := time.NewTimer(5 * time.Second).C

	for {
		if called == nil {
			break
		}
		select {
		case <-called:
			called = nil
		case <-timeoutChan:
			s.Fail("Timeout")
			return
		}
	}

	s.True(strings.HasPrefix(msg, "ERROR"))
}
func (s *EventNodeNotifierTestSuite) EqualURLValues(expected, actual url.Values) {
	for k := range expected {
		expV, expA := expected[k], actual[k]
		s.Len(expV, 1)
		s.Len(expA, 1)
		s.Equal(expected.Get(k), actual.Get(k))
	}
}
