package main

import (
	"log"
	"io"
	"flag"
	"reflect"
	"os"
	"net/http"
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
	
	http.HandleFunc("/healthz", logInfo(healthz))
	err := http.ListenAndServe(*listenAddr, nil)
    if err != nil {
        log.Panicln("ListenAndServe:", err)
    }
}

func logInfo(f HandleFnc) HandleFnc {
	return func (w http.ResponseWriter, r *http.Request) {
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