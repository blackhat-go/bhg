package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type logger struct {
	Inner http.Handler
}

func (l *logger) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("start %s\n", time.Now().String())
	l.Inner.ServeHTTP(w, req)
	log.Printf("finish %s\n", time.Now().String())
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Hello\n")
}

func main() {
	f := http.HandlerFunc(hello)
	l := logger{Inner: f}
	http.ListenAndServe(":8000", &l)
}
