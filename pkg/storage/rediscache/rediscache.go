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

package rediscache

import (
	"fmt"

	"github.com/gomodule/redigo/redis"

	"github.com/lucasreed/smol/pkg/data/models"
)

// Store represents a rediscache storage location
type Store struct {
	Host string
	Port string
	Pool *redis.Pool
}

// NewStore represents a new instance of a rediscache storage location
func NewStore(host, port string) *Store {
	return &Store{
		Host: host,
		Port: port,
	}
}

// Open populates the pool field of the store if it has not already been set up
func (s *Store) Open() error {
	if s.Pool == nil {
		s.Pool = &redis.Pool{
			MaxIdle:   80,
			MaxActive: 12000,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", s.Host, s.Port))
				if err != nil {
					innerErr := fmt.Errorf("failed connecting to redis\n   %w", err)
					err = innerErr
				}
				return c, err
			},
		}
	}
	if !s.Health() {
		return fmt.Errorf("[redis] connection not established")
	}
	return nil
}

func (s *Store) Health() bool {
	conn := s.Pool.Get()
	data, err := redis.String(conn.Do("PING"))
	if err != nil || data != "PONG" {
		return false
	}
	return true
}

func (s *Store) Close() error {
	return s.Pool.Close()
}

func (s *Store) GetURL(shortCode string) (models.URL, error) {
	conn := s.Pool.Get()
	data, err := getValue(conn, shortCode)
	if err != nil {
		return models.URL{}, err
	}
	return models.URL{Destination: data, ShortCode: shortCode}, nil
}

func (s *Store) GetShortCode(destination string) (string, error) {
	conn := s.Pool.Get()
	return getValue(conn, destination)
}

func (s *Store) SetURL(shortCode, url string) error {
	conn := s.Pool.Get()
	err := conn.Send("SET", url, shortCode)
	if err != nil {
		return err
	}
	err = conn.Send("SET", shortCode, url)
	if err != nil {
		return err
	}
	err = conn.Flush()
	if err != nil {
		return err
	}
	return nil
}

func getValue(conn redis.Conn, key string) (string, error) {
	data, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", err
	}
	return data, nil
}
