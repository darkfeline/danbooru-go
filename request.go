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
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// newRequest creates an HTTP request for an API call with the client.
// The request is authenticated with the client auth info if set.  If
// the body implements contentTyper, the contentTyper interface is
// used to set the Content-Type header.
func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) *http.Request {
	u := url.URL{
		Scheme: "https",
		Host:   c.Host,
		Path:   path,
	}
	r, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		panic(err)
	}
	r = r.WithContext(ctx)
	if c.AuthInfo.Login != "" && c.AuthInfo.APIKey != "" {
		r.SetBasicAuth(c.AuthInfo.Login, c.AuthInfo.APIKey)
	}
	if b, ok := body.(contentTyper); ok {
		r.Header.Set("Content-Type", b.ContentType())
	}
	return r
}

// contentTyper is the interface for getting the Content-Type to use
// for an HTTP request body.
type contentTyper interface {
	ContentType() string
}

// jsonBody wraps a Reader with application/json content type.
type jsonBody struct {
	io.Reader
}

func (jsonBody) ContentType() string {
	return "application/json"
}

// newJSONBody returns an encoded jsonBody for the given value.  If
// the value cannot be marshaled, panic.
func newJSONBody(v interface{}) jsonBody {
	d, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return jsonBody{bytes.NewReader(d)}
}

// A PostID is the ID of a post.
type PostID int

func newPostIDBody(id PostID) jsonBody {
	type body struct {
		PostID PostID `json:"post_id"`
	}
	b := body{
		PostID: id,
	}
	return newJSONBody(&b)
}
