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
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/lucasreed/smol/pkg/data/models"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "Home\n")
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleIgnore(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleAdd(w http.ResponseWriter, r *http.Request) {
	var path string
	var urlModel models.URL
	js := json.NewDecoder(r.Body)
	err := js.Decode(&urlModel)
	if err != nil {
		log.Println("error decoding json")
		w.WriteHeader(http.StatusInternalServerError)
		_, innerErr := w.Write([]byte("error decoding json"))
		if innerErr != nil {
			log.Printf("ERROR: %v", innerErr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if len(urlModel.Destination) == 0 {
		log.Println("destination field not provided")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("destination field not provided"))
		if err != nil {
			log.Printf("ERROR: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if !urlModel.ValidateURL() {
		message := fmt.Sprintf("url is not valid: %s", urlModel.Destination)
		log.Println(message)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(message))
		if err != nil {
			log.Printf("ERROR: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if p, exists := s.urlRegistered(urlModel.Destination); exists {
		w.WriteHeader(http.StatusFound)
		message := fmt.Sprintf("This url is already registered: %s -> %s", p, urlModel.Destination)
		log.Println(message)
		_, err = w.Write([]byte(message))
		if err != nil {
			log.Printf("ERROR: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	for i := 0; i < 3; i++ {
		p := createShortCode(7)
		if !s.pathRegistered(p) {
			path = p
			break
		}
		if i == 2 {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte("error creating a shortCode path"))
			if err != nil {
				log.Printf("ERROR: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}
	}
	urlModel.ShortCode = path
	err = s.Storage.SetURL(path, urlModel.Destination)
	if err != nil {
		log.Printf("failed to store url - %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		_, innerErr := w.Write([]byte("failed to store url"))
		if innerErr != nil {
			log.Printf("ERROR: %v", innerErr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	w.WriteHeader(http.StatusAccepted)
	log.Printf("Added path: %s, url: %s\n", urlModel.ShortCode, urlModel.Destination)
}

func (s *Server) handleShortCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]
	url, err := s.Storage.GetURL(shortCode)
	if err != nil {
		log.Printf("error finding shortcode, maybe it does not exist: %s\n", shortCode)
		w.WriteHeader(http.StatusNotFound)
		_, innerErr := w.Write([]byte("error finding shortcode, maybe it does not exist: " + shortCode))
		if innerErr != nil {
			log.Printf("ERROR: %v", innerErr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	log.Printf("Redirecting from %s to %s\n", r.URL.EscapedPath(), url.Destination)
	http.Redirect(w, r, url.Destination, http.StatusPermanentRedirect)

}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]
	if exists := s.pathRegistered(shortCode); !exists {
		w.WriteHeader(http.StatusNotFound)
		message := fmt.Sprintf("This short code is not registered: %s", shortCode)
		log.Println(message)
		_, err := w.Write([]byte(message))
		if err != nil {
			log.Printf("ERROR: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	err := s.Storage.Delete(shortCode)
	if err != nil {
		log.Printf("error deleting shortcode: %s\n", shortCode)
		w.WriteHeader(http.StatusNotFound)
		_, innerErr := w.Write([]byte("error deleting shortcode: " + shortCode))
		if innerErr != nil {
			log.Printf("ERROR: %v", innerErr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	w.WriteHeader(http.StatusAccepted)
	log.Printf("Deleted shortcode: %s\n", shortCode)
}

func (s *Server) urlRegistered(url string) (string, bool) {
	data, err := s.Storage.GetShortCode(url)
	if err != nil {
		return "", false
	}
	return data, true
}

func (s *Server) pathRegistered(shortCode string) bool {
	_, err := s.Storage.GetURL(shortCode)
	return err == nil
}

func createShortCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
