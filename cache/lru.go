package cache

import (
	"github.com/wazazaby/gs/container/list"
)

// LRU is a lightweight least-recently-used tracker keyed by KT and storing VT.
//
// It keeps a map from key to a list element and a doubly linked list that
// represents recency: the front is LRU, the back is MRU. Operations update
// the list to reflect access or mutation:
//   - Upsert inserts a new entry or updates an existing entry and makes it MRU.
//   - Insert adds only when the key is missing.
//   - Update updates only when the key exists and makes it MRU.
//   - MakeMRU and MakeLRU move an existing entry to the back or front.
//   - GetMRU and GetLRU read the current extremes without modifying the list.
//   - Remove deletes an entry from both map and list.
//   - Clear resets the list and clears the map while keeping allocations for reuse.
//
// Internally, a map lookup yields the list element pointer; list operations are
// O(1). This type is not concurrency-safe and does not evict by capacity.
// Values are stored only in the list nodes; the map holds keys to list elements.
type LRU[KT comparable, VT any] struct {
	ll   list.List[VT]
	keys map[KT]*list.Element[VT]
}

func (l *LRU[KT, VT]) lazyInit() {
	if l.keys == nil {
		l.keys = make(map[KT]*list.Element[VT])
	}
}

func (l *LRU[KT, VT]) Len() int {
	return max(l.ll.Len(), len(l.keys))
}

// Clear removes all entries, keeping internal allocations for reuse.
func (l *LRU[KT, VT]) Clear() {
	l.ll.Init()
	clear(l.keys)
}

// Upsert inserts a new entry or updates an existing entry and makes it MRU.
func (l *LRU[KT, VT]) Upsert(key KT, value VT) {
	if !l.update(key, value) {
		l.lazyInit()
		l.keys[key] = l.ll.PushBack(value)
	}
}

// Insert adds a new entry if the key does not already exist.
func (l *LRU[KT, VT]) Insert(key KT, value VT) {
	if _, ok := l.keys[key]; !ok {
		l.lazyInit()
		l.keys[key] = l.ll.PushBack(value)
	}
}

func (l *LRU[KT, VT]) update(key KT, value VT) bool {
	if e := l.keys[key]; e != nil {
		e.Value = value
		l.ll.MoveToBack(e)
		return true
	}
	return false
}

// Update changes the value for an existing key and makes it MRU.
func (l *LRU[KT, VT]) Update(key KT, value VT) { l.update(key, value) }

// Remove deletes the entry for key if present.
func (l *LRU[KT, VT]) Remove(key KT) {
	if e := l.keys[key]; e != nil {
		l.ll.Remove(e)
	}
	delete(l.keys, key)
}

// MakeMRU moves key to the most-recently-used position.
func (l *LRU[KT, VT]) MakeMRU(key KT) {
	if e := l.keys[key]; e != nil {
		l.ll.MoveToBack(e)
	}
}

// GetMRU returns the most-recently-used value.
func (l *LRU[KT, VT]) GetMRU() (VT, bool) {
	if e := l.ll.Back(); e != nil {
		return e.Value, true
	}
	return *new(VT), false
}

// MakeLRU moves key to the least-recently-used position.
func (l *LRU[KT, VT]) MakeLRU(key KT) {
	if e := l.keys[key]; e != nil {
		l.ll.MoveToFront(e)
	}
}

// GetLRU returns the least-recently-used value.
func (l *LRU[KT, VT]) GetLRU() (VT, bool) {
	if e := l.ll.Front(); e != nil {
		return e.Value, true
	}
	return *new(VT), false
}
