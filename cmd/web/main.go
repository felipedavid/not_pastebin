package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func methodNotAllowedResponse(w http.ResponseWriter, r *http.Request, methods ...string) {
	allowed := strings.Join(methods, ", ")
	w.Header().Set("Allow", allowed)
	errorResponse(w, r, http.StatusMethodNotAllowed)
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	w.Write([]byte(http.StatusText(status)))
}

func getQueryInt(r *http.Request, paramName string) (int, error) {
	return strconv.Atoi(r.URL.Query().Get(paramName))
}

func parseEnvVariable(variable, defaultVal string) string {
	value := os.Getenv(variable)
	if value == "" {
		value = defaultVal
	}
	return value
}

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	addr := parseEnvVariable("ADDR", "localhost:8080")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	server := &http.Server{
		Addr:     addr,
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}

	infoLog.Printf("Starting server on %s\n", addr)
	err := server.ListenAndServe()
	errorLog.Fatal(err)
}
