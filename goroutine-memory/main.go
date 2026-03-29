package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func measureWithMemStats() {
	runtime.GC()

	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	const numGoroutines = 10000

	for range numGoroutines {
		go func() {
			time.Sleep(10 * time.Second)
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

func measureWithPprof() {
	const numGoroutines = 100

	for range numGoroutines {
		go func() {
			time.Sleep(10 * time.Second)
		}()
	}

	f, err := os.Create("goroutine.prof")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create profile file: %v\n", err)
		return
	}
	defer f.Close()

	if err := pprof.WriteHeapProfile(f); err != nil {
		fmt.Fprintf(os.Stderr, "Could not write heap profile: %v\n", err)
		return
	}

	fmt.Println("Heap profile written to goroutine.prof")
	fmt.Println("Analyze with: go tool pprof goroutine.prof")
}

func main() {
	fmt.Println("=== Measuring goroutine memory with MemStats ===")
	measureWithMemStats()

	fmt.Println("\n=== Measuring with pprof heap profile ===")
	measureWithPprof()
}
