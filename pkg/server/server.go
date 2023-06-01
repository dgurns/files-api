package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Run() error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	username := os.Getenv("BASIC_AUTH_USERNAME")
	pwd := os.Getenv("BASIC_AUTH_PASSWORD")
	if username == "" || pwd == "" {
		fmt.Println("BASIC_AUTH_USERNAME and BASIC_AUTH_PASSWORD must be set")
		os.Exit(1)
	}
	r.Use(middleware.BasicAuth("user", map[string]string{
		username: pwd,
	}))

	r.Post("/files", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/files/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Delete("/files/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/files/search", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	fmt.Println("files-api server listening on port 8080")
	return http.ListenAndServe(":8080", r)
}
