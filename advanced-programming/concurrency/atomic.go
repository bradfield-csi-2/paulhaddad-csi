package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type counterService interface {
	// Returns values in ascending order; it should be safe to call
	// getNext() concurrently without any additional synchronization.
	getNext() uint64
}

type counter struct {
	val uint64
}

func (c *counter) getNext() uint64 {
	return atomic.AddUint64(&c.val, 1)
}

func main() {
	var wg sync.WaitGroup

	wg.Add(2)
	counter := counter{0}

	go func() {
		for i := 0; i < 100; i++ {
			counter.getNext()
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 100; i++ {
			counter.getNext()
		}
		wg.Done()
	}()

	wg.Wait()

	fmt.Printf("The final value of counters is %d\n", counter.val)
}
