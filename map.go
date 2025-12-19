package gs

import (
	"github.com/go4org/hashtriemap"
)

type Map[K comparable, V any] struct {
	u hashtriemap.HashTrieMap[K, V]
}

func (m *Map[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	return m.u.CompareAndDelete(key, old)
}

func (m *Map[K, V]) CompareAndSwap(key K, old V, new V) (swaped bool) {
	return m.u.CompareAndSwap(key, old, new)
}

func (m *Map[K, _]) Delete(key K) {
	m.u.Delete(key)
}

func (m *Map[K, V]) Load(key K) (value V, found bool) {
	return m.u.Load(key)
}

func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	return m.u.LoadAndDelete(key)
}

func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	return m.u.LoadOrStore(key, value)
}

func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.u.Range(f)
}

func (m *Map[K, V]) Store(key K, value V) {
	m.u.Store(key, value)
}

func (m *Map[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	return m.u.Swap(key, value)
}
