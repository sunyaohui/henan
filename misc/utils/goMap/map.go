//解决golang中map非并发安全问题
package goMap

import (
	"sync"
)

type Map struct {
	Map map[string]interface{}
	sync.RWMutex
}

func NewMap() *Map {
	return &Map{
		Map: make(map[string]interface{}),
	}
}

func (m *Map) Put(k string, v interface{}) {
	m.Lock()
	defer m.Unlock()
	if v == nil {
		v = "null"
	}
	m.Map[k] = v
}

func (m *Map) Remove(k string) {
	m.Lock()
	defer m.Unlock()
	delete(m.Map, k)
}

func (m *Map) Get(k string) interface{} {
	return m.Map[k]
}

func (m *Map) GetDefault(k, defaultValue string) string {
	if v, ok := m.Map[k]; ok {
		switch v.(type) {
		case string:
			return v.(string)
		default:
			return defaultValue
		}
	}
	return defaultValue
}

func (m *Map) Contains(k string) bool {
	if _, ok := m.Map[k]; ok {
		return true
	}
	return false
}
