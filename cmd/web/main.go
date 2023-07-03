package main

import (
	"log"
	"net/http"
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

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/snippet/view", viewSnippetHandler)
	mux.HandleFunc("/snippet/create", createSnippetHandler)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe("localhost:8080", mux)
	log.Fatal(err)
}
