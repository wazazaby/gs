package singleflight

import (
	"golang.org/x/sync/singleflight"
)

type Result[T any] struct {
	Val    T
	Err    error
	Shared bool
}

type Group[T any] struct {
	underlying singleflight.Group
}

func (g *Group[T]) Do(key string, fn func() (T, error)) (value T, err error, shared bool) {
	v, err, shared := g.underlying.Do(key, func() (any, error) {
		return fn()
	})
	if v != nil {
		return v.(T), err, shared
	}
	return value, err, shared
}

func (g *Group[T]) DoChan(key string, fn func() (T, error)) <-chan Result[T] {
	ch := g.underlying.DoChan(key, func() (any, error) {
		return fn()
	})

	resultCh := make(chan Result[T], 1)
	go func() {
		result := <-ch

		var v T
		if result.Val != nil {
			v = result.Val.(T)
		}

		resultCh <- Result[T]{
			Val:    v,
			Err:    result.Err,
			Shared: result.Shared,
		}
	}()

	return resultCh
}

func (g *Group[T]) Forget(key string) {
	g.underlying.Forget(key)
}
