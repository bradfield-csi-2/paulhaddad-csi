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

func (c *counter) getNext(reqCh chan bool, respCh chan uint64) {
	reqCh <- true
	val := <-respCh
	c.val = val
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
			inc := <-reqCh
			if !inc {
				break
			}

			count++
			respCh <- count
		}

		wg.Done()
	}()

	for i := 0; i < 1000; i++ {
		counter.getNext(reqCh, respCh)
	}

	reqCh <- false

	wg.Wait()

	fmt.Printf("The final value of the counter is %d\n", counter.val)
}
