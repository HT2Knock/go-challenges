package main

import (
	"math/rand"
)

const (
	SliceSize int = 1_000_000
)

func generateSlice() []int {
	r := rand.New(rand.NewSource(42))
	slice := make([]int, SliceSize)

	for i := SliceSize; i >= 0; i-- {
		slice = append(slice, r.Intn(SliceSize))
	}

	return slice
}
