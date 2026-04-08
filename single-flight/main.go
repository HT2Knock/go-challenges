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
	for i := 0; i < 5; i++ {
		go func(i int) {
			v, err, _ := group.Do("the_key", func() (any, error) {
				return expensivOperation("the_key")
			})

			if err == nil {
				fmt.Printf("Goroutine %d got result: %v\n", i, v)
			}
		}(i)
	}

	time.Sleep(3 * time.Second)
}
