package handler

import (
	"github.com/dgurns/files-api/internal/db"
	"github.com/dgurns/files-api/internal/storage"
)

type Handler struct {
	db      *db.SQLiteClient
	storage *storage.FilesystemClient
}

func NewHandler(dbClient *db.SQLiteClient, storageClient *storage.FilesystemClient) *Handler {
	return &Handler{
		db:      dbClient,
		storage: storageClient,
	}
}
