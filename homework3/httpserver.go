package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"reflect"
)

var (
	listenAddr = flag.String("http", ":80", "http listen address")
)

type HandleFnc func(http.ResponseWriter, *http.Request)

var version string

func init() {
	version = os.Getenv("VERSION")
	if version == "" {
		version = "unknown"
	}
}

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", logInfo(healthz))
	srv := http.Server{
		Addr:    *listenAddr,
		Handler: mux,
	}
	
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen err: %s\n", err)
		}
	}()
	log.Print("Server is running")
	
	<-done
	
	log.Print("Server is exiting")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		log.Print("Server is exited")
		cancel()
	}()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}

func logInfo(f HandleFnc) HandleFnc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				log.Printf("[%v] caught panic: %v", r.RemoteAddr, x)
			}
		}()

		log.Printf("Request: %s, %s\n", r.RemoteAddr, r.URL.String())

		for k, v := range r.Header {
			for i := 0; i < len(v); i++ {
				w.Header().Set(k, v[i])
			}
		}

		w.Header().Set("VERSION", version)
		f(w, r)

		value := reflect.ValueOf(w).Elem()

		log.Printf("Response: status %d", value.FieldByName("status"))
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello\n")
}

