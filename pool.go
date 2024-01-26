package gs

import "sync"

type Pool[T any] struct {
	underlying sync.Pool
}

func NewPool[T any](new func() T) Pool[T] {
	return Pool[T]{
		underlying: sync.Pool{
			New: func() any { return new() },
		},
	}
}

func (p *Pool[T]) Get() T {
	return p.underlying.Get().(T)
}

func (p *Pool[T]) Put(x T) {
	p.underlying.Put(x)
}
