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

var counter int

func (m *mutex) Lock() {
	for {
		if atomic.CompareAndSwapInt64(&m.val, 0, 1) {
			return
		}
	}
}

func (m *mutex) Unlock() {
	swapped := atomic.CompareAndSwapInt64(&m.val, 1, 0)
	if swapped == false {
		panic("tried to unlocked an unlocked mutex")
	}
}

func acqAndRelease(i int, m *mutex, wg *sync.WaitGroup) {
	// Acquire mutex
	m.Lock()
	fmt.Printf("Goroutine %d: The mutex at address %p is acquired\n", i, &m)

	// do work
	time.Sleep(1 * time.Second)
	counter++

	// Release mutex
	m.Unlock()
	fmt.Printf("Goroutine %d: The mutex is unlocked\n", i)

	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	var mutex mutex

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go acqAndRelease(i, &mutex, &wg)
	}

	wg.Wait()

	fmt.Printf("The value of the counter is %d\n", counter)
}
