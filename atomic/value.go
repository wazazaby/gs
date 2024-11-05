package atomic

import "sync/atomic"

type Value[T any] struct {
	underlying atomic.Value
}

func (v *Value[T]) CompareAndSwap(old T, new T) (swapped bool) {
	return v.underlying.CompareAndSwap(old, new)
}

func (v *Value[T]) Load() (val T, loaded bool) {
	l := v.underlying.Load()
	if l == nil {
		return val, false
	}
	return l.(T), true
}

func (v *Value[T]) Store(val T) {
	v.underlying.Store(val)
}

func (v *Value[T]) Swap(new T) (old T, swapped bool) {
	prev := v.underlying.Swap(new)
	if prev == nil {
		return old, false
	}
	return prev.(T), true
}
