package persistence

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/dgraph-io/badger"
	"log"
	"strings"
)

type Cache struct {
	db       *badger.DB
	bucket   string
	sequence *badger.Sequence
}

func NewCache(db *badger.DB, bucket string) *Cache {
	cache := new(Cache)

	cache.db = db
	cache.bucket = bucket

	sequence, err := db.GetSequence([]byte(bucket), 100)

	if err != nil {
		log.Fatal(err.Error())
	}

	cache.sequence = sequence

	return cache
}

func appendIfMissing(slice []string, str string) []string {
	for _, element := range slice {
		if element == str {
			return slice
		}
	}
	return append(slice, str)
}

func marshal(value interface{}) ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(&value)

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func unmarshal(data []byte) (interface{}, error) {
	var value interface{}

	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(&value)

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (cache *Cache) GetValidID(key string) string {

	originalKey := key

	var newKey string

	err := cache.db.Update(func(txn *badger.Txn) error {
		iter := 1

		for {
			fullKey := []byte(fmt.Sprintf("%s_%s", cache.bucket, key))

			val, _ := txn.Get(fullKey)

			if val == nil {
				newKey = key
				txn.Set(fullKey, fullKey)
				return nil
			} else {
				key = fmt.Sprintf("%s-%d", originalKey, iter)
			}

			iter++
		}

		return nil

	})

	if err != nil {
		log.Fatal(err)
	}

	return newKey
}

func (cache *Cache) Get(key string) (interface{}, error) {
	var returnValue interface{}

	fullKey := []byte(fmt.Sprintf("%s_%s", cache.bucket, key))

	err := cache.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(fullKey)

		if err != nil {
			log.Println(err.Error())
			return err
		}

		valueData, err := item.Value()

		if err != nil {
			log.Println(err.Error())
			return err
		}

		value, err := unmarshal(valueData)

		if err != nil {
			log.Println(err.Error())
			return err
		}

		returnValue = value

		return nil
	})

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return returnValue, nil
}

func (cache *Cache) Delete(key string) error {
	fullKey := []byte(fmt.Sprintf("%s_%s", cache.bucket, key))

	err := cache.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(fullKey)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (cache *Cache) Set(key string, value interface{}) {
	fullKey := []byte(fmt.Sprintf("%s_%s", cache.bucket, key))

	encodedValue, err := marshal(value)

	if err != nil {
		log.Fatal(err)
	}

	err = cache.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(fullKey, encodedValue)
		return err
	})

	if err != nil {
		log.Fatal(err)
	}
}

func (cache *Cache) TopLevelKeys() []string {
	keys := []string{}

	prefix := []byte(fmt.Sprintf("%s_", cache.bucket))

	err := cache.db.View(func(tx *badger.Txn) error {

		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := tx.NewIterator(opts)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {

			item := it.Item()
			keyAndPrefix := item.Key()

			keyOnly := strings.Split(string(keyAndPrefix), string(prefix))[1]

			keys = appendIfMissing(keys, keyOnly)
		}
		return nil

	})

	if err != nil {
		log.Fatal(err.Error())
	}

	return keys

}
