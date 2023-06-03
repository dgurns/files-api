package db

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestSQLiteClient(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE files (
    id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
    metadata JSON
	);`)
	assert.NoError(t, err)
	defer db.Close()

	client, err := NewSQLiteClient(db)
	if err != nil {
		t.Fatal(err)
	}

	// Test SaveFile
	id, err := client.SaveFile("test.txt", `{"size": 1024}`)
	assert.NoError(t, err)
	assert.Equal(t, 1, id)

	// Test SearchFiles
	id2, err := client.SaveFile("test2.txt", `{"size": 2048}`)
	assert.NoError(t, err)
	files, err := client.SearchFiles("2048")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(files))
	assert.Equal(t, id2, files[0].ID)

	// Test GetFile
	f, err := client.GetFile(id2)
	assert.NoError(t, err)
	assert.Equal(t, id2, f.ID)
}
