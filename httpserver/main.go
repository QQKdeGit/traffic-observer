package main

import (
	"net/http"
)

func main() {
	err := http.ListenAndServe(":8001", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world\n"))
	}))
	if err != nil {
		panic(err)
	}
}
