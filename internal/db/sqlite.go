package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type SQLiteClient struct {
	db *sql.DB
}

var _ DBClient = SQLiteClient{}

func NewSQLiteClient(db *sql.DB) (*SQLiteClient, error) {
	return &SQLiteClient{db: db}, nil
}

func (s SQLiteClient) SaveFile(name string, metadata string) (
	id int, err error,
) {
	result, err := s.db.Exec(
		"INSERT INTO files (name, metadata) VALUES (?, ?)",
		name,
		metadata,
	)
	if err != nil {
		fmt.Printf("Error saving file: %s", err)
		return 0, err
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastInsertID), nil
}

func (s SQLiteClient) GetFile(id int) (*File, error) {
	rows, err := s.db.Query("SELECT * FROM files WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	f := &File{}
	if rows.Next() {
		err = rows.Scan(&f.ID, &f.Name, &f.Metadata)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("file not found")
	}
	return f, nil
}

func (s SQLiteClient) DeleteFile(id int) error {
	_, err := s.db.Exec("DELETE FROM files WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (s SQLiteClient) SearchFiles(query string) ([]*File, error) {
	rows, err := s.db.Query(
		"SELECT * FROM files WHERE metadata LIKE ?",
		"%"+query+"%",
	)
	if err != nil {
		return nil, err
	}
	files := []*File{}
	for rows.Next() {
		f := &File{}
		err = rows.Scan(&f.ID, &f.Name, &f.Metadata)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func JSONStrToMap(metadata string) (map[string]interface{}, error) {
	var meta map[string]interface{}
	if metadata != "" {
		err := json.Unmarshal([]byte(metadata), &meta)
		if err != nil {
			return nil, err
		}
	}
	return meta, nil
}
