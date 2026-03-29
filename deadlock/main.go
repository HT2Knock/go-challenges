// See NOTES.md for detailed explanation
package main

import (
	"fmt"
	"sync"
)

func main() {
	ch1 := make(chan struct{})
	ch2 := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		fmt.Println("A: Trying to send ch1...")
		ch1 <- struct{}{}
		fmt.Println("A: Sent to ch1, now waiting for ch2...")
		<-ch2
	}()

	go func() {
		defer wg.Done()
		fmt.Println("B: Trying to send ch2...")
		ch2 <- struct{}{}
		fmt.Println("B: Sent to ch2, now waiting for ch1...")
		<-ch1
	}()

	wg.Wait()
}
