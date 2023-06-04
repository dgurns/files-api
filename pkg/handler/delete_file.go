package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid file id")
		return
	}

	err = h.db.DeleteFile(id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error deleting file from database")
		return
	}

	err = h.storage.DeleteFile(id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error deleting file from storage")
		return
	}

	c.Status(http.StatusOK)
}
