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

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/felipedavid/not_pastebin/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

// application is the state that will be shared between all handlers
type application struct {
	infoLog, errorLog *log.Logger
	snippets          *models.SnippetModel
	users             *models.UsersModel
	templateCache     map[string]*template.Template
	sessionManager    *scs.SessionManager
}

func main() {
	// Parse command line flags
	addr := flag.String("addr", "127.0.0.1:4000", "HTTP network address")
	dsn := flag.String("dsn", "web:123@/not_pastebin?parseTime=true", "MySQL data source name")
	flag.Parse()

	// Creating different loggers to make it easier to redirect the application
	// output to different files based on unix's default streams
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	// Make sure the client only sends cookies over HTTPS
	sessionManager.Cookie.Secure = true

	app := &application{
		infoLog:        infoLog,
		errorLog:       errorLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UsersModel{DB: db},
		templateCache:  templateCache,
		sessionManager: sessionManager,
	}

	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errorLog,
		TLSConfig: &tls.Config{
			CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		},

		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %v\n", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

// openDB just create a database connection pool and ping it to check if we can
// establish a connection
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
