package main

import (
	"fmt"
	"sync"
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
	val := c.val
	val++
	c.val = val

	return val
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	counter := counter{0}

	reqCh := make(chan bool)
	respCh := make(chan uint64)

	go func() {
		var count uint64

		for {
			done := <-reqCh
			if done {
				break
			}

			counter.getNext()
			count++
			fmt.Printf("The private counter value is: %d\n", count)
			respCh <- count
		}

		wg.Done()
	}()

	for i := 0; i < 1000; i++ {

		reqCh <- false
		curVal := <-respCh
		fmt.Printf("The current value of the counter is: %d\n", curVal)
	}

	reqCh <- true

	wg.Wait()

	fmt.Printf("The final value of the counter is %d\n", counter.val)
}
