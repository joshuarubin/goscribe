package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/yvasiyarov/gorelic"
)

func init() {
	if key := os.Getenv("NEW_RELIC_LICENSE_KEY"); key != "" {
		agent := gorelic.NewAgent()
		agent.Verbose = true
		agent.NewrelicLicense = key
		agent.Run()
	}
}

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
