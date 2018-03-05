package service

import (
	"log"
	"net/url"
)

// NotifyType is the type of notification to send
type NotifyType string

const (
	// NotifyServiceType is a node notification
	NotifyServiceType NotifyType = "service"
	// NotifyNodeType is a node notification
	NotifyNodeType NotifyType = "node"
)

// NotificationSender sends notifications to listeners
type NotificationSender interface {
	Create(urlValues url.Values) error
	Remove(urlValues url.Values) error
}

// Notifier implements `NotificationSender`
type Notifier struct {
	createAddrs []string
	removeAddrs []string
	notifyType  NotifyType
	retries     int
	internval   int
	log         *log.Logger
}

// NewNotifier returns a `Notifier`
func NewNotifier(
	createAddrs, removeAddrs []string, notifyType NotifyType,
	retries int, interval int, logger *log.Logger) *Notifier {
	return &Notifier{
		createAddrs: createAddrs,
		removeAddrs: removeAddrs,
		notifyType:  notifyType,
		retries:     retries,
		internval:   interval,
		log:         logger,
	}
}

// Create sends create notifications to listeners
func (m Notifier) Create(urlValues url.Values) error {
	return nil
}

// Remove sends remove notifications to listeners
func (m Notifier) Remove(urlValues url.Values) error {
	return nil
}
