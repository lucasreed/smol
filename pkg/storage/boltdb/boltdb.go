// Copyright 2020 Luke Reed <luke@lreed.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package boltdb

import (
	"fmt"
	"os"

	bolt "go.etcd.io/bbolt"

	"github.com/lucasreed/smol/pkg/data/models"
)

var (
	bucketName = "smol"
)

// Store represents a boltdb storage location
type Store struct {
	DB   *bolt.DB
	Path string
}

// NewStore represents a new instance of a rediscache storage location
func NewStore(path string) *Store {
	return &Store{
		Path: path,
	}
}

func (s *Store) Open() error {
	db, err := bolt.Open(s.Path, 0600, nil)
	if err != nil {
		return err
	}
	s.DB = db

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("[boltdb] error creating bucket: %s", err)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (s *Store) Close() error {
	return s.DB.Close()
}

func (s *Store) Health() bool {
	if _, err := os.Stat(s.Path); os.IsNotExist(err) {
		return false
	}
	return s.DB != nil
}

func (s *Store) GetURL(shortCode string) (models.URL, error) {
	data, err := s.getValue(shortCode)
	if err != nil {
		return models.URL{}, err
	}
	return models.URL{Destination: data, ShortCode: shortCode}, nil
}

func (s *Store) GetShortCode(destination string) (string, error) {
	return s.getValue(destination)
}

func (s *Store) SetURL(shortCode, url string) error {
	return s.DB.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if err := b.Put([]byte(shortCode), []byte(url)); err != nil {
			return err
		}
		if err := b.Put([]byte(url), []byte(shortCode)); err != nil {
			return err
		}
		return nil
	})
}

func (s *Store) getValue(key string) (string, error) {
	var value string
	err := s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		val := b.Get([]byte(key))
		value = string(val)
		return nil
	})
	if err != nil {
		return "", err
	}
	if value == "" {
		return "", fmt.Errorf("key not found: %s", key)
	}
	return value, nil
}
