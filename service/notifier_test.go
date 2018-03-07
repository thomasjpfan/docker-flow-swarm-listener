package service

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type NotifierTestSuite struct {
	suite.Suite
	Logger   *log.Logger
	LogBytes *bytes.Buffer
	Params   string
}

func TestNotifierUnitTestSuite(t *testing.T) {
	suite.Run(t, new(NotifierTestSuite))
}

func (s *NotifierTestSuite) SetupSuite() {
	s.LogBytes = new(bytes.Buffer)
	s.Logger = log.New(s.LogBytes, "", 0)

	cParams := url.Values{}
	cParams.Add("serviceName", "hello")
	s.Params = cParams.Encode()

}

func (s *NotifierTestSuite) TearDownTest() {
	s.LogBytes.Reset()
}

// Create

func (s *NotifierTestSuite) Test_Create_SendsRequests() {

	var query1 string
	httpSrv := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			switch r.URL.Path {
			case "/v1/docker-flow-proxy/reconfigure":
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				query1 = r.URL.Query().Encode()
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}
	}))
	defer httpSrv.Close()

	url1 := fmt.Sprintf("%s/v1/docker-flow-proxy/reconfigure", httpSrv.URL)

	n := NewNotifier(url1, "", "service", 5, 1, s.Logger)
	s.Equal(url1, n.GetCreateAddr())
	err := n.Create(s.Params)
	s.Require().NoError(err)

	s.Equal(s.Params, query1)

	urlObj1, err := url.Parse(url1)
	s.Require().NoError(err)

	urlObj1.RawQuery = s.Params

	logMsgs := s.LogBytes.String()
	s.Contains(logMsgs, fmt.Sprintf("Sending service created notification to %s", urlObj1.String()))
}

func (s *NotifierTestSuite) Test_Create_ReturnsAndLogsError_WhenUrlCannotBeParsed() {
	n := NewNotifier("%%%", "", "service", 5, 1, s.Logger)
	err := n.Create(s.Params)
	s.Error(err)

	logMsgs := s.LogBytes.String()
	s.True(strings.HasPrefix(logMsgs, "ERROR: "))
}

func (s *NotifierTestSuite) Test_Create_ReturnsAndLogsError_WhenHttpStatusIsNot200() {

	httpSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	n := NewNotifier(
		httpSrv.URL, "", "node", 1, 0, s.Logger)
	err := n.Create(s.Params)
	s.Error(err)

	logMsgs := s.LogBytes.String()
	s.Contains(logMsgs, "ERROR: ")
}

func (s *NotifierTestSuite) Test_Create_ReturnsNoError_WhenHttpStatusIs409() {

	httpSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
	}))

	n := NewNotifier(
		httpSrv.URL, "", "node", 1, 0, s.Logger)
	err := n.Create(s.Params)
	s.Require().NoError(err)
}

func (s *NotifierTestSuite) Test_Create_ReturnsAndLogsError_WhenHttpRequestErrors() {
	n := NewNotifier(
		"this-does-not-exist", "", "node", 2, 1, s.Logger)

	err := n.Create(s.Params)
	s.Require().Error(err)

	logMsgs := s.LogBytes.String()
	s.Contains(logMsgs, "ERROR: ")
}

func (s *NotifierTestSuite) Test_Create_RetriesRequests() {
	attempt := 0
	httpSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if attempt < 1 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
		}
		attempt++
	}))

	n := NewNotifier(
		httpSrv.URL, "", "service", 2, 1, s.Logger)
	n.Create(s.Params)

	s.Equal(2, attempt)

	logMsgs := s.LogBytes.String()
	expMsg := fmt.Sprintf("Retrying service created notification to %s", httpSrv.URL)
	s.Contains(logMsgs, expMsg)
}

// Remove

func (s *NotifierTestSuite) Test_Remove_SendsRequests() {
	var query1 string

	httpSrv := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			switch r.URL.Path {
			case "/v1/docker-flow-proxy/remove":
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				query1 = r.URL.Query().Encode()
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}
	}))
	defer httpSrv.Close()

	url1 := fmt.Sprintf("%s/v1/docker-flow-proxy/remove", httpSrv.URL)

	n := NewNotifier("", url1, "node", 5, 1, s.Logger)
	s.Equal(url1, n.GetRemoveAddr())
	err := n.Remove(s.Params)
	s.Require().NoError(err)

	s.Equal(s.Params, query1)

	urlObj1, err := url.Parse(url1)
	s.Require().NoError(err)

	urlObj1.RawQuery = s.Params

	logMsgs := s.LogBytes.String()
	s.Contains(logMsgs, fmt.Sprintf("Sending node removed notification to %s", urlObj1.String()))
}

func (s *NotifierTestSuite) Test_Remove_ReturnsAndLogsError_WhenUrlCannotBeParsed() {
	n := NewNotifier("", "%%%", "node", 5, 1, s.Logger)
	err := n.Remove(s.Params)
	s.Error(err)

	logMsgs := s.LogBytes.String()
	s.Contains(logMsgs, "ERROR: ")
}

func (s *NotifierTestSuite) Test_Remove_ReturnsAndLogsError_WhenHttpStatusIsNot200() {

	httpSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	n := NewNotifier(
		"", httpSrv.URL, "service", 1, 0, s.Logger)
	err := n.Remove(s.Params)
	s.Error(err)

	logMsgs := s.LogBytes.String()
	s.Contains(logMsgs, "ERROR: ")
}

func (s *NotifierTestSuite) Test_Remove_ReturnsAndLogsError_WhenHttpRequestReturnsError() {
	n := NewNotifier(
		"", "this-does-not-exist", "service", 2, 1, s.Logger)
	err := n.Remove(s.Params)
	s.Error(err)

	logMsgs := s.LogBytes.String()
	s.Contains(logMsgs, "ERROR: ")
}

func (s *NotifierTestSuite) Test_Remove_RetriesRequests() {
	attempt := 0
	httpSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if attempt < 1 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
		}
		attempt++
	}))

	n := NewNotifier(
		"", httpSrv.URL, "node", 2, 1, s.Logger)
	err := n.Remove(s.Params)
	s.Require().NoError(err)

	s.Equal(2, attempt)

	logMsgs := s.LogBytes.String()
	expMsg := fmt.Sprintf("Retrying node removed notification to %s", httpSrv.URL)
	s.Contains(logMsgs, expMsg)
}

func (s *NotifierTestSuite) EqualURLValues(expected, actual url.Values) {
	for k := range expected {
		expV, expA := expected[k], actual[k]
		s.Len(expV, 1)
		s.Len(expA, 1)
		s.Equal(expected.Get(k), actual.Get(k))
	}
}
