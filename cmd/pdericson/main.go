package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"crawshaw.io/littleboss"
	"github.com/gorilla/mux"

	"github.com/pdericson/pdericson/pkg/count"
	"github.com/pdericson/pdericson/pkg/ping"
)

var Version string

func VersionHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Fprintf(w, "%s\n", Version)
}

func main() {
	lb := littleboss.New("pdericson")
	flagHTTP := lb.Listener("http", "tcp", "127.0.0.1:8080", "address")
	lb.Run(func(ctx context.Context) {
		httpMain(ctx, flagHTTP.Listener())
	})
}

func httpMain(ctx context.Context, ln net.Listener) {
	r := mux.NewRouter()

	r.HandleFunc("/count", count.Handler).Methods("POST")
	r.HandleFunc("/count/{name}", count.Handler).Methods("GET")
	r.HandleFunc("/ping", ping.PingHandler)
	r.HandleFunc("/version", VersionHandler)

	srv := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      r,
	}

	go func() {
		if err := srv.Serve(ln); err != nil {
			if err == http.ErrServerClosed {
				return
			}
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	srv.Shutdown(ctx)
}
