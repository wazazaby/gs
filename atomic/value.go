package atomic

import "sync/atomic"

type Value[T any] struct {
	underlying atomic.Value
}

func (v *Value[T]) CompareAndSwap(old T, new T) (swapped bool) {
	return v.underlying.CompareAndSwap(old, new)
}

func (v *Value[T]) Load() T {
	val := v.underlying.Load()
	if val == nil {
		var zero T
		return zero
	}
	return val.(T)
}

func (v *Value[T]) Store(val T) {
	v.underlying.Store(val)
}

func (v *Value[T]) Swap(new T) T {
	old := v.underlying.Swap(new)
	if old == nil {
		var zero T
		return zero
	}
	return old.(T)
}
