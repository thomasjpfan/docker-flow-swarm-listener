package service

// SwarmServiceCacher caches sevices
type SwarmServiceCacher interface {
	InsertAndCheck(ss SwarmServiceMini) bool
	GetAndRemove(ID string) (SwarmServiceMini, bool)
}

// SwarmServiceCache implements `SwarmServiceCacher`
// Not threadsafe!
type SwarmServiceCache struct {
	cache map[string]SwarmServiceMini
}

// NewSwarmServiceCache creates a new `NewSwarmServiceCache`
func NewSwarmServiceCache() *SwarmServiceCache {
	return &SwarmServiceCache{
		cache: map[string]SwarmServiceMini{},
	}
}

// InsertAndCheck inserts `SwarmServiceMini` into cache
// If the service is new or updated `InsertAndCheck` returns true.
func (c *SwarmServiceCache) InsertAndCheck(ss SwarmServiceMini) bool {
	return false
}

// GetAndRemove removes `SwarmServiceMini` from cache
// If service was in cache, return the corresponding `SwarmServiceMini`,
// remove from cache, and return true
// If service is not in cache, return false
func (c *SwarmServiceCache) GetAndRemove(ID string) (SwarmServiceMini, bool) {
	return SwarmServiceMini{}, false
}

func (c SwarmServiceCache) get(ID string) (SwarmServiceMini, bool) {
	v, ok := c.cache[ID]
	return v, ok
}
