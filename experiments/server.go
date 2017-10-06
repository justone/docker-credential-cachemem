package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	l, err := net.Listen("unix", "/home/nate/listen.sock")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	http.HandleFunc("/", handler)
	if err := http.Serve(l, handlers.LoggingHandler(os.Stdout, http.DefaultServeMux)); err != nil {
		log.Fatal(err)
	}
}
