// Copyright (C) 2019  Allen Li
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
// limitations under the License.

package danbooru

import (
	"net/http"
)

// Client represent a Danbooru API client.
type Client struct {
	// Host is the Danbooru instance hostname, e.g., danbooru.donmai.us.
	Host string
	// AuthInfo is used for authentication if set.
	AuthInfo AuthInfo
	// HTTPClient is the HTTP client used for requests.  The zero
	// value is suitable for use.
	HTTPClient http.Client
	// Logger is used for diagnostic logging if set.
	Logger Logger
}

// AuthInfo contains authentication info for API requests.
type AuthInfo struct {
	Login  string
	APIKey string
}

// Logger describes the interface that Client uses for logging.
type Logger interface {
	Print(...interface{})
	Printf(string, ...interface{})
}

// log prints log messages using the client's logger.
func (c *Client) log(v ...interface{}) {
	if c.Logger != nil {
		c.Logger.Print(v...)
	}
}

// logf prints log messages using the client's logger.
func (c *Client) logf(format string, v ...interface{}) {
	if c.Logger != nil {
		c.Logger.Printf(format, v...)
	}
}
