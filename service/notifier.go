package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"../metrics"
)

// NotifyType is the type of notification to send
type NotifyType string

// NotificationSender sends notifications to listeners
type NotificationSender interface {
	Create(params string) error
	Remove(params string) error
	GetCreateAddrs() []string
	GetRemoveAddrs() []string
}

// Notifier implements `NotificationSender`
type Notifier struct {
	createAddrs       []string
	removeAddrs       []string
	notifyType        string
	retries           int
	interval          int
	createErrorMetric string
	removeErrorMetric string
	log               *log.Logger
}

// NewNotifier returns a `Notifier`
func NewNotifier(
	createAddrs, removeAddrs []string, notifyType string,
	retries int, interval int, logger *log.Logger) *Notifier {
	if createAddrs == nil {
		createAddrs = []string{}
	}
	if removeAddrs == nil {
		removeAddrs = []string{}
	}
	return &Notifier{
		createAddrs:       createAddrs,
		removeAddrs:       removeAddrs,
		notifyType:        notifyType,
		retries:           retries,
		interval:          interval,
		createErrorMetric: fmt.Sprintf("notificationSendCreate%sRequest", notifyType),
		removeErrorMetric: fmt.Sprintf("notificationSendRemove%sRequest", notifyType),
		log:               logger,
	}
}

// GetCreateAddrs returns create addresses
func (n Notifier) GetCreateAddrs() []string {
	return n.createAddrs
}

// GetRemoveAddrs returns create addresses
func (n Notifier) GetRemoveAddrs() []string {
	return n.removeAddrs
}

// Create sends create notifications to listeners
func (n Notifier) Create(params string) error {

	hasError := false
	wg := &sync.WaitGroup{}
	for _, addr := range n.createAddrs {
		wg.Add(1)
		go n.sendCreate(addr, params, wg, &hasError)
	}
	wg.Wait()

	if !hasError {
		return nil
	}
	return fmt.Errorf("At least one create %s request produced errors. Please consult logs for more details", n.notifyType)
}

// Remove sends remove notifications to listeners
func (n Notifier) Remove(params string) error {
	hasError := false
	wg := &sync.WaitGroup{}
	for _, addr := range n.removeAddrs {
		wg.Add(1)
		go n.sendRemove(addr, params, wg, &hasError)
	}
	wg.Wait()

	if !hasError {
		return nil
	}
	return fmt.Errorf("At least one remove %s request produced errors. Please consult logs for more details", n.notifyType)
}

func (n Notifier) sendCreate(addr string, params string, wg *sync.WaitGroup, hasError *bool) {
	defer wg.Done()

	urlObj, err := url.Parse(addr)
	if err != nil {
		n.log.Printf("ERROR: %v", err)
		metrics.RecordError(n.createErrorMetric)
		*hasError = true
		return
	}
	urlObj.RawQuery = params
	fullURL := urlObj.String()
	n.log.Printf("Sending %s created notification to %s", n.notifyType, fullURL)
	for i := 1; i <= n.retries; i++ {
		resp, err := http.Get(fullURL)
		if err == nil &&
			(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusConflict) {
			break
		} else if i < n.retries {
			if n.interval > 0 {
				n.log.Printf("Retrying %s created notification to %s (%d try)", n.notifyType, fullURL, i)
				time.Sleep(time.Second * time.Duration(n.interval))
			}
		} else {
			if err != nil {
				n.log.Printf("ERROR: %v", err)
				metrics.RecordError(n.createErrorMetric)
				*hasError = true
			} else if resp.StatusCode == http.StatusConflict {
				body, _ := ioutil.ReadAll(resp.Body)
				n.log.Printf("ERROR: Request %s returned status code %d\n%s", fullURL, resp.StatusCode, string(body[:]))
				metrics.RecordError(n.createErrorMetric)
				*hasError = true
			} else if resp.StatusCode != http.StatusOK {
				body, _ := ioutil.ReadAll(resp.Body)
				n.log.Printf("ERROR: Request %s returned status code %d\n%s", fullURL, resp.StatusCode, string(body[:]))
				metrics.RecordError(n.createErrorMetric)
				*hasError = true
			}
		}
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}
}

func (n Notifier) sendRemove(addr string, params string, wg *sync.WaitGroup, hasError *bool) {
	defer wg.Done()

	urlObj, err := url.Parse(addr)
	if err != nil {
		n.log.Printf("ERROR: %v", err)
		metrics.RecordError(n.removeErrorMetric)
		*hasError = true
		return
	}
	urlObj.RawQuery = params
	fullURL := urlObj.String()
	n.log.Printf("Sending %s removed notification to %s", n.notifyType, fullURL)
	for i := 1; i <= n.retries; i++ {
		resp, err := http.Get(fullURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		} else if i < n.retries {
			if n.interval > 0 {
				n.log.Printf("Retrying %s removed notification to %s (%d try)", n.notifyType, fullURL, i)
				time.Sleep(time.Second * time.Duration(n.interval))
			}
		} else {
			if err != nil {
				n.log.Printf("ERROR: %v", err)
				metrics.RecordError(n.removeErrorMetric)
				*hasError = true
			} else if resp.StatusCode != http.StatusOK {
				body, _ := ioutil.ReadAll(resp.Body)
				n.log.Printf("ERROR: Request %s returned status code %d\n%s", fullURL, resp.StatusCode, string(body[:]))
				metrics.RecordError(n.removeErrorMetric)
				*hasError = true
			}
		}
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}
}
