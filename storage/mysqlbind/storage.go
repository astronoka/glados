package mysqlbind

import (
	"encoding/json"
	"fmt"

	"github.com/astronoka/glados"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // golint skip
	"github.com/naoina/migu"
)

// NewStorage is create mysql storage instance
func NewStorage(c glados.Context) glados.Storage {
	dsn := c.Env("GLADOS_DATASTORE_MYSQL_DSN", "")
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		c.Logger().Fatalln("mysqlbind: " + err.Error())
	}
	migrations, err := migu.Diff(db.DB(), "storage/mysqlbind/schema.go", nil)
	if err != nil {
		c.Logger().Fatalln("mysqlbind: migration failed. " + err.Error())
	}
	for _, m := range migrations {
		c.Logger().Infoln("mysqlbind: migrate " + m)
	}
	return &mysqlStorage{
		context: c,
		db:      db,
	}
}

type mysqlStorage struct {
	context glados.Context
	db      *gorm.DB
}

func (s *mysqlStorage) Save(namespace, key string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("mysqlbind: serialize %s:%s data faild. %s", namespace, key, err.Error())
	}
	item := Item{
		Namespace: namespace,
		Key:       key,
		Value:     string(b),
	}
	err = s.db.Create(item).Error
	if err != nil {
		return fmt.Errorf("mysqlbind: save %s:%s failed. %s", namespace, key, err.Error())
	}
	return nil
}

func (s *mysqlStorage) Load(namespace, key string, value interface{}) (bool, error) {
	where := Item{
		Namespace: namespace,
		Key:       key,
	}
	item := Item{}
	result := s.db.Find(&item, where)
	if result.Error != nil {
		if result.RecordNotFound() {
			return false, nil
		}
		return false, fmt.Errorf("mysqlbind: load %s:%s failed. %s", namespace, key, result.Error.Error())
	}
	err := json.Unmarshal([]byte(item.Value), value)
	if err != nil {
		return false, fmt.Errorf("mysqlbind: deserialize %s:%s failed. %s", namespace, key, result.Error.Error())
	}
	return true, nil
}

func (s *mysqlStorage) Delete(namespace, key string) error {
	item := Item{
		Namespace: namespace,
		Key:       key,
	}
	err := s.db.Delete(item).Error
	if err != nil {
		return fmt.Errorf("mysqlbind: delete %s:%s failed. %s", namespace, key, err.Error())
	}
	return nil
}
