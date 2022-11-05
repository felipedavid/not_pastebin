package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// routes just match incoming requests with registered mappings and redirect
	// them to the right handler
	mux := http.NewServeMux() //Creating a new router

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home) // Adding a map between path "/" and the home handler
	mux.HandleFunc("/snippet/view/", viewSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	log.Println("Starting server on 127.0.0.1:4000")
	// Creates a web server that listens to a TCP port 4000 of every interface.
	err := http.ListenAndServe(":4000", mux)
	fmt.Errorf(err.Error())
}
