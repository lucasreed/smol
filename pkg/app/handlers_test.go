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

package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lucasreed/smol/pkg/data/models"
)

type storage struct {
	data map[string]string
}

func (s *storage) Open() error {
	return nil
}

func (s *storage) Close() error {
	return nil
}

func (s *storage) Health() bool {
	return true
}

func (s *storage) GetURL(shortCode string) (models.URL, error) {
	if dest, ok := s.data[shortCode]; ok {
		return models.URL{Destination: dest, ShortCode: shortCode}, nil
	}
	return models.URL{}, fmt.Errorf("code not found")
}

func (s *storage) GetShortCode(destination string) (string, error) {
	if code, ok := s.data[destination]; ok {
		return code, nil
	}
	return "", fmt.Errorf("destination not found")
}

func (s *storage) SetURL(shortCode, url string) error {
	s.data[shortCode] = url
	return nil
}

func (s *storage) Delete(shortCode string) error {
	delete(s.data, shortCode)
	return nil
}

var testStorage = storage{
	data: map[string]string{
		"abcd123":            "https://google.com",
		"https://google.com": "abcd123",
	},
}

var server = Server{
	Listen:  "",
	router:  nil,
	Storage: &testStorage,
}

func TestHandleAdd(t *testing.T) {
	requestBody, err := json.Marshal(map[string]string{
		"Destination": "https://lreed.net",
	})
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/add", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handleAdd)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// func TestHandleShortCode(t *testing.T) {
// 	req, err := http.NewRequest("GET", "/abcd123", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(server.handleShortCode)
// 	handler.ServeHTTP(rr, req)
//
// 	if status := rr.Code; status != http.StatusPermanentRedirect {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}
// }
