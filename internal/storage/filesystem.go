package storage

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type FilesystemClient struct {
	StoragePath string
}

var _ Client = FilesystemClient{}

func NewFilesystemClient(storagePath string) (*FilesystemClient, error) {
	return &FilesystemClient{
		StoragePath: storagePath,
	}, nil
}

func (f FilesystemClient) SaveFile(id int, data []byte) error {
	filePath := filepath.Join(f.StoragePath, strconv.Itoa(id))
	err := ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (f FilesystemClient) GetFile(id int) ([]byte, error) {
	filePath := filepath.Join(f.StoragePath, strconv.Itoa(id))
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (f FilesystemClient) DeleteFile(id int) error {
	filePath := filepath.Join(f.StoragePath, strconv.Itoa(id))
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}
