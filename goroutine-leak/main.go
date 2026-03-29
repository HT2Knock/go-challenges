// See NOTES.md for detailed explanation
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
	for i := range gen() {
		if i == 5 {
			break
		}
	}

	fmt.Printf("Number of Go routines: %v\n", runtime.NumGoroutine())

	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
}
