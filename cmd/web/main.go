package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type app struct {
	infoLogger *log.Logger
	errLogger  *log.Logger
}

func main() {
	addr := *flag.String("addr", ":4000", "server listen address")
	flag.Parse()

	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLogger := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	a := app{
		infoLogger: infoLogger,
		errLogger:  errLogger,
	}

	s := http.Server{
		Addr:     addr,
		ErrorLog: errLogger,
		Handler:  a.routes(),
	}

	infoLogger.Printf("Starting up server at %s\n", addr)
	err := s.ListenAndServe()
	errLogger.Fatal(err)
}
