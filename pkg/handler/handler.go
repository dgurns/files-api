package handler

import (
	"github.com/dgurns/files-api/internal/db"
	"github.com/dgurns/files-api/internal/storage"
)

type Handler struct {
	db      db.Client
	storage storage.Client
}

func New(db db.Client, storage storage.Client) *Handler {
	return &Handler{
		db:      db,
		storage: storage,
	}
}
