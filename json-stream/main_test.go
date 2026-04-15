package main

import "testing"

func BenchmarkConvertOld(b *testing.B) {
	for b.Loop() {
		if err := convert("./data/large-file.json"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConvertV2(b *testing.B) {
	for b.Loop() {
		if err := convertV2("./data/large-file.json"); err != nil {
			b.Fatal(err)
		}
	}
}
