package storage

type Storage interface {
	Put(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	List(prefix string) ([]string, error)
	Close() error
}
