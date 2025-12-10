package atomic

import (
	"sync/atomic"

	"golang.org/x/exp/constraints"
)

type Pointer[T any] = Val[*T]

type Val[T any] interface {
	Load() T
	Store(val T)
	Swap(new T) (old T)
	CompareAndSwap(old T, new T) (swapped bool)
}

type Integer[T constraints.Integer] interface {
	Val[T]
	Add(delta T) (new T)
	And(mask T) (old T)
	Or(mask T) (old T)
}

var (
	_ Integer[int32]   = (*atomic.Int32)(nil)
	_ Integer[int64]   = (*atomic.Int64)(nil)
	_ Integer[uint32]  = (*atomic.Uint32)(nil)
	_ Integer[uint64]  = (*atomic.Uint64)(nil)
	_ Integer[uintptr] = (*atomic.Uintptr)(nil)
	_ Val[bool]        = (*atomic.Bool)(nil)
	_ Pointer[any]     = (*atomic.Pointer[any])(nil)
)
