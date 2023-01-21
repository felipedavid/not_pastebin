package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/felipedavid/not_pastebin/internal/models"
	_ "github.com/lib/pq"
)

type app struct {
	debugMode      bool
	infoLogger     *log.Logger
	errLogger      *log.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", "127.0.0.1:8000", "HTTP network address")
	dsn := flag.String("dsn",
		"postgres://postgres:postgres@localhost/not_pastebin?sslmode=disable",
		"Domain service name")
	debug := flag.Bool("debug", false, "Debug mode")
	flag.Parse()

	errLogger := log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger := log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)

	// Setup the database
	db, err := openDB(*dsn)
	if err != nil {
		errLogger.Fatal(err)
	}
	defer db.Close()

	snippetModel, err := models.NewSnippetModel(db)
	if err != nil {
		errLogger.Fatal(err)
	}

	tc, err := newTemplateCache()
	if err != nil {
		errLogger.Fatal(err)
	}

	// Setup sessions
	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	a := app{
		debugMode:      *debug,
		errLogger:      errLogger,
		infoLogger:     infoLogger,
		snippets:       snippetModel,
		templateCache:  tc,
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	s := http.Server{
		Addr:         *addr,
		Handler:      a.routes(),
		ErrorLog:     errLogger,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLogger.Printf("Starting server on %s\n", *addr)
	err = s.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
