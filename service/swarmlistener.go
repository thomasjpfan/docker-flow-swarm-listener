package service

// SwarmListening provides public api for interacting with swarm listener
type SwarmListening interface {
	Run()
	NotifyServices()
	GetServices() ([]SwarmService, error)
}

// SwarmServiceMiniEvent are used to denote if `Service` should be
// ignored by cache
type SwarmServiceMiniEvent struct {
	Service     SwarmServiceMini
	IgnoreCache bool
}

// NodeMiniEvent are used to denote if `Node` should be
// ignored by cache
type NodeMiniEvent struct {
	Node        NodeMini
	IgnoreCache bool
}

// SwarmListener provides public api
type SwarmListener struct {
	SSListener         SwarmServiceListening
	SSClient           SwarmServiceInspector
	SSCache            SwarmServiceCacher
	SSEventChan        chan Event
	SSMiniChan         chan SwarmServiceMiniEvent
	SSNotificationChan chan Notification

	NodeListener         NodeListening
	NodeClient           NodeInspector
	NodeCache            NodeCacher
	NodeEvent            chan Event
	NodeMiniChan         chan NodeMiniEvent
	NodeNotificationChan chan Notification

	NotifyDistributor NotifyDistributing
}

func newSwarmListener(
	ssListener SwarmServiceListening,
	ssClient SwarmServiceInspector,
	ssCache SwarmServiceCacher,
	ssEventChan chan Event,
	ssMiniChan chan SwarmServiceMiniEvent,

	nodeListener NodeListening,
	nodeClient NodeInspector,
	nodeCache NodeCacher,
	nodeEvent chan Event,
	nodeMiniChan chan NodeMiniEvent,
) *SwarmListener {
	return &SwarmListener{
		SSListener:   ssListener,
		SSClient:     ssClient,
		SSCache:      ssCache,
		SSEventChan:  ssEventChan,
		SSMiniChan:   ssMiniChan,
		NodeListener: nodeListener,
		NodeClient:   nodeClient,
		NodeCache:    nodeCache,
		NodeEvent:    nodeEvent,
		NodeMiniChan: nodeMiniChan,
	}
}

// NewSwarmListenerFromEnv creats `SwarmListener` from environment variables
func NewSwarmListenerFromEnv() *SwarmListener {
	return nil
}

// Run starts swarm listener
func (l SwarmListener) Run() {

}

// NotifyServices places all services on queue to notify services
// Ignoring the cache
func (l SwarmListener) NotifyServices() {

}

// GetServices get all services
func (l SwarmListener) GetServices() ([]SwarmService, error) {
	return []SwarmService{}, nil
}
