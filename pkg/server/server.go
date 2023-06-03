package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/dgurns/files-api/internal/db"
	"github.com/dgurns/files-api/internal/storage"
	"github.com/dgurns/files-api/pkg/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Run() error {
	// init db client
	dbConn, err := sql.Open(
		"sqlite3",
		os.Getenv("LOCAL_DB_PATH")+"/"+os.Getenv("LOCAL_DB_NAME"),
	)
	if err != nil {
		fmt.Println("Error opening database connection")
		return err
	}
	defer dbConn.Close()
	dbClient, err := db.NewSQLiteClient(dbConn)
	if err != nil {
		fmt.Println("Error creating database client")
		return err
	}

	// init storage client
	storagePath := os.Getenv("LOCAL_FILES_PATH")
	storageClient, err := storage.NewFilesystemClient(storagePath)
	if err != nil {
		return err
	}

	r := chi.NewRouter()

	// set up middleware
	username := os.Getenv("BASIC_AUTH_USERNAME")
	pwd := os.Getenv("BASIC_AUTH_PASSWORD")
	if username == "" || pwd == "" {
		fmt.Println("BASIC_AUTH_USERNAME and BASIC_AUTH_PASSWORD must be set")
		os.Exit(1)
	}
	r.Use(middleware.BasicAuth("user", map[string]string{
		username: pwd,
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// handle routes
	h := handler.New(dbClient, storageClient)
	r.Post("/files/upload", h.UploadFile)
	r.Get("/files/{id}", h.GetFile)
	r.Delete("/files/{id}", h.DeleteFile)
	r.Get("/files/search", h.SearchFiles)

	fmt.Println("files-api server listening on port 8080")
	return http.ListenAndServe(":8080", r)
}
