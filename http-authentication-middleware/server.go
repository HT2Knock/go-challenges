package main

import (
	"fmt"
	"net/http"
)

func NewServer() http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux)

	return mux
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello!")
}

func addRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/hello", helloHandler)
}
