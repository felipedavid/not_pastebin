package main

import (
	"net/http"
	"log"
    "flag"
    "os"
    "database/sql"
    _ "github.com/lib/pq"
)

type app struct {
    infoLogger *log.Logger
    errLogger *log.Logger
}

func main() {
    addr := *flag.String("addr", "127.0.0.1:8000", "HTTP network address")
    dsn := *flag.String("dsn", "postgres://postgres:postgres@localhost/not_pastebin?sslmode=disable", "Domain service name")
    flag.Parse()

    errLogger := log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
    infoLogger := log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)

    db, err := openDB(dsn)
    if err != nil {
        errLogger.Fatal(err)
    }
    defer db.Close()

    a := app{
        errLogger: errLogger,
        infoLogger: infoLogger,
    }

    s := http.Server{
        Addr: addr,
        Handler: a.routes(),
        ErrorLog: errLogger,
    }

	infoLogger.Printf("Starting server on %s\n", addr)
	err = s.ListenAndServe()
	errLogger.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }

    if err = db.Ping(); err != nil {
        return nil, err
    }
    return db, nil
}
