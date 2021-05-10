package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type mutex struct {
	val int64
}

func (m *mutex) Lock() {
	for {
		mutexVal := atomic.LoadInt64(&m.val)
		if mutexVal == 0 {
			atomic.StoreInt64(&m.val, 1)
			return
		}
	}
}

func (m *mutex) Unlock() {
	atomic.StoreInt64(&m.val, 0)
}

func main() {
	var counter int
	var mutex mutex
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		// Acquire mutex
		mutex.Lock()

		// do work
		fmt.Printf("Goroutine 2: The mutex at address %p is acquired\n", &mutex)
		time.Sleep(3 * time.Second)
		counter++

		// Release mutex
		mutex.Unlock()
		fmt.Println("Goroutine 2: The mutex is unlocked")
		wg.Done()
	}()

	// Acquire mutex
	mutex.Lock()

	// do work
	fmt.Printf("Goroutine 1: The mutex at address %p is acquired\n", &mutex)
	time.Sleep(3 * time.Second)
	counter++

	// Release mutex
	mutex.Unlock()
	fmt.Println("Goroutine 1: The mutex is unlocked")

	wg.Wait()

	fmt.Printf("The value of the counter is %d\n", counter)
}
