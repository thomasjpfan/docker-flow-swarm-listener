package service

// NodeCacher caches sevices
type NodeCacher interface {
	InsertAndCheck(n NodeMini) bool
	GetAndRemove(ID string) (NodeMini, bool)
}

// NodeCache implements `NodeCacher`
// Not threadsafe!
type NodeCache struct {
	cache map[string]NodeMini
}

// NewNodeCache creates a new `NewNodeCache`
func NewNodeCache() *NodeCache {
	return &NodeCache{
		cache: map[string]NodeMini{},
	}
}

// InsertAndCheck inserts `NodeMini` into cache
// If the node is new or updated `InsertAndCheck` returns true.
func (c *NodeCache) InsertAndCheck(n NodeMini) bool {
	cachedNode, ok := c.cache[n.ID]
	c.cache[n.ID] = n

	return !ok || !n.Equal(cachedNode)
}

// GetAndRemove removes `NodeMini` from cache
// If node was in cache, return the corresponding `NodeMini`,
// remove from cache, and return true
// If node is not in cache, return false
func (c *NodeCache) GetAndRemove(ID string) (NodeMini, bool) {
	if cachedNode, ok := c.cache[ID]; ok {
		delete(c.cache, ID)
		return cachedNode, true
	}
	return NodeMini{}, false
}

func (c NodeCache) get(ID string) (NodeMini, bool) {
	v, ok := c.cache[ID]
	return v, ok
}
