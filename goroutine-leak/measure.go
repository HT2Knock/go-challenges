package main

import (
	"fmt"
	"runtime"
	"time"
)

func measureGoroutineMemory() {
	runtime.GC()

	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	const numGoroutines = 10000
	done := make(chan struct{})

	for range numGoroutines {
		go func() {
			time.Sleep(10 * time.Second)
			done <- struct{}{}
		}()
	}

	runtime.GC()
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	memPerGoroutine := (m2.StackInuse - m1.StackInuse) / numGoroutines
	fmt.Printf("Memory per goroutine: %d bytes\n", memPerGoroutine)
	fmt.Printf("Total memory for %d goroutines: %d MB\n",
		numGoroutines, (m2.StackInuse-m1.StackInuse)/(1024*1024))
}
