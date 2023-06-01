package db

type File struct {
	ID       int
	Name     string
	Metadata map[string]string
}

type DBClient interface {
	SaveFile(name string, metadata map[string]string) (id int, err error)
	GetFile(id int) (*File, error)
	DeleteFile(id int) error
	SearchFiles(query string) ([]*File, error)
}
