package cache

import (
	"container/list"
	"encoding/gob"
	"log"
	"os"
	"sync"
	"time"
)

const SaveFilePath = "cache.gob"

type entry struct {
	key        string
	value      []byte
	expiration time.Time
}

type Cache struct {
	mu      sync.Mutex
	maxSize int
	ll      *list.List
	cache   map[string]*list.Element
}

func NewCache(maxSize int, cleanupInterval time.Duration) *Cache {
	c := &Cache{
		maxSize: maxSize,
		ll:      list.New(),
		cache:   make(map[string]*list.Element),
	}

	if err := c.LoadFromFile(SaveFilePath); err != nil {
		log.Println("Could not load cache from file, starting fresh:", err)
	}

	if cleanupInterval > 0 {
		go c.cleanupLoop(cleanupInterval)
	}

	go c.saveLoop(5*time.Minute, SaveFilePath)

	return c
}

func (c *Cache) saveLoop(interval time.Duration, path string) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		if err := c.SaveToFile(path); err != nil {
			log.Println("Error saving cache to file:", err)
		}
	}
}

func (c *Cache) Set(key string, value []byte, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiration time.Time
	if duration > 0 {
		expiration = time.Now().Add(duration)
	}

	if element, found := c.cache[key]; found {
		c.ll.MoveToFront(element)
		entry := element.Value.(*entry)
		entry.value = value
		entry.expiration = expiration
		return
	}

	newEntry := &entry{key: key, value: value, expiration: expiration}
	element := c.ll.PushFront(newEntry)
	c.cache[key] = element

	if c.maxSize > 0 && c.ll.Len() > c.maxSize {
		c.removeOldest()
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	element, found := c.cache[key]
	if !found {
		return nil, false
	}

	entry := element.Value.(*entry)

	if !entry.expiration.IsZero() && time.Now().After(entry.expiration) {
		c.ll.Remove(element)
		delete(c.cache, key)
		log.Printf("Passively deleted expired key: %s", key)
		return nil, false
	}

	c.ll.MoveToFront(element)
	return entry.value, true
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

func (c *Cache) removeOldest() {
	element := c.ll.Back()
	if element != nil {
		c.ll.Remove(element)
		keyToRemove := element.Value.(*entry).key
		delete(c.cache, keyToRemove)
		log.Printf("Evicted oldest key: %s", keyToRemove)
	}
}

func (c *Cache) cleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		c.deleteExpired()
	}
}

func (c *Cache) deleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	deletedCount := 0
	for key, element := range c.cache {
		entry := element.Value.(*entry)

		if !entry.expiration.IsZero() && time.Now().After(entry.expiration) {
			c.ll.Remove(element)
			delete(c.cache, key)
			deletedCount++
		}
	}

	if deletedCount > 0 {
		log.Printf("Background cleanup deleted %d expired keys.", deletedCount)
	}
}

func (c *Cache) SaveToFile(path string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	tempStorage := make(map[string]entry)
	for key, element := range c.cache {
		tempStorage[key] = *element.Value.(*entry)
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(tempStorage); err != nil {
		return err
	}

	log.Printf("Cache successfully saved to file: %s", path)
	return nil
}

func (c *Cache) LoadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)

	tempStorage := make(map[string]entry)
	if err := decoder.Decode(&tempStorage); err != nil {
		return err
	}

	for key, entryData := range tempStorage {
		var duration time.Duration
		if !entryData.expiration.IsZero() {
			duration = time.Until(entryData.expiration)
		}
		if duration > 0 {
			c.Set(key, entryData.value, duration)
		}
	}

	log.Printf("Cache successfully loaded and rebuilt from file: %s", path)
	return nil
}
