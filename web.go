package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/yvasiyarov/gorelic"
)

var agent *gorelic.Agent

func init() {
	agent = gorelic.NewAgent()

	if key := os.Getenv("NEW_RELIC_LICENSE_KEY"); key != "" {
		agent.NewrelicLicense = key
		agent.Run()
	}
}

func main() {
	http.HandleFunc("/", agent.WrapHTTPHandlerFunc(handleIndex))

	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleIndex(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "hello, world")
}
