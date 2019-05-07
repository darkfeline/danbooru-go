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
	"fmt"
	"io"
	"net/http"

	"golang.org/x/xerrors"
)

type responseBody struct {
	Success   bool     `json:"success"`
	Message   string   `json:"message"`
	Backtrace []string `json:"backtrace"`
}

func parseBody(r io.Reader) (responseBody, error) {
	d := json.NewDecoder(r)
	var rb responseBody
	if err := d.Decode(&rb); err != nil {
		return rb, xerrors.Errorf("parse body: %w", err)
	}
	return rb, nil
}

// ErrThrottled is returned when an API call is throttled.  Use
// errors.Is to compare.
var ErrThrottled = responseError{StatusCode: 429}

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

type responseError struct {
	StatusCode int
	body       responseBody
}

var _ error = responseError{}
var _ fmt.Formatter = responseError{}
var _ xerrors.Formatter = responseError{}

func (e responseError) Error() string {
	return fmt.Sprintf("HTTP %d %s", e.StatusCode, http.StatusText(e.StatusCode))
}

func (e responseError) Is(target error) bool {
	e2, ok := target.(responseError)
	if !ok {
		return false
	}
	return e.StatusCode == e2.StatusCode
}

func (e responseError) Format(f fmt.State, c rune) {
	xerrors.FormatError(e, f, c)
}

func (e responseError) FormatError(p xerrors.Printer) error {
	p.Print(e.Error())
	if !p.Detail() {
		return nil
	}
	p.Printf("%+v", e.body)
	return nil
}
