package server

import (
	"net/http"
)

func ApplicationMiddlewareResponse(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	}
}

func HandleRoutesNotFound(mux *http.ServeMux) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, pattern := mux.Handler(r)
		if pattern == "" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"status":"error","message":"Route not found"}`))
			return
		}
		mux.ServeHTTP(w, r)
	}
}
