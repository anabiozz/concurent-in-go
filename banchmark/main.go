package main

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

var memoryAccess sync.Mutex
var value int

func lock() {
	go func() {
		memoryAccess.Lock()
		value++
		memoryAccess.Unlock()
	}()

	memoryAccess.Lock()
	if value == 0 {
		fmt.Printf("the value is %v.\n", value)
	} else {
		fmt.Printf("the value is %v.\n", value)
	}
	memoryAccess.Unlock()
}

func memConsumed() uint64 {
	runtime.GC()
	var s runtime.MemStats
	runtime.ReadMemStats(&s)
	return s.Sys
}

func BenchmarkContextSwitch(b *testing.B) {
	var wg sync.WaitGroup
	begin := make(chan struct{})
	c := make(chan struct{})

	var token struct{}
	sender := func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			c <- token
		}
	}
	reciever := func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			<-c
		}
	}
	wg.Add(2)
	go sender()
	go reciever()
	b.StartTimer()
	close(begin)
	wg.Wait()
}

func main() {
	var c <-chan interface{}
	var wg sync.WaitGroup
	noop := func() { wg.Done(); <-c }

	const numGorutines = 1e4
	wg.Add(numGorutines)
	before := memConsumed()
	for i := numGorutines; i > 0; i-- {
		go noop()
	}
	wg.Wait()
	after := memConsumed()
	fmt.Printf("\n%.3fkb\n", float64(after-before)/numGorutines/1000)
}
