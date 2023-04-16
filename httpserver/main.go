package main

import (
	"net/http"
)

func main() {
	err := http.ListenAndServe(":8000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!\n"))
	}))
	if err != nil {
		panic(err)
	}
}
