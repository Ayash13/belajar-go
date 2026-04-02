package main

import (
	"fmt"
	"sync"
)

// Simple in-memory cache to demonstrate caching fundamentals
type Cache struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]string),
	}
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists := c.data[key]
	return value, exists
}

func main() {
	cache := NewCache()

	fmt.Println("Storing data in cache...")
	cache.Set("user:1", "John Doe")

	fmt.Println("Retrieving data from cache...")
	if val, ok := cache.Get("user:1"); ok {
		fmt.Printf("Found: %s\n", val)
	} else {
		fmt.Println("Not found")
	}
}
