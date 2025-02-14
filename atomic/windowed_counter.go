package atomic

import (
	"sync/atomic"
	"time"

	"golang.org/x/exp/constraints"
)

type WindowedCounter[UintT constraints.Unsigned] struct {
	counter atomic.Uint64
	running atomic.Bool
	stop    chan struct{}
}

func NewUnstartedWindowedCounter[UintT constraints.Unsigned]() *WindowedCounter[UintT] {
	return new(WindowedCounter[UintT])
}

func NewWindowedCounter[UintT constraints.Unsigned](resetInterval time.Duration) *WindowedCounter[UintT] {
	counter := NewUnstartedWindowedCounter[UintT]()
	counter.Start(resetInterval)
	return counter
}

func (c *WindowedCounter[UintT]) Start(resetInterval time.Duration) {
	if c.running.CompareAndSwap(false, true) { // If it's already running, do nothing.
		c.stop = make(chan struct{})
		go c.resetLoop(resetInterval)
	}
}

func (c *WindowedCounter[UintT]) Stop() {
	if c.running.CompareAndSwap(true, false) { // If it has already been stopped, do nothing.
		close(c.stop)
		c.stop = nil
		c.Reset()
	}
}

func (c *WindowedCounter[UintT]) Add(delta UintT) UintT {
	return UintT(c.counter.Add(uint64(delta)))
}

func (c *WindowedCounter[UintT]) Inc() UintT {
	return c.Add(1)
}

func (c *WindowedCounter[UintT]) Load() UintT {
	return UintT(c.counter.Load())
}

func (c *WindowedCounter[UintT]) Reset() {
	c.counter.Store(0)
}

func (c *WindowedCounter[UintT]) resetLoop(resetInterval time.Duration) {
	tick := time.Tick(resetInterval)
	for {
		select {
		case <-c.stop:
			return
		case <-tick:
			c.Reset()
		}
	}
}
