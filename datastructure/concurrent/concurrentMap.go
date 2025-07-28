package concurrent

import "sync"

type ConcurrentMap[Key string | int | uint, Value interface{}] struct {
	lock    sync.RWMutex
	thisMap map[Key]Value
}

func (m *ConcurrentMap[Key, Value]) Get(key Key) Value {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.thisMap[key]
}

func (m *ConcurrentMap[Key, Value]) Set(key Key, value Value) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.thisMap[key] = value
}

func NewConcurrentMap[Key string | int | uint, Value interface{}]() *ConcurrentMap[Key, Value] {
	return &ConcurrentMap[Key, Value]{
		thisMap: make(map[Key]Value),
		lock:    sync.RWMutex{},
	}
}
