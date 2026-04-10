package main

import (
	"fmt"
	"time"

	"golang.org/x/sync/singleflight"
)

var group singleflight.Group

func expensivOperation(key string) (string, error) {
	time.Sleep(2 * time.Second)

	return fmt.Sprintf("Data for %s", key), nil
}

func main() {
	for i := range 5 {
		go func(i int) {
			v, err, shared := group.Do("the_key", func() (any, error) {
				return expensivOperation("the_key")
			})

			if err == nil {
				fmt.Printf("Goroutine %d got shared = %t result: %v\n", i, shared, v)
			}
		}(i)
	}

	time.Sleep(3 * time.Second)
}
