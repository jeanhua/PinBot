package concurrent

import "sync"

type ConcurrentMap[Key comparable, Value interface{}] struct {
	lock    sync.RWMutex
	thisMap map[Key]Value
}

func (m *ConcurrentMap[Key, Value]) Get(key Key) (Value, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	m.check()
	value, ok := m.thisMap[key]
	return value, ok
}

func (m *ConcurrentMap[Key, Value]) Set(key Key, value Value) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.check()
	m.thisMap[key] = value
}

func (m *ConcurrentMap[Key, Value]) Delete(key Key) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.check()
	delete(m.thisMap, key)
}

func NewConcurrentMap[Key string | int | uint, Value interface{}]() *ConcurrentMap[Key, Value] {
	return &ConcurrentMap[Key, Value]{
		thisMap: make(map[Key]Value),
		lock:    sync.RWMutex{},
	}
}

func (m *ConcurrentMap[Key, Value]) check() {
	if m.thisMap == nil {
		m.thisMap = make(map[Key]Value)
	}
}
