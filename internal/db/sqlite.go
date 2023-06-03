package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteClient struct {
	db *sql.DB
}

var _ DBClient = SQLiteClient{}

func NewSQLiteClient(dbPath string) (*SQLiteClient, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	return &SQLiteClient{db: db}, nil
}

func (s SQLiteClient) SaveFile(name string, metadata map[string]string) (
	id int, err error,
) {
	result, err := s.db.Exec(
		"INSERT INTO files (name, metadata) VALUES (?, ?)",
		name,
		metadata,
	)
	if err != nil {
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
	for rows.Next() {
		err = rows.Scan(&f.ID, &f.Name, &f.Metadata)
		if err != nil {
			return nil, err
		}
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
