package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/felipedavid/not_pastebin/internal/models"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type app struct {
	infoLogger     *log.Logger
	errLogger      *log.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
	env            string
}

func main() {
	addr := *flag.String("addr", ":4000", "server listen address")
	dsn := *flag.String("dsn",
		"postgres://postgres:secret@localhost/not_pastebin?sslmode=disable",
		"Data source name")
	env := *flag.String("env", "development", "Environment (development|production)")
	flag.Parse()

	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLogger := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(dsn)
	if err != nil {
		errLogger.Fatal(err)
	}

	templateCache, err := newTemplateCache()

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true // Ensures that the web browser will only send the cookie thought a TSL connection

	a := app{
		infoLogger:     infoLogger,
		errLogger:      errLogger,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		sessionManager: sessionManager,
		env:            env,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	s := http.Server{
		Addr:      addr,
		ErrorLog:  errLogger,
		Handler:   a.routes(),
		TLSConfig: tlsConfig,

		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLogger.Printf("Starting up server at %s\n", addr)
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
