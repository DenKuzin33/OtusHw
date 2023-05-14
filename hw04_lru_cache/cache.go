package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	node, exist := cache.items[key]

	if exist {
		node.Value = value
		cache.queue.MoveToFront(node)
		return true
	}

	if cache.queue.Len() >= cache.capacity {
		back := cache.queue.Back()
		cache.queue.Remove(back)
		delete(cache.items, back.Key)
	}

	cache.items[key] = cache.queue.PushFront(key, value)
	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	node, exist := cache.items[key]
	if exist {
		cache.queue.MoveToFront(node)
		return node.Value, true
	}
	return nil, false
}

func (cache *lruCache) Clear() {
	cache = &lruCache{
		capacity: cache.capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, cache.capacity),
	}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
