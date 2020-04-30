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
	"testing"
)

var testCases = []string{
	"https://example.com",
	"https://example.com/path",
	"https://example.com/abc?de=1&fg=true",
	"http://example.com",
	"http://example.com/path",
	"http://example.com/abc?de=1&fg=true",
	"example.com",
	"example.com/path",
	"example.com/abc?de=1&fg=true",
}

var failCases = []string{
	"http:/example.com",
	"http:example.com",
	"http//example.com",
	"://example.com/abc?de=1&sandwich=true",
	":/example.com",
	":example.com",
	"//example.com",
	"https:/example.com",
	"https:example.com",
	"https//example.com",
}

func TestURL_validateProtocolSuccess(t *testing.T) {
	for _, tc := range testCases {
		// t.Logf("Testing %s", tc)
		u := URL{
			Destination: tc,
		}
		if !u.ValidateURL() {
			t.Errorf("Expected success, but %s failed", tc)
			t.Fail()
		}
	}
}

func TestURL_validateProtocolFail(t *testing.T) {
	for _, tc := range failCases {
		// t.Logf("Testing %s", tc)
		u := URL{
			Destination: tc,
		}
		if u.ValidateURL() {
			t.Errorf("Expected failure, but %s succeeded", tc)
			t.Fail()
		}
	}
}
