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
	Host            string
	ServiceChan     chan Notification
	ServiceNotifier *Notifier
	NodeChan        chan Notification
	NodeNotifier    *Notifier
}

// NotifyDistributing takes a stream of `Notification` and
// NodeNotifiction and distributes it listeners
type NotifyDistributing interface {
	Run(serviceChan <-chan Notification, nodeChan <-chan Notification)
}

// NotifyDistributor distributes service and node notifications to `NotifyEnpoint`s
type NotifyDistributor struct {
	NotifyEndpoints []NotifyEndpoint
}

func newNotifyDistributor(notifyEndpoint []NotifyEndpoint) *NotifyDistributor {
	return &NotifyDistributor{NotifyEndpoints: notifyEndpoint}
}

// NewNotifyDistributorFromEnv creates `NotifyDistributor` from environment variables
func NewNotifyDistributorFromEnv() *NotifyDistributor {
	return nil
}
