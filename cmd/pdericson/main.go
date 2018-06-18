package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pdericson/pdericson/pkg/ping"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/ping", ping.PingHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}
