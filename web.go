package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", handleIndex)

	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func handleIndex(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "hello, world")
}
