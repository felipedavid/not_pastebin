package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/felipedavid/not_pastebin/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

// Just a neat way to do dependency injection. If a handler or helper function
// needs some kind of dependency we just add the dependency into the app struct,
// and then we make the procedure a method of the struct
type app struct {
	errLogger     *log.Logger
	infoLogger    *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	// Parsing command line flags
	addr := *flag.String("addr", "127.0.0.1:4000", "Server address")
	dsn := *flag.String("dsn", "web:pass@tcp(127.0.0.1:3306)/not_pastebin?parseTime=true", "Database service name")
	flag.Parse()

	// Creating application's loggers
	errLog := log.New(os.Stderr, "ERROR\t", log.Lshortfile|log.Ldate|log.Ltime)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Connecting to the database
	db, err := openDatabase(dsn)
	if err != nil {
		errLog.Fatal(err)
	}
	defer db.Close()

	// Creating a template cache, so we don't need to read the template files
	// from disk and parse them for every request
	templateCache, err := newTemplateCache()
	if err != nil {
		errLog.Fatal(err)
	}

	// Instantiating application's dependencies
	a := &app{
		errLogger:     errLog,
		infoLogger:    infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// Creating a new server and listening in 'addr'
	server := &http.Server{
		Addr:     addr,
		Handler:  a.routes(),
		ErrorLog: errLog,
	}
	infoLog.Printf("Starting server on %s\n", addr)
	err = server.ListenAndServe()
	errLog.Fatal(err.Error())
}

// openDatabase creates a database connection pool and then checks if the database is reachable
func openDatabase(dsn string) (*sql.DB, error) {
	// Create a database connection pool
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Check connection to the database
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
