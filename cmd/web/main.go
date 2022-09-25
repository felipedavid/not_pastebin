package main

import (
	"database/sql"
	"flag"
	"github.com/felipedavid/not_pastebin/internal/data"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type app struct {
	infoLogger *log.Logger
	errLogger  *log.Logger
	snippets   *data.SnippetModel
}

func main() {
	addr := *flag.String("addr", ":4000", "server listen address")
	dsn := *flag.String("dsn", "postgres://root:secret@localhost/not_pastebin?sslmode=disable", "Data source name")
	flag.Parse()

	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLogger := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(dsn)
	if err != nil {
		errLogger.Fatal(err)
	}

	a := app{
		infoLogger: infoLogger,
		errLogger:  errLogger,
		snippets:   &data.SnippetModel{DB: db},
	}

	s := http.Server{
		Addr:     addr,
		ErrorLog: errLogger,
		Handler:  a.routes(),
	}

	infoLogger.Printf("Starting up server at %s\n", addr)
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
