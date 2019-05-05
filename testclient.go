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
	"testing"

	"golang.org/x/xerrors"
)

func testClient(t *testing.T) *Client {
	return &Client{
		Host: "api.example.com",
		AuthInfo: AuthInfo{
			Login:  "username",
			APIKey: "secret",
		},
		HTTPClient: http.Client{
			Transport: &transportStub{
				err: xerrors.New("not implemented"),
			},
		},
		Logger: testLogger{t},
	}
}

// testLogger is a wrapper implementing Logger.
type testLogger struct {
	t *testing.T
}

func (t testLogger) Print(v ...interface{}) {
	t.t.Log(v...)
}

func (t testLogger) Printf(format string, v ...interface{}) {
	t.t.Logf(format, v...)
}

var _ Logger = testLogger{}

// transportStub is a stub implementing http.RoundTripper for tests.
type transportStub struct {
	resp *http.Response
	err  error
}

var _ http.RoundTripper = &transportStub{}

func (t *transportStub) RoundTrip(*http.Request) (*http.Response, error) {
	return t.resp, t.err
}

// transportSpy is a spy implementing http.RoundTripper for tests.
type transportSpy struct {
	transportStub
	req *http.Request
}

var _ http.RoundTripper = &transportSpy{}

func (t *transportSpy) RoundTrip(r *http.Request) (*http.Response, error) {
	t.req = r
	return t.transportStub.RoundTrip(r)
}
