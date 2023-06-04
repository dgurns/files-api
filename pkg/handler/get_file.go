package handler

import (
	"net/http"
	"strconv"

	"github.com/dgurns/files-api/internal/db"
	"github.com/gin-gonic/gin"
)

type GetFileResponse struct {
	ID       int                    `json:"id"`
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
	Data     []byte                 `json:"data"`
}

func (h *Handler) GetFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid file id")
		return
	}

	file, err := h.db.GetFile(id)
	if err != nil {
		c.String(http.StatusNotFound, "File not found in database")
		return
	}

	meta, err := db.JSONStrToMap(file.Metadata)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error unmarshalling metadata")
		return
	}

	fileBytes, err := h.storage.GetFile(id)
	if err != nil {
		c.String(http.StatusNotFound, "File not found in storage")
		return
	}

	c.JSON(http.StatusOK, &GetFileResponse{
		ID:       file.ID,
		Name:     file.Name,
		Metadata: meta,
		Data:     fileBytes,
	})
}
