package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadFileResponse struct {
	ID int `json:"id"`
}

func (h *Handler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Error retrieving file from form data")
		return
	}

	allowed := []string{"application/pdf", "image/jpeg", "image/png"}
	for _, a := range allowed {
		if file.Header.Get("Content-Type") == a {
			break
		}
		c.String(http.StatusBadRequest, "Invalid file type: only PDF, JPEG, or PNG are allowed")
		return
	}
	f, err := file.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error opening file contents")
		return
	}
	fileBytes, err := ioutil.ReadAll(f)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error reading file")
		return
	}

	metadata := c.PostForm("metadata")
	if metadata != "" {
		var m map[string]interface{}
		err = json.Unmarshal([]byte(metadata), &m)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid metadata: must be valid JSON")
			return
		}
	}

	id, err := h.db.SaveFile(file.Filename, metadata)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error saving file to database")
		return
	}
	err = h.storage.SaveFile(id, fileBytes)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error saving file to storage")
		return
	}

	c.JSON(http.StatusCreated, &UploadFileResponse{ID: id})
}
