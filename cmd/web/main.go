package main

import (
	"database/sql"
	"flag"
	"github.com/felipedavid/not_pastebin/pkg/models/mysql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *mysql.SnippetModel
}

func main() {
	// Parse command line flags
	addr := flag.String("addr", "127.0.0.1:4000", "HTTP network address")
	// parseTime is set so that the driver convert the database's types TIME and DATE
	// to Go's time.Time
	dsn := flag.String("dsn", "web:123@/not_pastebin?parseTime=true", "MySQL data source name")
	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create a database connection pool
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Establishing dependencies for the handlers
	app := application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &mysql.SnippetModel{DB: db},
	}

	// Setting up an HTTP server that uses our app's error logger
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.errorLog,
		Handler:  app.routes(),
	}

	app.infoLog.Printf("Starting web server at %s\n", *addr)
	err = srv.ListenAndServe()
	app.errorLog.Fatal(err)
}

// Create a database connection pool and check if it's possible to create
// connections
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// The 'sql.Open' don't create any connections, it just creates a pool
	// for future use. So we need to check if we can create a connection
	// to the database by calling Ping.
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}
