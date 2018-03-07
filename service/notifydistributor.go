package service

import (
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

// Notification is a node notification
type Notification struct {
	EventType  EventType
	Parameters string
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
	log             *log.Logger
	interval        int
}

func newNotifyDistributor(notifyEndpoints map[string]NotifyEndpoint, interval int, logger *log.Logger) *NotifyDistributor {
	return &NotifyDistributor{NotifyEndpoints: notifyEndpoints, interval: interval, log: logger}
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

	return newNotifyDistributor(notifyEndpoints, interval, logger)
}

func insertAddrStringIntoMap(tempEP map[string]map[string]string, key, addrs string) {
	for _, v := range strings.Split(addrs, ",") {
		urlObj, err := url.Parse(v)
		if err != nil {
			continue
		}
		host := urlObj.Host
		if len(host) == 0 {
			continue
		}
		if tempEP[host] == nil {
			tempEP[host] = map[string]string{}
		}
		tempEP[host][key] = v
	}
}

// NewNotifyDistributorFromEnv creates `NotifyDistributor` from environment variables
func NewNotifyDistributorFromEnv(retries, interval int, logger *log.Logger) *NotifyDistributor {
	var createServiceAddr, removeServiceAddr string
	if len(os.Getenv("DF_NOTIF_CREATE_SERVICE_URL")) > 0 {
		createServiceAddr = os.Getenv("DF_NOTIF_CREATE_SERVICE_URL")
	} else if len(os.Getenv("DF_NOTIFY_CREATE_SERVICE_URL")) > 0 {
		createServiceAddr = os.Getenv("DF_NOTIFY_CREATE_SERVICE_URL")
	} else {
		createServiceAddr = os.Getenv("DF_NOTIFICATION_URL")
	}
	if len(os.Getenv("DF_NOTIF_REMOVE_SERVICE_URL")) > 0 {
		removeServiceAddr = os.Getenv("DF_NOTIF_REMOVE_SERVICE_URL")
	} else if len(os.Getenv("DF_NOTIFY_REMOVE_SERVICE_URL")) > 0 {
		removeServiceAddr = os.Getenv("DF_NOTIFY_REMOVE_SERVICE_URL")
	} else {
		removeServiceAddr = os.Getenv("DF_NOTIFICATION_URL")
	}
	createNodeAddr := os.Getenv("DF_NOTIFY_CREATE_NODE_URL")
	removeNodeAddr := os.Getenv("DF_NOTIFY_REMOVE_NODE_URL")

	return newNotifyDistributorfromStrings(
		createServiceAddr, removeServiceAddr, createNodeAddr, removeNodeAddr, retries, interval, logger)

}

// Run starts the distributor
func (d NotifyDistributor) Run(serviceChan <-chan Notification, nodeChan <-chan Notification) {

	for _, endpoint := range d.NotifyEndpoints {
		go d.watchChannels(endpoint)
	}
	if serviceChan != nil {
		go func() {
			for n := range serviceChan {
				for _, endpoint := range d.NotifyEndpoints {
					endpoint.ServiceChan <- n
				}
			}
		}()
	}
	if nodeChan != nil {
		go func() {
			for n := range nodeChan {
				for _, endpoint := range d.NotifyEndpoints {
					endpoint.NodeChan <- n
				}
			}
		}()
	}
}

func (d NotifyDistributor) watchChannels(endpoint NotifyEndpoint) {
	for {
		select {
		case n := <-endpoint.ServiceChan:
			if n.EventType == EventTypeCreate {
				err := endpoint.ServiceNotifier.Create(n.Parameters)
				if err != nil {
					d.log.Printf("ServiceCreateNotify Error: addr: %s, params: %s back on queue", endpoint.ServiceNotifier.GetCreateAddr(), n.Parameters)
					go func() {
						time.Sleep(time.Second * time.Duration(d.interval))
						endpoint.ServiceChan <- n
					}()
				}
			} else if n.EventType == EventTypeRemove {
				err := endpoint.ServiceNotifier.Remove(n.Parameters)
				if err != nil {
					d.log.Printf("ServiceRemoveNotify Error: addr: %s, params: %s back on queue", endpoint.ServiceNotifier.GetRemoveAddr(), n.Parameters)
					go func() {
						time.Sleep(time.Second * time.Duration(d.interval))
						endpoint.ServiceChan <- n
					}()
				}
			}
		case n := <-endpoint.NodeChan:
			if n.EventType == EventTypeCreate {
				err := endpoint.NodeNotifier.Create(n.Parameters)
				if err != nil {
					d.log.Printf("NodeCreateNotify Error: addr: %s placing %s back on queue", endpoint.NodeNotifier.GetCreateAddr(), n.Parameters)
					go func() {
						time.Sleep(time.Second * time.Duration(d.interval))
						endpoint.NodeChan <- n
					}()
				}
			} else if n.EventType == EventTypeRemove {
				err := endpoint.NodeNotifier.Remove(n.Parameters)
				if err != nil {
					d.log.Printf("NodeRemoveNotify Error: addr: %s, params: %s back on queue", endpoint.NodeNotifier.GetRemoveAddr(), n.Parameters)
					go func() {
						time.Sleep(time.Second * time.Duration(d.interval))
						endpoint.NodeChan <- n
					}()
				}
			}
		}
	}
}

// HasServiceListeners when there exists service listeners
func (d NotifyDistributor) HasServiceListeners() bool {
	for _, endpoint := range d.NotifyEndpoints {
		if endpoint.ServiceNotifier != nil {
			return true
		}
	}
	return false
}

// HasNodeListeners when there exists node listeners
func (d NotifyDistributor) HasNodeListeners() bool {
	for _, endpoint := range d.NotifyEndpoints {
		if endpoint.NodeNotifier != nil {
			return true
		}
	}
	return false
}
