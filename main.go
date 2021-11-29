package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

// Healthz func API endpoint for backend healthchecks
func Healthz(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health is OK")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

func main() {
	log.Info("Starting Todolist API server")
	host, set := os.LookupEnv("TODOIST_HOST")
	if !set {
		host = "0.0.0.0"
	}
	port, set := os.LookupEnv("TODOIST_PORT")
	if !set {
		port = "8080"
	}

	/* mux is a simple library to build API's
	   - at first we add /heatlhz endpoint and specify GET method for it
	*/
	router := mux.NewRouter()
	router.HandleFunc("/healthz", Healthz).Methods("GET")

	// To share our routes we need a http server
	srv := &http.Server{
		Addr: fmt.Sprintf("%v:%v", host, port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}
	// TODO(vnaumov) handle graceful shutdown like: https://github.com/gorilla/mux#graceful-shutdown
	log.Fatal(srv.ListenAndServe())
}
