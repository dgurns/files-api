package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestGetFile(t *testing.T) {
	req, err := http.NewRequest("GET", "/files/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	r := gin.Default()
	h := NewTestHandler()
	r.GET("/files/:id", h.GetFile)

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	resBody := rr.Body.String()
	res := GetFileResponse{}
	err = json.Unmarshal([]byte(resBody), &res)
	assert.NoError(t, err)
	assert.Equal(t, 1, res.ID)
}
