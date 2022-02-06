package lru

import "sync"

type Cache interface {
	Set(key string, value string) bool
	Get(key string) (string, bool)
	Clear()
}

type lruCache struct {
	sync.Mutex
	capacity int
	queue    List
	items    map[string]*ListItem
}

func (c *lruCache) Set(key string, value string) bool {
	defer c.Unlock()

	c.Lock()
	_, exists := c.items[key]

	if exists {
		c.queue.Remove(c.queue.Front())
	}

	if len(c.items) > c.capacity {
		c.queue.Remove(c.queue.Back())
	}

	c.queue.PushFront(value)

	c.items[key] = c.queue.Front()

	return exists
}

func (c *lruCache) Get(key string) (string, bool) {
	defer c.Unlock()
	c.Lock()
	val, exists := c.items[key]

	if exists {
		c.queue.MoveToFront(val)

		return val.Value, exists
	}

	return "", false
}

func (c *lruCache) Clear() {
	defer c.Unlock()
	c.Lock()

	c.queue = NewList()
	c.items = make(map[string]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[string]*ListItem, capacity),
	}
}
