package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestSearchFilesMissingQuery(t *testing.T) {
	req, err := http.NewRequest("GET", "/files/search", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	r := gin.Default()
	h := NewTestHandler()
	r.GET("/files/search", h.SearchFiles)

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Missing query parameter")
}

func TestSearchFiles(t *testing.T) {
	req, err := http.NewRequest("GET", "/files/search?query=test", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	r := gin.Default()
	h := NewTestHandler()
	r.GET("/files/search", h.SearchFiles)

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, rr.Body.String(), "{\"results\":[]}")
}
