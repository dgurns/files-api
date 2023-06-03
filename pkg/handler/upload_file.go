package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type UploadFileResponse struct {
	ID int `json:"id"`
}

func (h *Handler) UploadFile(w http.ResponseWriter, r *http.Request) {
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

	id, err := h.db.SaveFile(handler.Filename, metadata)
	if err != nil {
		http.Error(w, "Error saving file to database", http.StatusInternalServerError)
		return
	}
	err = h.storage.SaveFile(id, fileBytes)
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
}
