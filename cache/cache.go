package cache

import (
	"container/list"
	"log"
	"sync"
)

type entry struct {
	key   string
	value []byte
}

type Cache struct {
	mu      sync.Mutex
	maxSize int
	ll      *list.List
	cache   map[string]*list.Element
}

func NewCache(maxSize int) *Cache {
	return &Cache{
		maxSize: maxSize,
		ll:      list.New(),
		cache:   make(map[string]*list.Element),
	}
}

func (c *Cache) Set(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, found := c.cache[key]; found {
		c.ll.MoveToFront(element)
		element.Value.(*entry).value = value
		return
	}

	newEntry := &entry{key: key, value: value}
	element := c.ll.PushFront(newEntry)
	c.cache[key] = element

	if c.ll.Len() > c.maxSize && c.maxSize > 0 {
		c.removeOldest()
	}
}

func (c *Cache) removeOldest() {
	element := c.ll.Back()
	if element != nil {
		c.ll.Remove(element)
		keyToRemove := element.Value.(*entry).key
		delete(c.cache, keyToRemove)
		log.Printf("Evicted oldest key: %s", keyToRemove)
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, found := c.cache[key]; found {
		c.ll.MoveToFront(element)
		return element.Value.(*entry).value, true
	}

	return nil, false
}

func (c *Cache) Delete(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, found := c.cache[key]; found {
		c.ll.Remove(element)
		delete(c.cache, key)
		return true
	}
	return false
}
