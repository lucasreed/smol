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
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/lucasreed/smol/pkg/data"
)

type Server struct {
	Listen  string
	router  *mux.Router
	Storage data.StorageReadWrite
}

func NewServer(storageRW data.StorageReadWrite, listenAddress string) *Server {
	return &Server{
		Listen:  listenAddress,
		router:  mux.NewRouter(),
		Storage: storageRW,
	}
}

func (s *Server) Run() {
	// Handle basic root paths
	s.router.HandleFunc("/", logHandler(s.handleIndex))
	s.router.HandleFunc("/favicon.ico", s.handleIgnore)
	s.router.HandleFunc("/{shortCode}", logHandler(s.handleShortCode)).Methods("GET")

	// Set up a subrouter for /api and then each version as more subrouters below /api
	api := s.router.PathPrefix("/api").Subrouter()
	v1 := api.PathPrefix("/v1").Subrouter()
	versionedApiRoutes(v1, s)

	log.Println("Starting server:", s.Listen)
	if err := http.ListenAndServe(s.Listen, s.router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
