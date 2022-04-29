package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

var (
	delayRequestsDuration    = 1 * time.Millisecond
	gracefulShutdownDuration = 1 * time.Millisecond
	delayWhenReady           = 15 * time.Second
)

func hello(w http.ResponseWriter, r *http.Request) {
	print("Hello request\n")
	time.Sleep(delayRequestsDuration)
	w.WriteHeader(200)
	w.Write([]byte("hello\n"))
}

func healthy(w http.ResponseWriter, r *http.Request) {
	print("Health check request\n")
	w.WriteHeader(200)
	w.Write([]byte("healthy\n"))
}

func ready(w http.ResponseWriter, r *http.Request) {
	print("Ready check request\n")
	w.WriteHeader(200)
	w.Write([]byte("ready\n"))
}

func main() {

	if drd := os.Getenv("DELAY_REQUESTS_DURATION"); drd != "" {
		delayRequestsDuration, _ = time.ParseDuration(drd)
	}

	if gsd := os.Getenv("GRACEFUL_SHUTDOWN_DURATION"); gsd != "" {
		gracefulShutdownDuration, _ = time.ParseDuration(gsd)
	}

	if dwr := os.Getenv("DELAY_UNTIL_READY"); dwr != "" {
		delayWhenReady, _ = time.ParseDuration(dwr)
	}

	router := mux.NewRouter()
	router.HandleFunc("/healthy", healthy).Methods("GET")

	time.AfterFunc(delayWhenReady, func() {
		router.HandleFunc("/ready", ready).Methods("GET")
		router.HandleFunc("/hello", hello).Methods("GET")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownDuration)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
