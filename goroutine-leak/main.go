package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
)

func gen() <-chan int {
	ch := make(chan int)

	go func() {
		for i := 0; ; i++ {
			ch <- i
		}
	}()

	return ch
}

func main() {
	// We could solve it using ctx.Done() signal or create a done channel

	for i := range gen() {
		if i == 5 {
			break
		}
	}

	// Log the number of go routines
	fmt.Printf("Number of Go routines: %v\n", runtime.NumGoroutine())

	// Lookup through the profilling. It's collect CPU traces, memory allocation traces, and mostly go-routine stack traces.
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
}
