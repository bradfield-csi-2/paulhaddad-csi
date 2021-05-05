package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	done := make(chan struct{}, 1)
	wg.Add(1)
	go func() {
		fmt.Println("performing initialization...")
		<-done
		wg.Done()
	}()

	done <- struct{}{}
	fmt.Println("initialization done, continuing with rest of program")
	wg.Wait()
}
