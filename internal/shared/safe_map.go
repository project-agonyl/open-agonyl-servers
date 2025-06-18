package shared

import "sync"

type SafeMap[K comparable, V any] struct {
	m sync.Map
}

func (m *SafeMap[K, V]) Store(k K, v V) {
	m.m.Store(k, v)
}

func (m *SafeMap[K, V]) Set(k K, v V) {
	m.Store(k, v)
}

func (m *SafeMap[K, V]) Load(k K) (V, bool) {
	v, found := m.m.Load(k)
	if !found {
		var empty V
		return empty, found
	}

	return v.(V), found
}

func (m *SafeMap[K, V]) Get(k K) (V, bool) {
	return m.Load(k)
}

func (m *SafeMap[K, V]) Delete(k K) {
	m.m.Delete(k)
}

func (m *SafeMap[K, V]) Range(f func(k K, v V) bool) {
	m.m.Range(func(k, v interface{}) bool {
		return f(k.(K), v.(V))
	})
}

func (m *SafeMap[K, V]) Len() int {
	length := 0
	m.Range(func(k K, v V) bool {
		length++
		return true
	})

	return length
}

func (m *SafeMap[K, V]) Has(k K) bool {
	_, found := m.Load(k)
	return found
}

func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{}
}
