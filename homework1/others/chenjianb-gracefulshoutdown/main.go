package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main(){
	c := make(chan os.Signal,2)
	defer close(c)
	signal.Notify(c,syscall.SIGHUP,syscall.SIGUSR1,syscall.SIGUSR2,syscall.SIGQUIT)
	srv := NewServer()
	go WaitReturn(c,srv)
	err := srv.ListenAndServe()
	if err!=nil{
		fmt.Printf("http server : %s",err)
	}

}

func NewServer() *http.Server{
	router := http.NewServeMux()
	router.HandleFunc("/healthz",healthFunc)
	router.HandleFunc("/",SingeFunc)
	ser := http.Server{
		Addr: ":8011",
		Handler: router,
	}
	return &ser
}

func SingeFunc(w http.ResponseWriter,r *http.Request){
	for headName,strSlice := range  r.Header {
		w.Header().Set(headName,SliceToString(strSlice))
	}
	Version := os.Getenv("VERSION")
	w.Header().Set("VERSION",Version)
	w.WriteHeader(http.StatusOK)
	fmt.Printf("ip %s ,code %d,",r.Host,http.StatusOK)
}

func healthFunc(w http.ResponseWriter,r * http.Request){
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200"))
	fmt.Printf("ip %s ,code %d, health:ok",r.Host,http.StatusOK)
}

func WaitReturn(c chan  os.Signal,srv *http.Server) {
	s :=  <-c
	switch s  {
	case syscall.SIGHUP, syscall.SIGUSR1:
		fmt.Println("server exist")
		srv.Shutdown(context.Background())
		return
	case syscall.SIGUSR2:
		fmt.Println("server ok ")
	}
}


func SliceToString(s []string) string{
	result := ""
	for n,value := range s{
		if n ==0{
			result += value
		}
		result += ", " + value
	}
	return result
}
