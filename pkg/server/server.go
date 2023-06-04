package server

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"github.com/dgurns/files-api/internal/db"
	"github.com/dgurns/files-api/internal/storage"
	"github.com/dgurns/files-api/pkg/handler"
)

func Run() error {
	// init db client
	dbConn, err := sql.Open(
		"sqlite3",
		os.Getenv("LOCAL_DB_PATH")+"/"+os.Getenv("LOCAL_DB_NAME"),
	)
	if err != nil {
		fmt.Println("Error opening database connection")
		return err
	}
	defer dbConn.Close()
	dbClient, err := db.NewSQLiteClient(dbConn)
	if err != nil {
		fmt.Println("Error creating database client")
		return err
	}

	// init storage client
	storagePath := os.Getenv("LOCAL_FILES_PATH")
	storageClient, err := storage.NewFilesystemClient(storagePath)
	if err != nil {
		return err
	}

	// set up router
	r := gin.Default()
	r.SetTrustedProxies(nil)

	username := os.Getenv("BASIC_AUTH_USERNAME")
	pwd := os.Getenv("BASIC_AUTH_PASSWORD")
	if username == "" || pwd == "" {
		fmt.Println("BASIC_AUTH_USERNAME and BASIC_AUTH_PASSWORD must be set")
		os.Exit(1)
	}
	a := r.Group("/", gin.BasicAuth(gin.Accounts{
		username: pwd,
	}))

	// handle routes
	h := handler.New(dbClient, storageClient)
	a.POST("/files/upload", h.UploadFile)
	a.GET("/files/:id", h.GetFile)
	a.DELETE("/files/:id", h.DeleteFile)
	a.GET("/files/search", h.SearchFiles)

	fmt.Println("files-api server listening on port 8080")
	return r.Run(":8080")
}
