package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Just a neat way to do dependency injection. If a handler or helper function
// needs some kind of dependency we just add the dependency into the app struct,
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

	// Instantiating application's dependencies
	a := &app{
		errLogger:  errLog,
		infoLogger: infoLog,
	}

	// Creating a new server and listening in 'addr'
	server := &http.Server{
		Addr:     addr,
		Handler:  a.routes(),
		ErrorLog: errLog,
	}
	infoLog.Printf("Starting server on %s\n", addr)
	err := server.ListenAndServe()
	errLog.Fatal(err.Error())
}
