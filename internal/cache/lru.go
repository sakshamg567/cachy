package cache

import (
	"errors"
	"sync"
)

type dllNode struct {
	key   string
	value string
	next  *dllNode
	prev  *dllNode
}

type DLL struct {
	front *dllNode
	back  *dllNode
}

func (d *DLL) moveToFront(node *dllNode) {
	if d.front == node {
		return
	}

	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	if d.back == node {
		d.back = node.prev
	}

	node.prev = nil
	node.next = d.front
	if d.front != nil {
		d.front.prev = node
	}
	d.front = node
	if d.back == nil {
		d.back = d.front
	}
}

func (d *DLL) evictLRU() *dllNode {
	if d.back == nil {
		return nil
	}
	node := d.back
	if node.prev != nil {
		node.prev.next = nil
	}
	d.back = node.prev
	if d.back == nil {
		d.front = nil
	}
	return node
}

type LruCache struct {
	capacity int
	cache    map[string]*dllNode
	dll      *DLL
	mu       sync.RWMutex
}

func NewLruCache(cap int) *LruCache {
	return &LruCache{
		capacity: cap,
		cache:    map[string]*dllNode{},
		dll:      &DLL{},
	}
}

const (
	ERRKEYNOTFOUND = "key not found"
)

func (c *LruCache) get(key string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if node, ok := c.cache[key]; ok {
		c.dll.moveToFront(node)
		return node.value, nil
	}
	return "", errors.New(ERRKEYNOTFOUND)
}

func (c *LruCache) set(key string, value string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if node, ok := c.cache[key]; ok {
		node.value = value
		c.dll.moveToFront(node)
		return true
	}
	if len(c.cache) >= c.capacity {
		evicted := c.dll.evictLRU()
		if evicted != nil {
			delete(c.cache, evicted.key)
		}
	}

	newNode := &dllNode{
		key:   string(key),
		value: value,
	}
	c.dll.moveToFront(newNode)
	c.cache[key] = newNode
	return true
}

func (c *LruCache) GetAllKeys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.cache))
	for k := range c.cache {
		keys = append(keys, k)
	}
	return keys
}

func (c *LruCache) delete(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, ok := c.cache[key]
	if !ok {
		return false
	}

	if node.prev != nil {
		node.prev.next = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	}

	if c.dll.front == node {
		c.dll.front = node.next
	}

	if c.dll.back == node {
		c.dll.back = node.prev
	}

	delete(c.cache, key)
	return true
}
