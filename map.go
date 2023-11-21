package gs

import (
	"sync"
)

type Map[K comparable, V any] struct {
	underlying sync.Map
}

func (m *Map[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	return m.underlying.CompareAndDelete(key, old)
}

func (m *Map[K, V]) CompareAndSwap(key K, old V, new V) (swaped bool) {
	return m.underlying.CompareAndSwap(key, old, new)
}

func (m *Map[K, _]) Delete(key K) {
	m.underlying.Delete(key)
}

func (m *Map[K, V]) Load(key K) (value V, found bool) {
	v, found := m.underlying.Load(key)
	if !found {
		return value, found
	}
	return v.(V), found
}

func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, loaded := m.underlying.LoadAndDelete(key)
	if !loaded {
		return value, loaded
	}
	return v.(V), loaded
}

func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	v, loaded := m.underlying.LoadOrStore(key, value)
	return v.(V), loaded
}

func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.underlying.Range(func(k, v any) bool {
		return f(k.(K), v.(V))
	})
}

func (m *Map[K, V]) Store(key K, value V) {
	m.underlying.Store(key, value)
}

func (m *Map[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	v, loaded := m.underlying.Swap(key, value)
	if !loaded {
		return previous, loaded
	}
	return v.(V), loaded
}
