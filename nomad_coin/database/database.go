package database

import (
	"errors"
	"sync"

	bolt "go.etcd.io/bbolt"
)

const DB_FILE_NAME = "database.db"

type Database struct {
	database *bolt.DB
	mutex    sync.Mutex
}

var openedBoltDBInstances []*bolt.DB

var ErrDatabaseAccessWithoutOpen = errors.New("database access without open")

func (db *Database) OpenDB(db_file_name string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	var err_db_open error
	if db.database == nil {
		db.database, err_db_open = bolt.Open(db_file_name, 0600, nil)
		if err_db_open == nil {
			openedBoltDBInstances = append(openedBoltDBInstances, db.database)
		}
	}
	return err_db_open
}

func (db *Database) CreateBucketWithStringName(name string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	if db.database != nil {
		err := db.database.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(name))
			return err
		})
		return err
	}
	return ErrDatabaseAccessWithoutOpen
}

func (db *Database) ReadByteDataFromBucket(bucketName string, key string) ([]byte, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	var data []byte
	if db.database != nil {
		err := db.database.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			data = b.Get([]byte(key))
			return nil
		})
		return data, err
	}
	return nil, ErrDatabaseAccessWithoutOpen
}

func (db *Database) WriteByteDataToBucket(bucketName string, key string, data []byte) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	if db.database != nil {
		err := db.database.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			return b.Put([]byte(key), data)
		})
		return err
	}
	return ErrDatabaseAccessWithoutOpen
}

func CloseAllOpenedDB() {
	for _, boltDB := range openedBoltDBInstances {
		if boltDB != nil {
			boltDB.Close()
		}
	}
}
