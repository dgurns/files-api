package db

import "fmt"

type SQLiteClient struct {
	dbPath string
}

var _ DBClient = SQLiteClient{}

func NewSQLiteClient(dbPath string) (*SQLiteClient, error) {
	return &SQLiteClient{
		dbPath: dbPath,
	}, nil
}

func (s SQLiteClient) SaveFile(name string, metadata map[string]string) (
	id int, err error,
) {
	fmt.Println("SAVING", name, metadata)
	return 1, nil
}

func (s SQLiteClient) GetFile(id int) (*File, error) {
	return &File{
		ID:       id,
		Name:     "myfile.pdf",
		Metadata: map[string]string{},
	}, nil
}

func (s SQLiteClient) DeleteFile(id int) error {
	return nil
}

func (s SQLiteClient) SearchFiles(query string) ([]*File, error) {
	return nil, nil
}
