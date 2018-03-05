package service

// SwarmListening provides public api for interacting with swarm listener
type SwarmListening interface {
	Run()
	NotifyServices()
	GetServices() ([]SwarmService, error)
}

// SwarmListener provides public api
type SwarmListener struct {
	SSListener         SwarmServiceListening
	SSClient           SwarmServiceInspector
	SSCache            SwarmServiceCacher
	SSEventChan        chan Event
	SSMiniChan         chan SwarmServiceMini
	SSNotificationChan chan Notification

	NodeListener         NodeListening
	NodeClient           NodeInspector
	NodeCache            NodeCacher
	NodeEvent            chan Event
	NodeMiniChan         chan NodeMini
	NodeNotificationChan chan Notification

	NotifyDistributor NotifyDistributing
}

func newSwarmListener(
	ssListener SwarmServiceListening,
	ssClient SwarmServiceInspector,
	ssCache SwarmServiceCacher,
	ssEventChan chan Event,
	ssMiniChan chan SwarmServiceMini,

	nodeListener NodeListening,
	nodeClient NodeInspector,
	nodeCache NodeCacher,
	nodeEvent chan Event,
	nodeMiniChan chan NodeMini,
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
func (l SwarmListener) NotifyServices() {

}

// GetServices get all services
func (l SwarmListener) GetServices() ([]SwarmService, error) {
	return []SwarmService{}, nil
}
