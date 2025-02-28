package main

import (
	"net/http"
	
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	loadNextTaskID() // load next possible task ID

	// start router

	r := chi.NewRouter()

	r.Use(middleware.Logger) // logger middleware

	r.Use(JSONMiddleWare) // json handler

	r.Get("/", func(w http.ResponseWriter,  r *http.Request){ // http request
		w.Write([]byte("Hello"))
	})

	r.Get("/tasks", getTasks)

	r.Post("/tasks", createTask)

	r.Patch("/tasks/{taskID}", updateTask) // handle updating task

	r.Delete("/tasks/{taskID}", deleteTask) // delete task

	http.ListenAndServe("localhost:9000", r)
}

func JSONMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application-json") // to handle json files
		next.ServeHTTP(w, r)
	})
}
