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

type NodeEventNotifierTestSuite struct {
	suite.Suite
}

func TestNodeEventNotifierUnitTestSuite(t *testing.T) {
	logPrintfOrig := logPrintf
	defer func() {
		logPrintf = logPrintfOrig
	}()
	logPrintf = func(format string, v ...interface{}) {}

	s := new(NodeEventNotifierTestSuite)
	suite.Run(t, s)
}

func (s *NodeEventNotifierTestSuite) Test_NewNotificationFromEnv_ParseENV() {
	defer func() {
		os.Unsetenv("DF_NOTIFY_CREATE_NODE_URL")
		os.Unsetenv("DF_NOTIFY_UPDATE_NODE_URL")
	}()
	os.Setenv("DF_NOTIFY_CREATE_NODE_URL", "create_url1,create_url2")
	os.Setenv("DF_NOTIFY_UPDATE_NODE_URL", "update_url1")

	n := NewNodeEventNotifierFromEnv()
	s.Require().NotNil(n)

	s.Require().Len(n.CreateAddrs, 2)
	s.Equal("create_url1", n.CreateAddrs[0])
	s.Equal("create_url2", n.CreateAddrs[1])

	s.Require().Len(n.UpdateAddrs, 1)
	s.Equal("update_url1", n.UpdateAddrs[0])

	s.Len(n.RemoveAddrs, 0)
}

func (s *NodeEventNotifierTestSuite) Test_CreateNodes_SendRequests() {
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
	node1 := s.getNode(nodeHN1, "managerNodeID", swarm.NodeRoleManager, label1)

	nodeHN2 := "workerHostname"
	label2 := map[string]string{
		"com.df.birds":  "fly",
		"com.df.zebras": "run",
	}
	node2 := s.getNode(nodeHN2, "workerNodeID", swarm.NodeRoleWorker, label2)

	url1 := fmt.Sprintf("%s/v1/docker-flow-proxy/reconfigure", httpSrv.URL)
	n := newNodeEventNotifier(
		[]string{url1}, []string{}, []string{})
	n.NotifyCreateNodes([]swarm.Node{node1, node2}, 50, 5)

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

	s.Equal(nodeHN1, query1.Get("hostname"))
	s.Equal("grass", query1.Get("cows"))
	s.Equal("", query1.Get("bears"))

	s.Equal(nodeHN2, query2.Get("hostname"))
	s.Equal("fly", query1.Get("birds"))
	s.Equal("run", query1.Get("zebras"))

}

func (s *NodeEventNotifierTestSuite) Test_CreateNode_SendRequests_TwoURLs() {

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
	managerNode := s.getNode(hostname, "managerNodeID", swarm.NodeRoleManager, labels)

	n := newNodeEventNotifier(
		[]string{url1, url2}, []string{}, []string{})
	n.NotifyCreateNode(managerNode, 50, 5)

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
	s.Equal(hostname, query1.Get("hostname"))
	s.Equal("true", query1.Get("manager"))
	s.Equal("world", query1.Get("hello"))
	s.Equal("", query1.Get("areducks"))

	s.Equal(hostname, query2.Get("hostname"))
	s.Equal("true", query2.Get("manager"))
	s.Equal("world", query2.Get("hello"))
	s.Equal("", query2.Get("areducks"))
}

func (s *NodeEventNotifierTestSuite) Test_UpdateNode_SendRequests_OneURL() {

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
	workerNode := s.getNode(hostname, "workerNodeID", swarm.NodeRoleWorker, labels)

	url1 := fmt.Sprintf("%s/v1/docker-flow-proxy/reconfigure", httpSrv.URL)
	n := newNodeEventNotifier(
		[]string{}, []string{url1}, []string{})
	n.NotifyUpdateNode(workerNode, 50, 5)

	timeoutChan := time.NewTimer(5 * time.Second).C

	for {
		if queryChan == nil {
			break
		}
		select {
		case q := <-queryChan:
			s.Equal(hostname, q.Get("hostname"))
			s.Equal("false", q.Get("manager"))
			s.Equal("world", q.Get("hello"))
			queryChan = nil
		case <-timeoutChan:
			s.Fail("Timeout")
			return
		}
	}

}

func (s *NodeEventNotifierTestSuite) Test_RemoveNode_LogsError_WhenHTTPStatusIsNot200() {

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
	workerNode := s.getNode(
		hostname, "workerNodeID", swarm.NodeRoleWorker, map[string]string{})

	n := newNodeEventNotifier(
		[]string{}, []string{}, []string{httpSrv.URL})
	n.NotifyRemoveNode(workerNode, 50, 5)

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

func (s *NodeEventNotifierTestSuite) getNode(
	hostname string, nodeID string,
	role swarm.NodeRole, labels map[string]string) swarm.Node {

	annotations := swarm.Annotations{
		Labels: labels,
	}
	nodeSpec := swarm.NodeSpec{
		Annotations: annotations,
		Role:        role,
	}
	nodeDescription := swarm.NodeDescription{
		Hostname: hostname,
	}
	return swarm.Node{
		ID:          nodeID,
		Description: nodeDescription,
		Spec:        nodeSpec,
	}
}
