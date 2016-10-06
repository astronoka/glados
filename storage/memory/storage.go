package memory

import (
	"encoding/json"
	"fmt"

	"github.com/astronoka/glados"
)

// NewStorage is create mysql storage instance
func NewStorage(c glados.Context) glados.Storage {
	return &memoryStorage{
		data: map[string][]byte{},
	}
}

type memoryStorage struct {
	data map[string][]byte
}

func (s *memoryStorage) Save(namespace, key string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("memory : serialize %s:%s data faild. %s", namespace, key, err.Error())
	}
	s.data[namespace+":"+key] = b
	return nil
}

func (s *memoryStorage) Load(namespace, key string, value interface{}) (bool, error) {
	if _, exist := s.data[namespace+":"+key]; !exist {
		return false, nil
	}
	err := json.Unmarshal(s.data[namespace+":"+key], value)
	if err != nil {
		return false, fmt.Errorf("memory: deserialize %s:%s failed. %s", namespace, key, err.Error())
	}
	return true, nil
}

func (s *memoryStorage) Delete(namespace, key string) error {
	delete(s.data, namespace+":"+key)
	return nil
}
