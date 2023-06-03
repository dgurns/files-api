package db

type File struct {
	ID       int
	Name     string
	Metadata string
}

type Client interface {
	SaveFile(name string, metadata string) (id int, err error)
	GetFile(id int) (*File, error)
	DeleteFile(id int) error
	SearchFiles(query string) ([]*File, error)
}
