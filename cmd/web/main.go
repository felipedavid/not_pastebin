package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	addr := *flag.String("addr", "127.0.0.1:4000", "Server address")
	flag.Parse()

	// routes just match incoming requests with registered mappings and redirect
	// them to the right handler
	mux := http.NewServeMux() //Creating a new router

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home) // Adding a map between path "/" and the home handler
	mux.HandleFunc("/snippet/view/", viewSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	log.Printf("Starting server on %s\n", addr)
	// Creates a web server that listens to a TCP port 4000 of every interface.
	err := http.ListenAndServe(addr, mux)
	fmt.Errorf(err.Error())
}
