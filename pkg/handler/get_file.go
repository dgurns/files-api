package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dgurns/files-api/internal/db"
	"github.com/go-chi/chi/v5"
)

type GetFileResponse struct {
	ID       int                    `json:"id"`
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
	Data     []byte                 `json:"data"`
}

func (h *Handler) GetFile(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid file id", http.StatusBadRequest)
		return
	}
	file, err := h.db.GetFile(id)
	if err != nil {
		http.Error(w, "File not found in database", http.StatusNotFound)
		return
	}
	meta, err := db.JSONStrToMap(file.Metadata)
	if err != nil {
		http.Error(w, "Error unmarshalling metadata", http.StatusInternalServerError)
		return
	}
	fileBytes, err := h.storage.GetFile(id)
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
}
