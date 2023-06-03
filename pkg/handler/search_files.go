package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dgurns/files-api/internal/db"
)

type SearchResult struct {
	ID       int                    `json:"id"`
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
}

type SearchFilesResponse struct {
	Results []*SearchResult `json:"results"`
}

func (h *Handler) SearchFiles(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}

	files, err := h.db.SearchFiles(query)
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
}
