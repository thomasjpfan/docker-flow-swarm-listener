package service

import (
	"log"
	"net/url"
)

// NotificationSender sends notifications to listeners
type NotificationSender interface {
	Create(urlValues url.Values) error
	Remove(urlValues url.Values) error
	HasListeners() bool
}

// Notifier implements `NotificationSender`
type Notifier struct {
	createAddrs  []string
	removeAddrs  []string
	notifierType string
	retries      int
	internval    int
	log          *log.Logger
}

// NewNotifier returns a `Notifier`
func NewNotifier(
	createAddrs, removeAddrs []string, notifierType string,
	retries int, interval int, logger *log.Logger) *Notifier {
	return &Notifier{
		createAddrs:  createAddrs,
		removeAddrs:  removeAddrs,
		notifierType: notifierType,
		retries:      retries,
		internval:    interval,
		log:          logger,
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

// HasListeners when there are listeners
func (m Notifier) HasListeners() bool {
	return (len(m.createAddrs) + len(m.removeAddrs)) > 0
}
