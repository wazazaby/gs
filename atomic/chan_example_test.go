package atomic_test

import (
	"fmt"

	"github.com/wazazaby/gs/atomic"
)

func Example() {
	ch := atomic.MakeCloseSafeChan[string](8)

	fmt.Println("Cap:", ch.Cap())
	fmt.Println("Len:", ch.Len())

	sent := ch.Send("foo")
	fmt.Println("Sent foo ->", sent)

	sent = ch.Send("bar")
	fmt.Println("Sent bar ->", sent)

	fmt.Println("Len:", ch.Len())

	foo, ok := ch.Receive()
	fmt.Printf("Got %s - chan open -> %t\n", foo, ok)

	bar, ok := ch.Receive()
	fmt.Printf("Got %s - chan open -> %t\n", bar, ok)

	sent = ch.Send("john")
	fmt.Println("Sent john ->", sent)

	sent = ch.Send("doe")
	fmt.Println("Sent doe ->", sent)

	fmt.Println(ch.Close())

	fmt.Println("Len:", ch.Len())

	for value := range ch.Iter() {
		fmt.Printf("Got %s in range\n", value)
	}

	sent = ch.Send("baz")
	fmt.Println("Sent baz ->", sent)

	baz, ok := ch.Receive()
	fmt.Printf("Is empty baz -> %t - chan open -> %t\n", baz == "", ok)

	// Output:
	// Cap: 8
	// Len: 0
	// Sent foo -> true
	// Sent bar -> true
	// Len: 2
	// Got foo - chan open -> true
	// Got bar - chan open -> true
	// Sent john -> true
	// Sent doe -> true
	// <nil>
	// Len: 2
	// Got john in range
	// Got doe in range
	// Sent baz -> false
	// Is empty baz -> true - chan open -> false
}
