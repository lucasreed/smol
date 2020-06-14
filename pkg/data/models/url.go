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

package models

import (
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
)

type URL struct {
	gorm.Model
	Destination string `gorm:"index"`
	ShortCode   string `gorm:"type:varchar(7);index"`
	User        string `gorm:"index"`
}

func (urlPath *URL) ValidateURL() bool {
	return govalidator.IsURL(urlPath.Destination) && urlPath.validateProtocol()
}

func (urlPath *URL) validateProtocol() bool {
	regex := regexp.MustCompile(`^(https|http)?(.*)?`)
	subs := regex.FindStringSubmatch(urlPath.Destination)
	if subs[1] != "" {
		return strings.Contains(subs[2], "://")
	}
	urlPath.addProtocol()
	return true
}

func (urlPath *URL) addProtocol() {
	urlPath.Destination = "http://" + urlPath.Destination
}
