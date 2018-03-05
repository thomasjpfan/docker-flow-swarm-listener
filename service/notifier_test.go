package service

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
)

type NotifierTestSuite struct {
	suite.Suite
	Logger       *log.Logger
	LogBytes     *bytes.Buffer
	CreateValues url.Values
	RemoveValues url.Values
}

func TestNotifierUnitTestSuite(t *testing.T) {
	suite.Run(t, new(NotifierTestSuite))
}

func (s *NotifierTestSuite) SetupSuite() {
	s.LogBytes = new(bytes.Buffer)
	s.Logger = log.New(s.LogBytes, "", 0)

	cParams := url.Values{}
	cParams.Add("replicas", "3")
	cParams.Add("serviceName", "hello")
	s.CreateValues = cParams

	rParams := url.Values{}
	rParams.Add("serviceName", "hello")
}

func (s *NotifierTestSuite) TearDownTest() {
	s.LogBytes.Reset()
}

// Create

func (s *NotifierTestSuite) Test_Create_SendsRequests() {

	var query1, query2 url.Values
	httpSrv := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			switch r.URL.Path {
			case "/v1/docker-flow-proxy/reconfigure":
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				query1 = r.URL.Query()
			case "/something/else":
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				query2 = r.URL.Query()
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}
	}))
	defer httpSrv.Close()

	url1 := fmt.Sprintf("%s/v1/docker-flow-proxy/reconfigure", httpSrv.URL)
	url2 := fmt.Sprintf("%s/something/else", httpSrv.URL)

	n := NewNotifier([]string{url1, url2}, []string{}, NotifyServiceType, 5, 1, s.Logger)
	err := n.Create(s.CreateValues)
	s.Require().NoError(err)

	s.EqualURLValues(s.CreateValues, query1)
	s.EqualURLValues(s.CreateValues, query2)

	urlObj1, err := url.Parse(url1)
	s.Require().NoError(err)
	urlObj2, err := url.Parse(url2)
	s.Require().NoError(err)

	urlObj1.RawQuery = s.CreateValues.Encode()
	urlObj2.RawQuery = s.CreateValues.Encode()

	logMessages := s.LogBytes.String()
	s.Contains(logMessages, fmt.Sprintf("Sending service created notification to %s", urlObj1.String()))
	s.Contains(logMessages, fmt.Sprintf("Sending service created notification to %s", urlObj2.String()))
}

func (s *NotifierTestSuite) Test_Create_ReturnsError_WhenUrlCannotBeParsed() {
	n := NewNotifier([]string{"%%%"}, []string{}, NotifyServiceType, 5, 1, s.Logger)
	err := n.Create(s.CreateValues)
	s.Error(err)
}

func (s *NotifierTestSuite) Test_Create_ReturnsError_WhenHttpStatusIsNot200() {

	httpSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	n := NewNotifier(
		[]string{httpSrv.URL}, []string{}, NotifyNodeType, 1, 0, s.Logger)
	err := n.Create(s.CreateValues)
	s.Require().NoError(err)
}

func (s *NotifierTestSuite) Test_Create_DoesNotReturnError_WhenHttpStatusIs409() {

	httpSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
	}))

	n := NewNotifier(
		[]string{httpSrv.URL}, []string{}, NotifyNodeType, 1, 0, s.Logger)
	err := n.Create(s.CreateValues)

	s.NoError(err)
}

func (s *NotifierTestSuite) Test_Create_ReturnsError_WhenHttpRequestReturnsError() {
	n := NewNotifier(
		[]string{}, []string{"this-does-not-exist"}, NotifyNodeType, 2, 1, s.Logger)

	err := n.Create(s.CreateValues)
	s.Error(err)
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
		[]string{httpSrv.URL}, []string{}, NotifyServiceType, 2, 1, s.Logger)
	err := n.Create(s.CreateValues)
	s.Require().NoError(err)

	s.Equal(2, attempt)

}

// Remove

func (s *NotifierTestSuite) Test_Remove_SendsRequests() {
	var query1, query2 url.Values

	httpSrv := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			switch r.URL.Path {
			case "/v1/docker-flow-proxy/remove":
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				query1 = r.URL.Query()
			case "/something/else":
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				query2 = r.URL.Query()
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}
	}))
	defer httpSrv.Close()

	url1 := fmt.Sprintf("%s/v1/docker-flow-proxy/remove", httpSrv.URL)
	url2 := fmt.Sprintf("%s/something/else", httpSrv.URL)

	n := NewNotifier([]string{url1, url2}, []string{}, NotifyServiceType, 5, 1, s.Logger)
	err := n.Remove(s.RemoveValues)
	s.Require().NoError(err)

	s.EqualURLValues(s.RemoveValues, query1)
	s.EqualURLValues(s.RemoveValues, query2)

	urlObj1, err := url.Parse(url1)
	s.Require().NoError(err)
	urlObj2, err := url.Parse(url2)
	s.Require().NoError(err)

	urlObj1.RawQuery = s.RemoveValues.Encode()
	urlObj2.RawQuery = s.RemoveValues.Encode()

	logMessages := s.LogBytes.String()
	s.Contains(logMessages, fmt.Sprintf("Sending node removed notification to %s", urlObj1.String()))
	s.Contains(logMessages, fmt.Sprintf("Sending node removed notification to %s", urlObj2.String()))
}

func (s *NotifierTestSuite) Test_Remove_ReturnsError_WhenUrlCannotBeParsed() {
	n := NewNotifier([]string{}, []string{"%%%"}, NotifyNodeType, 5, 1, s.Logger)
	err := n.Remove(s.RemoveValues)
	s.Error(err)
}

func (s *NotifierTestSuite) Test_Remove_ReturnsError_WhenHttpStatusIsNot200() {

	httpSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	n := NewNotifier(
		[]string{}, []string{httpSrv.URL}, NotifyServiceType, 1, 0, s.Logger)
	err := n.Remove(s.CreateValues)
	s.NoError(err)
}

func (s *NotifierTestSuite) Test_Remove_ReturnsError_WhenHttpRequestReturnsError() {
	n := NewNotifier(
		[]string{"this-does-not-exist"}, []string{}, NotifyServiceType, 2, 1, s.Logger)
	err := n.Remove(s.RemoveValues)
	s.Error(err)
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
		[]string{}, []string{httpSrv.URL}, NotifyNodeType, 2, 1, s.Logger)
	err := n.Remove(s.CreateValues)
	s.Require().NoError(err)

	s.Equal(2, attempt)
}

func (s *NotifierTestSuite) EqualURLValues(expected, actual url.Values) {
	for k := range expected {
		expV, expA := expected[k], actual[k]
		s.Len(expV, 1)
		s.Len(expA, 1)
		s.Equal(expected.Get(k), actual.Get(k))
	}
}
