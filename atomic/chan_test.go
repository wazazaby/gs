package atomic

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestCloseSafeChanSendReceive(t *testing.T) {
	data := []string{"foo", "bar", "baz", "doe"}

	ch := MakeCloseSafeChan[string]()

	var (
		eg       errgroup.Group
		received atomic.Uint32
	)

	eg.Go(func() error {
		for {
			value, ok := ch.Receive()
			if !ok {
				return nil
			}
			require.Contains(t, data, value)
			received.Add(1)
		}
	})

	for _, s := range data {
		require.True(t, ch.Send(s))
	}
	require.NoError(t, ch.Close())

	require.NoError(t, eg.Wait())
	require.Equal(t, uint32(len(data)), received.Load())
}

func TestCloseSafeChanSendReceiveBuffered(t *testing.T) {
	data := []string{"foo", "bar", "baz", "doe"}

	ch := MakeCloseSafeChan[string](len(data))

	for _, s := range data {
		require.True(t, ch.Send(s))
	}
	require.NoError(t, ch.Close())

	var received int
	for {
		value, ok := ch.Receive()
		if !ok {
			break
		}
		require.Contains(t, data, value)
		received++
	}

	require.Equal(t, len(data), received)
}

func TestCloseSafeChanClose(t *testing.T) {
	ch := MakeCloseSafeChan[struct{}]()

	require.NoError(t, ch.Close())

	require.False(t, ch.Send(struct{}{}))

	_, ok := ch.Receive()
	require.False(t, ok)

	require.Zero(t, ch.sending.Load())
	require.Equal(t, uint32(2), ch.state.Load())
}

func TestCloseSafeChanConcurrentSendsReceivesCloses(t *testing.T) {
	const (
		nbParallelClosingGoroutines = 64
		nbParallelSendingGoroutines = 256
		nbSend                      = 128
	)

	t.Logf(
		"testing with %d closing goroutines and %d sending goroutines (%d payloads each)",
		nbParallelClosingGoroutines,
		nbParallelSendingGoroutines,
		nbSend,
	)

	require.NotPanics(t, func() {
		ch := MakeCloseSafeChan[struct{}]()

		var sent, skipped, received atomic.Uint32
		var eg errgroup.Group

		// Spinning a receive goroutine.
		eg.Go(func() error {
			for {
				_, ok := ch.Receive()
				if !ok {
					return nil
				}
				received.Add(1)
			}
		})

		t.Logf("waiting for receive goroutine to be scheduled")
		time.Sleep(time.Second)

		start := make(chan struct{})

		// Closing goroutines.
		for range nbParallelClosingGoroutines {
			eg.Go(func() error {
				<-start
				return ch.Close()
			})
		}

		// Sending goroutines.
		for range nbParallelSendingGoroutines {
			eg.Go(func() error {
				<-start
				for range nbSend {
					if ch.Send(struct{}{}) {
						sent.Add(1)
					} else {
						skipped.Add(1)
					}
				}
				return nil
			})
		}

		t.Log("waiting for sending/closing goroutines to be scheduled")
		time.Sleep(time.Second)
		t.Log("starting everything")
		close(start)

		require.NoError(t, eg.Wait())

		t.Logf("sent %d, received %d, skipped %d", sent.Load(), received.Load(), skipped.Load())

		require.Equal(t, sent.Load(), received.Load())
		require.Equal(t, uint32(nbParallelSendingGoroutines*nbSend), sent.Load()+skipped.Load())
	})
}

func TestCloseSafeChanLenCap(t *testing.T) {
	const (
		buffer = 64
		nbSend = 16
	)

	ch := MakeCloseSafeChan[struct{}](buffer)

	for range nbSend {
		require.True(t, ch.Send(struct{}{}))
	}

	require.NoError(t, ch.Close())

	require.Equal(t, buffer, ch.Cap())
	require.Equal(t, nbSend, ch.Len())
}

func TestCloseSafeChanSendIter(t *testing.T) {
	data := []string{"foo", "bar", "baz", "doe", "qux", "hello", "world", "ted"}

	ch := MakeCloseSafeChan[string]()

	var (
		eg       errgroup.Group
		received atomic.Uint32
	)

	eg.Go(func() error {
		for value := range ch.Iter() {
			require.Contains(t, data, value)
			received.Add(1)
		}
		return nil
	})

	for _, s := range data {
		require.True(t, ch.Send(s))
	}
	require.NoError(t, ch.Close())

	require.NoError(t, eg.Wait())
	require.Equal(t, uint32(len(data)), received.Load())
}

func TestCloseSafeChanSendIterBuffered(t *testing.T) {
	data := []string{"foo", "bar", "baz", "doe", "qux", "hello", "world", "ted"}

	ch := MakeCloseSafeChan[string](len(data))

	for _, s := range data {
		require.True(t, ch.Send(s))
	}
	require.NoError(t, ch.Close())

	var received int
	for value := range ch.Iter() {
		require.Contains(t, data, value)
		received++
	}

	require.Equal(t, len(data), received)
}
func BenchmarkChanSendReceiveSeq(b *testing.B) {
	b.ReportAllocs()

	ch := make(chan struct{}, 1)
	s := struct{}{}

	b.Cleanup(func() {
		close(ch)
	})

	b.ResetTimer()

	for range b.N {
		ch <- s
		<-ch
	}
}

func BenchmarkChanSendReceivePar(b *testing.B) {
	b.ReportAllocs()

	ch := make(chan struct{}, 64)
	s := struct{}{}

	b.Cleanup(func() {
		close(ch)
	})

	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			ch <- s
			<-ch
		}
	})
}

func BenchmarkChanCloseSafeSendReceiveSeq(b *testing.B) {
	b.ReportAllocs()

	ch := MakeCloseSafeChan[struct{}](1)
	s := struct{}{}

	b.Cleanup(func() {
		require.NoError(b, ch.Close())
	})

	b.ResetTimer()

	for range b.N {
		_ = ch.Send(s)
		_, _ = ch.Receive()
	}
}

func BenchmarkChanCloseSafeSendReceivePar(b *testing.B) {
	b.ReportAllocs()

	ch := MakeCloseSafeChan[struct{}](64)
	s := struct{}{}

	b.Cleanup(func() {
		require.NoError(b, ch.Close())
	})

	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_ = ch.Send(s)
			_, _ = ch.Receive()
		}
	})
}
