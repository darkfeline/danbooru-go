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
	"fmt"
	"testing"

	"golang.org/x/xerrors"
)

func TestResponseError_Error(t *testing.T) {
	t.Parallel()
	assert := func(t *testing.T, got, want string) {
		if got != want {
			t.Errorf("Error() = %#v; want %#v", got, want)
		}
	}
	t.Run("no message", func(t *testing.T) {
		t.Parallel()
		err := responseError{
			StatusCode: 500,
		}
		got := err.Error()
		want := `HTTP 500 Internal Server Error`
		assert(t, got, want)
	})
	t.Run("message", func(t *testing.T) {
		t.Parallel()
		err := responseError{
			StatusCode: 500,
			body: responseBody{
				Message: "blah",
			},
		}
		got := err.Error()
		want := `HTTP 500 Internal Server Error: blah`
		assert(t, got, want)
	})
	t.Run("wrapped with message", func(t *testing.T) {
		t.Parallel()
		var err error
		err = responseError{
			StatusCode: 500,
			body: responseBody{
				Message: "blah",
			},
		}
		err = xerrors.Errorf("add favorite %d: %w", 123, err)
		got := err.Error()
		want := `add favorite 123: HTTP 500 Internal Server Error: blah`
		assert(t, got, want)
	})
}

func TestResponseError_Is(t *testing.T) {
	t.Parallel()
	t.Run("ErrThrottled", func(t *testing.T) {
		t.Parallel()
		err := responseError{
			StatusCode: 429,
			body: responseBody{
				Success: false,
				Message: "throttled",
			},
		}
		if !xerrors.Is(err, ErrThrottled) {
			t.Errorf("%#v should be ErrThrottled", err)
		}
	})
}

func TestResponseError_FormatError(t *testing.T) {
	t.Parallel()
	assert := func(t *testing.T, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf(`Unexpected format
Got:
%s

Want:
%s`, got, want)
		}
	}
	t.Run("string format", func(t *testing.T) {
		t.Parallel()
		err := responseError{
			StatusCode: 429,
		}
		got := fmt.Sprintf("%v", err)
		want := `HTTP 429 Too Many Requests`
		assert(t, got, want)
	})
	t.Run("string format with message", func(t *testing.T) {
		t.Parallel()
		err := responseError{
			StatusCode: 429,
			body: responseBody{
				Success: false,
				Message: "throttled",
			},
		}
		got := fmt.Sprintf("%v", err)
		want := `HTTP 429 Too Many Requests: throttled`
		assert(t, got, want)
	})
	t.Run("detail format with body", func(t *testing.T) {
		t.Parallel()
		err := responseError{
			StatusCode: 429,
			body: responseBody{
				Success: false,
				Message: "throttled",
			},
		}
		got := fmt.Sprintf("%+v", err)
		want := `HTTP 429 Too Many Requests: throttled:
    {Success:false Message:throttled Backtrace:[]}`
		assert(t, got, want)
	})
}
