package main

import "fmt"

//go:generate go run golang.org/x/tools/cmd/stringer@latest -type=Level

type Level int

const (
	Info Level = iota
	Error
	Fatal
)

func main() {
	fmt.Printf("%v: Hello world!\n", Info)
}
