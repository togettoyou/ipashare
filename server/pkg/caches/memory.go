package caches

import "sync"

var globalCache = &memory{kv: make(map[string]interface{}, 0)}

type memory struct {
	kv map[string]interface{}
	mu sync.RWMutex
}

func (m *memory) Set(k string, v interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.kv[k] = v
}

func (m *memory) Get(k string) interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.kv[k]
}
