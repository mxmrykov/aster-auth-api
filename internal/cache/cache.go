package cache

import (
	"sync"
	"time"
)

type ICache interface {
	Get(key string) *Client
	Set(key string, client *Client) bool
}

type Client struct {
	IAID string

	// rate limiting
	rateLimitRemain uint8
	lastReq         time.Time

	// Client properties
	IsBanned bool
	Login    string

	// Inner properties
	LastUpdated time.Time
}

type Cache struct {
	Storage   map[string]*Client
	RWm       *sync.RWMutex
	TempUsers []string
}

func NewCache() *Cache {
	return &Cache{
		Storage:   make(map[string]*Client),
		RWm:       new(sync.RWMutex),
		TempUsers: make([]string, 0),
	}
}

func (c *Cache) Get(key string) *Client {
	c.RWm.RLock()
	defer c.RWm.RUnlock()
	if client, ok := c.Storage[key]; ok {
		return client
	}

	return nil
}

func (c *Cache) Set(key string, client *Client) bool {
	c.RWm.Lock()
	defer c.RWm.Unlock()
	if _, ok := c.Storage[key]; ok {
		return false
	}

	c.Storage[key] = client
	return true
}
