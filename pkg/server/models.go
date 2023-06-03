package server

import "github.com/dgurns/files-api/internal/db"

type UploadFileResponse struct {
	ID int `json:"id"`
}

type GetFileResponse struct {
	ID       int               `json:"id"`
	Name     string            `json:"name"`
	Metadata map[string]string `json:"metadata"`
	Data     []byte            `json:"data"`
}

type SearchFilesRequest struct {
	Query string `json:"query"`
}

type SearchFilesResponse struct {
	Files []*db.File `json:"files"`
}
