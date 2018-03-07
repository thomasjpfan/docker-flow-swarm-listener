package service

import (
	"log"
	"net/url"
	"strings"
)

// Notification is a node notification
type Notification struct {
	eventType  EventType
	parameters string
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

// NotifyDistributor distributes service and node notifications to `NotifyEndpoints`
// `NotifyEndpoints` are keyed by hostname to send notifications to
type NotifyDistributor struct {
	NotifyEndpoints map[string]NotifyEndpoint
}

func newNotifyDistributor(notifyEndpoints map[string]NotifyEndpoint) *NotifyDistributor {
	return &NotifyDistributor{NotifyEndpoints: notifyEndpoints}
}

func newNotifyDistributorfromStrings(serviceCreateAddrs, serviceRemoveAddrs, nodeCreateAddrs, nodeRemoveAddrs string, retries, interval int, logger *log.Logger) *NotifyDistributor {
	tempNotifyEP := map[string]map[string]string{}

	insertAddrStringIntoMap(tempNotifyEP, "createService", serviceCreateAddrs)
	insertAddrStringIntoMap(tempNotifyEP, "removeService", serviceRemoveAddrs)
	insertAddrStringIntoMap(tempNotifyEP, "createNode", nodeCreateAddrs)
	insertAddrStringIntoMap(tempNotifyEP, "removeNode", nodeRemoveAddrs)

	notifyEndpoints := map[string]NotifyEndpoint{}

	for hostname, addrMap := range tempNotifyEP {
		ep := NotifyEndpoint{}
		if len(addrMap["createService"]) > 0 || len(addrMap["removeService"]) > 0 {
			ep.ServiceChan = make(chan Notification)
			ep.ServiceNotifier = NewNotifier(
				addrMap["createService"],
				addrMap["removeService"],
				"service",
				retries,
				interval,
				logger,
			)
		}
		if len(addrMap["createNode"]) > 0 || len(addrMap["removeNode"]) > 0 {
			ep.NodeChan = make(chan Notification)
			ep.NodeNotifier = NewNotifier(
				addrMap["createNode"],
				addrMap["removeNode"],
				"node",
				retries,
				interval,
				logger,
			)
		}
		notifyEndpoints[hostname] = ep
	}

	return newNotifyDistributor(notifyEndpoints)
}

func insertAddrStringIntoMap(tempEP map[string]map[string]string, key, addrs string) {
	for _, v := range strings.Split(addrs, ",") {
		urlObj, err := url.Parse(v)
		if err != nil {
			continue
		}
		hostname := urlObj.Host
		if len(hostname) == 0 {
			continue
		}
		if tempEP[hostname] == nil {
			tempEP[hostname] = map[string]string{}
		}
		tempEP[hostname][key] = v
	}
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
