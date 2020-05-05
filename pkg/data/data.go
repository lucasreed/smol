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

package data

import (
	"github.com/lucasreed/smol/pkg/data/models"
)

type StorageReader interface {
	GetURL(shortCode string) (models.URL, error)
	GetShortCode(destination string) (string, error)
	Health() bool
}

type StorageWriter interface {
	Open() error
	Close() error
	SetURL(shortCode, url string) error
}

type StorageReadWrite interface {
	StorageReader
	StorageWriter
}
