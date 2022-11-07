package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Just a neat way to do dependecy injection. If a handler or helper function
// needs some kind of dependecy we just add the dependacy into the app struct
// and then we make the procedure a method of the struct
type app struct {
	errLogger  *log.Logger
	infoLogger *log.Logger
}

func main() {
	// Parsing command line flags
	addr := *flag.String("addr", "127.0.0.1:4000", "Server address")
	flag.Parse()

	// Creating application's loggers
	errLog := log.New(os.Stderr, "ERROR\t", log.Lshortfile|log.Ldate|log.Ltime)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	a := &app{
		errLogger:  errLog,
		infoLogger: infoLog,
	}

	// Creating a router and setting up routes
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view/", a.viewSnippet)
	mux.HandleFunc("/snippet/create", a.createSnippet)

	// Creating a new server and listening in 'addr'
	server := &http.Server{
		Addr:     addr,
		Handler:  mux,
		ErrorLog: errLog,
	}
	infoLog.Printf("Starting server on %s\n", addr)
	err := server.ListenAndServe()
	errLog.Fatal(err.Error())
}
