package main

import (
    "fmt"
	"log"
	"net/http"
	"github.com/nelkinda/health-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	fmt.Printf("example app %s, commit %s, built at %s by %s", version, commit, date, builtBy)

	h := health.New(health.Health{})
	http.HandleFunc("/health", h.Handler)
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}