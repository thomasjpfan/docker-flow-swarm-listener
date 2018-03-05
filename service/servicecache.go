package service

// SwarmServiceCacher caches sevices
type SwarmServiceCacher interface {
	InsertAndCheck(ss SwarmServiceMini) bool
	GetAndRemove(ID string) *SwarmServiceMini
}

// SwarmServiceCache implements `SwarmServiceCacher`
// Not threadsafe!
type SwarmServiceCache struct {
	Cache map[string]SwarmServiceMini
}

// InsertAndCheck inserts `SwarmServiceMini` into cache
// If the service is new, created, `InsertAndCheck` returns true.
func (c SwarmServiceCache) InsertAndCheck(ss SwarmServiceMini) bool {
	return false
}

// GetAndRemove removes `SwarmServiceMini` from cache
// If service was in cache, return the corresponding `SwarmServiceMini`
// IF service is not in cache, return nil
func GetAndRemove(ID string) *SwarmServiceMini {
	return nil
}
