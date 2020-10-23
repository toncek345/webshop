package storage

type File struct {
	Key  string
	Data []byte
}

type Storage interface {
	Put(key string, data []byte) (*File, error)
	Get(key string) (*File, error)
	Delete(key string) error
}
