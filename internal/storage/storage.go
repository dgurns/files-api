package storage

type Client interface {
	SaveFile(id int, data []byte) error
	GetFile(id int) ([]byte, error)
	DeleteFile(id int) error
}
