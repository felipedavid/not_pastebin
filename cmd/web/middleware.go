package main

import "net/http"

func (a *app) logRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		re := newResponseWriter(w)
		next.ServeHTTP(re, r)

		a.infoLogger.Printf("%s \"%s %s\" -> %d %s\n", r.RemoteAddr, r.Method, r.URL.Path, re.status, http.StatusText(re.status))
	}
	return http.HandlerFunc(fn)
}
