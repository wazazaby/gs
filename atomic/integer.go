package atomic

import (
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
