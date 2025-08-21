package atomic

import (
	"sync/atomic"

	"golang.org/x/exp/constraints"
)

type Integer[IntT constraints.Integer] interface {
	Add(delta IntT) (new IntT)
	And(mask IntT) (old IntT)
	CompareAndSwap(old IntT, new IntT) (swapped bool)
	Load() IntT
	Or(mask IntT) (old IntT)
	Store(val IntT)
	Swap(new IntT) (old IntT)
}

var (
	_ Integer[int32]   = (*atomic.Int32)(nil)
	_ Integer[int64]   = (*atomic.Int64)(nil)
	_ Integer[uint32]  = (*atomic.Uint32)(nil)
	_ Integer[uint64]  = (*atomic.Uint64)(nil)
	_ Integer[uintptr] = (*atomic.Uintptr)(nil)
)
