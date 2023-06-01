package storage

type FilesystemClient struct{}

var _ StorageClient = FilesystemClient{}

func NewFilesystemClient() (*FilesystemClient, error) {
	return &FilesystemClient{}, nil
}

func (f FilesystemClient) SaveFile(id int, data []byte) error {
	return nil
}

func (f FilesystemClient) GetFile(id int) ([]byte, error) {
	return nil, nil
}

func (f FilesystemClient) DeleteFile(id int) error {
	return nil
}
