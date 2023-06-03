package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid file id", http.StatusBadRequest)
		return
	}

	err = h.db.DeleteFile(id)
	if err != nil {
		http.Error(w, "Error deleting file from database", http.StatusInternalServerError)
		return
	}

	err = h.storage.DeleteFile(id)
	if err != nil {
		http.Error(w, "Error deleting file from storage", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
