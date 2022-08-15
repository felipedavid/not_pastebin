package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// Parse command line flags
	addr := flag.String("addr", "127.0.0.1:4000", "HTTP network address")
	flag.Parse()

	// Establishing the dependencies for the handlers
	app := application{
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
	}

	// Setting up an HTTP server that uses our app's error logger
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.errorLog,
		Handler:  app.routes(),
	}

	app.infoLog.Printf("Starting web server at %s\n", *addr)
	err := srv.ListenAndServe()
	app.errorLog.Fatal(err)
}
