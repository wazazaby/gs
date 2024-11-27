package singleflight

import (
	"golang.org/x/sync/singleflight"
)

func mapResult[T any](inCh <-chan singleflight.Result, outCh chan<- Result[T]) {
	result := <-inCh
	var v T
	if result.Val != nil {
		v = result.Val.(T)
	}
	outCh <- Result[T]{
		Val:    v,
		Err:    result.Err,
		Shared: result.Shared,
	}
}

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
	resultCh := g.underlying.DoChan(key, func() (any, error) {
		return fn()
	})
	ch := make(chan Result[T], 1)
	go mapResult(resultCh, ch)
	return ch
}

func (g *Group[T]) Forget(key string) {
	g.underlying.Forget(key)
}
