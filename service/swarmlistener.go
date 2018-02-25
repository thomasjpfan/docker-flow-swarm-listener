package service

import (
	"context"
	"log"
	"os"
)

// SwarmListening provides public api for interacting with swarm listener
type SwarmListening interface {
	Run()
	NotifyServices()
	GetServicesParameters(ctx context.Context) ([]map[string]string, error)
}

// SwarmListener provides public api
type SwarmListener struct {
	SSListener         SwarmServiceListening
	SSClient           SwarmServiceInspector
	SSCache            SwarmServiceCacher
	SSEventChan        chan Event
	SSNotificationChan chan Notification

	NodeListener         NodeListening
	NodeClient           NodeInspector
	NodeCache            NodeCacher
	NodeEventChan        chan Event
	NodeNotificationChan chan Notification

	NotifyDistributor NotifyDistributing

	IncludeNodeInfo bool
	IgnoreKey       string
	IncludeKey      string
	Log             *log.Logger
}

func newSwarmListener(
	ssListener SwarmServiceListening,
	ssClient SwarmServiceInspector,
	ssCache SwarmServiceCacher,

	nodeListener NodeListening,
	nodeClient NodeInspector,
	nodeCache NodeCacher,

	notifyDistributor NotifyDistributing,
	includeNodeInfo bool,
	ignoreKey string,
	includeKey string,
	logger *log.Logger,
) *SwarmListener {

	return &SwarmListener{
		SSListener:           ssListener,
		SSClient:             ssClient,
		SSCache:              ssCache,
		SSEventChan:          make(chan Event),
		SSNotificationChan:   make(chan Notification),
		NodeListener:         nodeListener,
		NodeClient:           nodeClient,
		NodeCache:            nodeCache,
		NodeEventChan:        make(chan Event),
		NodeNotificationChan: make(chan Notification),
		NotifyDistributor:    notifyDistributor,
		IncludeNodeInfo:      includeNodeInfo,
		IgnoreKey:            ignoreKey,
		IncludeKey:           includeKey,
		Log:                  logger,
	}
}

// NewSwarmListenerFromEnv creats `SwarmListener` from environment variables
func NewSwarmListenerFromEnv(retries, interval int, logger *log.Logger) (*SwarmListener, error) {
	ignoreKey := os.Getenv("DF_NOTIFY_LABEL")
	includeNodeInfo := os.Getenv("DF_INCLUDE_NODE_IP_INFO") == "true"

	dockerClient, err := NewDockerClientFromEnv()
	if err != nil {
		return nil, err
	}
	ssListener := NewSwarmServiceListener(dockerClient, logger)
	ssClient := NewSwarmServiceClient(dockerClient, ignoreKey, "com.df.scrapeNetwork")
	ssCache := NewSwarmServiceCache()

	nodeListener := NewNodeListener(dockerClient, logger)
	nodeClient := NewNodeClient(dockerClient)
	nodeCache := NewNodeCache()

	notifyDistributor := NewNotifyDistributorFromEnv(retries, interval, logger)

	return newSwarmListener(
		ssListener,
		ssClient,
		ssCache,
		nodeListener,
		nodeClient,
		nodeCache,
		notifyDistributor,
		includeNodeInfo,
		ignoreKey,
		"com.docker.stack.namespace",
		logger,
	), nil

}

// Run starts swarm listener
func (l *SwarmListener) Run() {
	l.connectServiceChannels()
	l.connectNodeChannels()

	if l.SSEventChan != nil {
		l.SSListener.ListenForServiceEvents(l.SSEventChan)
	}
	if l.NodeEventChan != nil {
		l.NodeListener.ListenForNodeEvents(l.NodeEventChan)
	}

	l.NotifyDistributor.Run(l.SSNotificationChan, l.NodeNotificationChan)
}

func (l *SwarmListener) connectServiceChannels() {

	// Remove service channels if there are no service listeners
	if !l.NotifyDistributor.HasServiceListeners() {
		l.SSEventChan = nil
		l.SSNotificationChan = nil
		return
	}

	go func() {
		for event := range l.SSEventChan {
			if event.Type == EventTypeCreate {
				service, err := l.SSClient.SwarmServiceInspect(event.ID, l.IncludeNodeInfo)
				if err != nil {
					l.Log.Printf("ERROR: %v", err)
					continue
				}
				// Ignored service (filtered by `com.df.notify`)
				if service == nil {
					continue
				}
				ssm := MinifySwarmService(*service, l.IgnoreKey, l.IncludeKey)

				// Store in cache
				isUpdated := l.SSCache.InsertAndCheck(ssm)
				if !isUpdated {
					continue
				}
				params := GetSwarmServiceMiniCreateParameters(ssm)
				paramsEncoded := ConvertMapStringStringToURLValues(params).Encode()
				l.placeOnNotificationChan(l.SSNotificationChan, event.Type, paramsEncoded)
			} else {
				// EventTypeRemove
				ssm, ok := l.SSCache.Get(event.ID)
				if !ok {
					continue
				}
				l.SSCache.Delete(ssm.ID)
				params := GetSwarmServiceMiniRemoveParameters(ssm)
				paramsEncoded := ConvertMapStringStringToURLValues(params).Encode()
				l.placeOnNotificationChan(l.SSNotificationChan, event.Type, paramsEncoded)
			}
		}
	}()
}

func (l *SwarmListener) connectNodeChannels() {

	// Remove node channels if there are no service listeners
	if !l.NotifyDistributor.HasNodeListeners() {
		l.NodeEventChan = nil
		l.NodeNotificationChan = nil
		return
	}

	go func() {
		for event := range l.NodeEventChan {
			if event.Type == EventTypeCreate {
				node, err := l.NodeClient.NodeInspect(event.ID)
				if err != nil {
					l.Log.Printf("ERROR: %v", err)
					continue
				}
				nm := MinifyNode(node)

				// Store in cache
				isUpdated := l.NodeCache.InsertAndCheck(nm)
				if !isUpdated {
					continue
				}
				params := GetNodeMiniCreateParameters(nm)
				paramsEncoded := ConvertMapStringStringToURLValues(params).Encode()
				l.placeOnNotificationChan(l.NodeNotificationChan, event.Type, paramsEncoded)
			} else {
				// EventTypeRemove
				nm, ok := l.NodeCache.Get(event.ID)
				if !ok {
					continue
				}
				l.NodeCache.Delete(nm.ID)
				params := GetNodeMiniRemoveParameters(nm)
				paramsEncoded := ConvertMapStringStringToURLValues(params).Encode()
				l.placeOnNotificationChan(l.NodeNotificationChan, event.Type, paramsEncoded)
			}
		}
	}()
}

// NotifyServices places all services on queue to notify services
// Ignoring the cache
func (l SwarmListener) NotifyServices() {
	services, err := l.SSClient.SwarmServiceList(context.Background(), l.IncludeNodeInfo)
	if err != nil {
		l.Log.Printf("ERROR: NotifyService, %v", err)
		return
	}
	for _, s := range services {
		ssm := MinifySwarmService(s, l.IgnoreKey, l.IncludeKey)
		params := GetSwarmServiceMiniCreateParameters(ssm)
		paramsEncoded := ConvertMapStringStringToURLValues(params).Encode()
		l.placeOnNotificationChan(l.SSNotificationChan, EventTypeCreate, paramsEncoded)
	}
}

func (l SwarmListener) placeOnNotificationChan(notiChan chan<- Notification, eventType EventType, parameters string) {
	go func() {
		notiChan <- Notification{
			EventType:  eventType,
			Parameters: parameters,
		}
	}()
}

// GetServicesParameters get all services
func (l SwarmListener) GetServicesParameters(ctx context.Context) ([]map[string]string, error) {
	services, err := l.SSClient.SwarmServiceList(ctx, l.IncludeNodeInfo)
	if err != nil {
		return []map[string]string{}, err
	}
	params := []map[string]string{}
	for _, s := range services {
		ssm := MinifySwarmService(s, l.IgnoreKey, l.IncludeKey)
		newParams := GetSwarmServiceMiniCreateParameters(ssm)
		if len(newParams) > 0 {
			params = append(params, newParams)
		}
	}
	return params, nil
}
