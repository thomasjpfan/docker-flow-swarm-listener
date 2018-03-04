package service

import (
	"github.com/docker/docker/api/types/swarm"
)

// NodeCacher caches sevices
type NodeCacher interface {
	InsertAndCheck(ss swarm.Node, eventType NodeEventType) bool
}

// NodeCache implements `NodeCacher`
// Not threadsafe!
type NodeCache struct {
	Cache map[string]SwarmServiceMini
}

// InsertAndCheck inserts `swarm.Node` into cache if the service is updated or created
// If the service is removed, it will be removed from the cache
// If the service is new, created, or removed, `InsertAndCheck` returns true.
func (c NodeCache) InsertAndCheck(ss swarm.Node, eventType NodeEventType) bool {
	return false
}
