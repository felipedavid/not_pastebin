package main

import (
	"net/http"
	"log"
)

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", view)
	mux.HandleFunc("/snippet/create", create)

	log.Println("Starting server on 127.0.0.1:8080")
	err := http.ListenAndServe("127.0.0.1:8080", mux)
	log.Fatal(err)
}
