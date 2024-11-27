package atomic

import (
	"io"
	"iter"
	"sync/atomic"
)

const (
	stateOpen uint32 = iota
	stateClosing
	stateClosed
)

// CloseSafeChan is a wrapper around a channel, bringing thread-safety on
// concurrently executing send and close operations.
//
// Any goroutine can safely close a [CloseSafeChan] instance without worrying
// if other goroutines are currently sending values.
//
// This is guaranteed by the implementation's inner-state, using atomic
// compare-and-swap operations to transition from the channel open state to the
// closing, and then closed state.
//
// During this transition, sending goroutines won't be able to send data on the
// channel, and will be notified by [CloseSafeChan.Send] returning false.
type CloseSafeChan[T any] struct {
	// ch is the underlying channel.
	ch chan T
	// state stores the current state of a [CloseSafeChan] instance :
	// 	- 0 means it's currently open
	// 	- 1 means it's closing (a goroutine has called [CloseSafeChan.Close])
	// 	- 2 means it's closed ([CloseSafeChan.Close] is done)
	state atomic.Uint32
	// sending stores the number of goroutines that are currently trying to
	// send data on the channel.
	//
	// It works as a gauge and will be incremented / decremented by each
	// [CloseSafeChan.Send] call.
	sending atomic.Int32
}

// MakeCloseSafeChan initializes a new [CloseSafeChan] instance.
//
// The size argument works identically to the size argument when creating
// a standard channel using [make].
func MakeCloseSafeChan[T any](size ...int) *CloseSafeChan[T] {
	var s int
	if len(size) > 0 {
		s = size[0]
	}
	return &CloseSafeChan[T]{
		ch: make(chan T, s),
	}
}

// Close closes the [CloseSafeChan] instance.
//
// It is safe to call it from concurrently running goroutines, as only the
// first caller will be able to close the underlying channel.
//
// It never returns an error.
func (c *CloseSafeChan[T]) Close() error {
	// 1. Try to transition to a closing state (1).
	// If it succeeds, it means that the calling goroutine is the first to call
	// [CloseSafeChan.Close] and now owns the closing process (this is atomic
	// of course). If it fails, it means another goroutine has called
	// [CloseSafeChan.Close] beforehand, so we can simply let it handle the
	// closing operation and return.
	if !c.state.CompareAndSwap(stateOpen, stateClosing) {
		return nil
	}
	// 2. Wait for all sending goroutines to finish their sending operations.
	// Further operations won't be registered, as sending is allowed only if the
	// state of the [CloseSafeChan] is open (0) and we previously transitioned
	// to the closing state (1).
	for c.sending.Load() != 0 {
	}
	// 3. Try to transition to the closed state (2).
	// If it succeeds, it means that the goroutine own the closing operation
	// and can safely close the underlying channel without worrying of
	// duplicate [close] calls. This is because we ensure that :
	// 	- a single goroutine executed the transition to the closing state (1)
	// 	- we drained all sending operations and won't register them any further
	if c.state.CompareAndSwap(stateClosing, stateClosed) {
		close(c.ch)
	}
	return nil
}

// Send sends a value to the [CloseSafeChan] instance.
//
// It returns true if the [CloseSafeChan] instance is still open and the send
// operation succeeded, false if the [CloseSafeChan] instance is closing or is
// closed.
func (c *CloseSafeChan[T]) Send(value T) bool {
	c.sending.Add(1)
	defer c.sending.Add(-1)
	if c.state.Load() == stateOpen {
		c.ch <- value
		return true
	}
	return false
}

// Receive receives a value from the [CloseSafeChan] instance.
//
// It returns the value, and a bool indicating if the [CloseSafeChan] instance
// is open (true) or closed (false).
//
// It works similarly to pulling from a goroutine using the arrow operator.
//
//	ch := make(chan struct{})
//	v, ok := <-ch
func (c *CloseSafeChan[T]) Receive() (T, bool) {
	value, ok := <-c.ch
	return value, ok
}

// Iter returns an iterator ranging on all the values buffered or sent in the
// underlying channel.
//
// It works exactly like ranging over a channel :
//   - The iteration will stop once the [CloseSafeChan] instance is closed and
//     there are no more values to be received
//   - You can skip values using the continue statement
//   - You can break out of the iteration using the break statement
func (c *CloseSafeChan[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for value := range c.ch {
			if !yield(value) {
				return
			}
		}
	}
}

// Len returns the number of elements queued (unread) in the [CloseSafeChan]
// instance's buffer.
//
// It simply calls [len] on the underlying channel.
func (c *CloseSafeChan[T]) Len() int {
	return len(c.ch)
}

// Cap returns the [CloseSafeChan] instance's buffer capacity, in units of
// elements.
//
// It simply calls [cap] on the underlying channel.
func (c *CloseSafeChan[T]) Cap() int {
	return cap(c.ch)
}

var (
	_ io.Closer = &CloseSafeChan[any]{}
)
