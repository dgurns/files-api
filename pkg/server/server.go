package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgurns/files-api/internal/db"
	"github.com/dgurns/files-api/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Run() error {
	dbClient, err := db.NewSQLiteClient("/files")
	if err != nil {
		return err
	}

	storageClient, err := storage.NewFilesystemClient()
	if err != nil {
		return err
	}

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

	r.Post("/files/upload", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(32 << 20) // 32 MB max memory.
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error retrieving file from form data", http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}

		var metadata map[string]string
		raw := r.FormValue("metadata")
		if raw != "" {
			json.Unmarshal([]byte(raw), &metadata)
		}

		id, err := dbClient.SaveFile(handler.Filename, metadata)
		if err != nil {
			http.Error(w, "Error saving file to database", http.StatusInternalServerError)
			return
		}
		err = storageClient.SaveFile(id, fileBytes)
		if err != nil {
			http.Error(w, "Error saving file to storage", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		res, err := json.Marshal(&UploadFileResponse{ID: id})
		if err != nil {
			http.Error(w, "Error marshalling response", http.StatusInternalServerError)
			return
		}
		w.Write(res)
	})

	r.Get("/files/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "Invalid file id", http.StatusBadRequest)
			return
		}
		file, err := dbClient.GetFile(id)
		if err != nil {
			http.Error(w, "Error retrieving file from database", http.StatusInternalServerError)
			return
		}
		fileBytes, err := storageClient.GetFile(id)
		if err != nil {
			http.Error(w, "Error retrieving file from storage", http.StatusInternalServerError)
			return
		}
		res, err := json.Marshal(&GetFileResponse{
			ID:       file.ID,
			Name:     file.Name,
			Metadata: file.Metadata,
			Data:     fileBytes,
		})
		if err != nil {
			http.Error(w, "Error marshalling response", http.StatusInternalServerError)
			return
		}
		w.Write(res)
	})

	r.Delete("/files/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "Invalid file id", http.StatusBadRequest)
			return
		}
		err = dbClient.DeleteFile(id)
		if err != nil {
			http.Error(w, "Error deleting file from database", http.StatusInternalServerError)
			return
		}
		err = storageClient.DeleteFile(id)
		if err != nil {
			http.Error(w, "Error deleting file from storage", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	r.Get("/files/search", func(w http.ResponseWriter, r *http.Request) {
		var req SearchFilesRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}
		files, err := dbClient.SearchFiles(req.Query)
		if err != nil {
			http.Error(w, "Error searching files", http.StatusInternalServerError)
			return
		}
		res, err := json.Marshal(&SearchFilesResponse{Files: files})
		if err != nil {
			http.Error(w, "Error marshalling response", http.StatusInternalServerError)
			return
		}
		w.Write(res)
	})

	fmt.Println("files-api server listening on port 8080")
	return http.ListenAndServe(":8080", r)
}
