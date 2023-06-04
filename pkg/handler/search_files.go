package handler

import (
	"net/http"

	"github.com/dgurns/files-api/internal/db"
	"github.com/gin-gonic/gin"
)

type SearchResult struct {
	ID       int                    `json:"id"`
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
}

type SearchFilesResponse struct {
	Results []*SearchResult `json:"results"`
}

func (h *Handler) SearchFiles(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.String(http.StatusBadRequest, "Missing query parameter")
		return
	}

	files, err := h.db.SearchFiles(query)
	if err != nil {
		c.String(http.StatusNotFound, "No files found")
		return
	}

	results := []*SearchResult{}
	for _, f := range files {
		meta, err := db.JSONStrToMap(f.Metadata)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error unmarshalling metadata")
			return
		}
		results = append(results, &SearchResult{
			ID:       f.ID,
			Name:     f.Name,
			Metadata: meta,
		})
	}

	c.JSON(http.StatusOK, &SearchFilesResponse{Results: results})
}
