package service

// ServicCacher caches sevices
type ServicCacher interface {
	InsertAndCheck(ss SwarmService, eventType EventType) bool
}

// ServicCache implements `ServicCacher`
// Not threadsafe!
type ServicCache struct {
	Cache map[string]SwarmServiceMini
}

// InsertAndCheck inserts `SwarmService` into cache if the service is updated or created
// If the service is removed, it will be removed from the cache
// If the service is new, created, or removed, `InsertAndCheck` returns true.
func (c ServicCache) InsertAndCheck(ss SwarmService, eventType EventType) bool {
	return false
}
