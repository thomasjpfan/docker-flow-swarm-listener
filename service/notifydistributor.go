package service

import (
	"net/url"
)

// NotificationValue is a node notification
type NotificationValue struct {
	eventType EventType
	urlValues url.Values
}

// NotifyEndpoint holds Notifiers and channels to watch
type NotifyEndpoint struct {
	Host            string
	ServiceChan     chan NotificationValue
	ServiceNotifier *Notifier
	NodeChan        chan NotificationValue
	NodeNotifier    *Notifier
}

// NotifyDistributing takes a stream of `NotificationValue` and
// NodeNotifiction and distributes it listeners
type NotifyDistributing interface {
	Run(serviceChan <-chan NotificationValue, nodeChan <-chan NotificationValue)
}
