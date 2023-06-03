package handler

import (
	"github.com/dgurns/files-api/internal/db"
	"github.com/dgurns/files-api/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

type MockDBClient struct{}

var _ db.Client = MockDBClient{}

func (m MockDBClient) GetFile(id int) (*db.File, error) {
	return &db.File{ID: id}, nil
}
func (m MockDBClient) SaveFile(name string, metadata string) (int, error) {
	return 1, nil
}
func (m MockDBClient) DeleteFile(id int) error {
	return nil
}
func (m MockDBClient) SearchFiles(query string) ([]*db.File, error) {
	return []*db.File{}, nil
}

type MockStorageClient struct{}

var _ storage.Client = MockStorageClient{}

func (m MockStorageClient) GetFile(id int) ([]byte, error) {
	return []byte{}, nil
}
func (m MockStorageClient) SaveFile(id int, data []byte) error {
	return nil
}
func (m MockStorageClient) DeleteFile(id int) error {
	return nil
}

func NewTestHandler() *Handler {
	return &Handler{
		db:      &MockDBClient{},
		storage: &MockStorageClient{},
	}
}
