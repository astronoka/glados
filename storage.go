package glados

// Storage is glados data store interface
type Storage interface {
	Save(namespace, key string, value interface{}) error
	Load(namespace, key string, value interface{}) (bool, error)
	Delete(namespace, key string) error
}
