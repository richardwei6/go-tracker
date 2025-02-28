package main

import (
	"net/http"
	
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger) // logger middleware

	r.Use(JSONMiddleWare) // json handler

	r.Get("/", func(w http.ResponseWriter,  r *http.Request){ // https request
		w.Write([]byte("Hello test"))
	})

	r.Get("/tasks", getTasks)

	r.Post("/tasks", createTask)

	http.ListenAndServe("localhost:9000", r)
}

func JSONMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application-json") // to handle json files
		next.ServeHTTP(w, r)
	})
}