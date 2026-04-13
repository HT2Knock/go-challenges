package main

import (
	"slices"
	"sort"
	"testing"
)

func BenchmarkSortSlice(b *testing.B) {
	slice := generateSlice()
	b.ReportAllocs()

	for b.Loop() {
		sort.Slice(slice, func(i, j int) bool {
			return slice[i] < slice[j]
		})
	}
}

func BenchmarkSliceSort(b *testing.B) {
	slice := generateSlice()
	b.ReportAllocs()

	for b.Loop() {
		slices.Sort(slice)
	}
}
