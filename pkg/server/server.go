package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/dgurns/files-api/internal/db"
	"github.com/dgurns/files-api/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Run() error {
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

	storagePath := os.Getenv("LOCAL_FILES_PATH")
	storageClient, err := storage.NewFilesystemClient(storagePath)
	if err != nil {
		return err
	}

	r := chi.NewRouter()

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

		allowed := []string{"application/pdf", "image/jpeg", "image/png"}
		for _, a := range allowed {
			if handler.Header.Get("Content-Type") == a {
				break
			}
			http.Error(w, "Invalid file type: only PDF, JPEG, or PNG are allowed", http.StatusBadRequest)
			return
		}
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}

		metadata := r.FormValue("metadata")
		if metadata != "" {
			var m map[string]interface{}
			err = json.Unmarshal([]byte(metadata), &m)
			if err != nil {
				http.Error(w, "Invalid metadata: must be valid JSON", http.StatusBadRequest)
				return
			}
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
			http.Error(w, "File not found in database", http.StatusNotFound)
			return
		}
		meta, err := db.JSONStrToMap(file.Metadata)
		if err != nil {
			http.Error(w, "Error unmarshalling metadata", http.StatusInternalServerError)
			return
		}
		fileBytes, err := storageClient.GetFile(id)
		if err != nil {
			http.Error(w, "File not found in storage", http.StatusNotFound)
			return
		}
		res, err := json.Marshal(&GetFileResponse{
			ID:       file.ID,
			Name:     file.Name,
			Metadata: meta,
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
		query := r.URL.Query().Get("query")
		if query == "" {
			http.Error(w, "Missing query parameter", http.StatusBadRequest)
			return
		}
		files, err := dbClient.SearchFiles(query)
		if err != nil {
			http.Error(w, "No files found", http.StatusNotFound)
			return
		}
		results := []*SearchResult{}
		for _, f := range files {
			meta, err := db.JSONStrToMap(f.Metadata)
			if err != nil {
				http.Error(w, "Error unmarshalling metadata", http.StatusInternalServerError)
				return
			}
			results = append(results, &SearchResult{
				ID:       f.ID,
				Name:     f.Name,
				Metadata: meta,
			})
		}
		res, err := json.Marshal(&SearchFilesResponse{Results: results})
		if err != nil {
			http.Error(w, "Error marshalling response", http.StatusInternalServerError)
			return
		}
		w.Write(res)
	})

	fmt.Println("files-api server listening on port 8080")
	return http.ListenAndServe(":8080", r)
}
