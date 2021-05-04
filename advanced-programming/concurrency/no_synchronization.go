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

	fmt.Println() // this creates a system call that illustrates the problem with this code

	c.val = val
	return c.val
}

func main() {
	var wg sync.WaitGroup

	wg.Add(2)
	counter := counter{0}

	go func() {
		for i := 0; i < 100; i++ {
			fmt.Println("Function 1", counter.getNext())
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 100; i++ {
			fmt.Println("Function 2", counter.getNext())
		}
		wg.Done()
	}()

	wg.Wait()
}
