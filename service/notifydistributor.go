package service

import (
	"net/url"
)

// Notification is a node notification
type Notification struct {
	eventType EventType
	urlValues url.Values
}

// NotifyEndpoint holds Notifiers and channels to watch
type NotifyEndpoint struct {
	ServiceChan     chan Notification
	ServiceNotifier NotificationSender
	NodeChan        chan Notification
	NodeNotifier    NotificationSender
}

// NotifyDistributing takes a stream of `Notification` and
// NodeNotifiction and distributes it listeners
type NotifyDistributing interface {
	Run(serviceChan <-chan Notification, nodeChan <-chan Notification)
	HasServiceListeners() bool
	HasNodeListeners() bool
}

// NotifyDistributor distributes service and node notifications to `NotifyEnpoint`s
type NotifyDistributor struct {
	NotifyEndpoints map[string]NotifyEndpoint
}

func newNotifyDistributor(notifyEndpoints map[string]NotifyEndpoint) *NotifyDistributor {
	return &NotifyDistributor{NotifyEndpoints: notifyEndpoints}
}

// NewNotifyDistributorFromEnv creates `NotifyDistributor` from environment variables
func NewNotifyDistributorFromEnv() *NotifyDistributor {
	return nil
}

// Run starts distributor
func (n NotifyDistributor) Run(serviceChan <-chan Notification, nodeChan <-chan Notification) {

}

// HasServiceListeners when there exists service listeners
func (n NotifyDistributor) HasServiceListeners() bool {
	return false
}

// HasNodeListeners when there exists node listeners
func (n NotifyDistributor) HasNodeListeners() bool {
	return false
}
