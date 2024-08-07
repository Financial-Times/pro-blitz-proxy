package storage

type IStorage interface {
	Exists(key string) bool
	Save(key string, data []byte, headers map[string][]string) error
	Get(key string) ([]byte, map[string][]string, error)
	DeleteById(id string) error
}
