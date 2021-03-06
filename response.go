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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// responseBody represents an API response.
type responseBody struct {
	Success   bool     `json:"success"`
	Message   string   `json:"message"`
	Backtrace []string `json:"backtrace"`
}

// parseBody parses the HTTP response of an API call.  This function
// only supports JSON responses.
func parseBody(r io.Reader) (responseBody, error) {
	d := json.NewDecoder(r)
	var rb responseBody
	if err := d.Decode(&rb); err != nil {
		return rb, fmt.Errorf("parse body: %w", err)
	}
	return rb, nil
}

// ErrThrottled is returned when an API call is throttled.  Use
// errors.Is to compare.
var ErrThrottled = errors.New("throttled")

// responseError represents an error API response.
type responseError struct {
	StatusCode int
	body       responseBody
}

var _ error = responseError{}

// getResponseError returns the error for the API HTTP response.  If
// the response is not an error, return nil.
func getResponseError(r *http.Response, rb responseBody) error {
	if r.StatusCode < 400 {
		return nil
	}
	if rb.Success {
		return nil
	}
	return responseError{
		StatusCode: r.StatusCode,
		body:       rb,
	}
}

func (e responseError) Error() string {
	s := fmt.Sprintf("HTTP %d %s", e.StatusCode, http.StatusText(e.StatusCode))
	if m := e.body.Message; m != "" {
		s = s + ": " + m
	}
	return s
}

func (e responseError) Is(target error) bool {
	switch target {
	case ErrThrottled:
		return e.StatusCode == 429
	default:
		return false
	}
}
