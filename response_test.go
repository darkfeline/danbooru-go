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

func TestResponseError(t *testing.T) {
	t.Parallel()
	t.Run("is ErrThrottled", func(t *testing.T) {
		t.Parallel()
		err := responseError{
			StatusCode: 429,
			body: &responseBody{
				Success: false,
				Message: "throttled",
			},
		}
		if !xerrors.Is(err, ErrThrottled) {
			t.Errorf("%#v should be ErrThrottled", err)
		}
	})
	t.Run("formatting", func(t *testing.T) {
		testResponseErrorFormatting(t)
	})
}

func testResponseErrorFormatting(t *testing.T) {
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
	t.Run("detail format", func(t *testing.T) {
		t.Parallel()
		err := responseError{
			StatusCode: 429,
		}
		got := fmt.Sprintf("%+v", err)
		want := `HTTP 429 Too Many Requests:
    <nil>`
		assert(t, got, want)
	})
	t.Run("detail format with body", func(t *testing.T) {
		t.Parallel()
		err := responseError{
			StatusCode: 429,
			body: &responseBody{
				Success: false,
				Message: "throttled",
			},
		}
		got := fmt.Sprintf("%+v", err)
		want := `HTTP 429 Too Many Requests:
    &{Success:false Message:throttled Backtrace:[]}`
		assert(t, got, want)
	})
}
