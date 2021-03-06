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
	"context"
	"testing"
)

func TestClient_AddFavorite(t *testing.T) {
	t.Parallel()
	t.Run("transport error", func(t *testing.T) {
		t.Parallel()
		c := testClient(t)
		ctx := context.Background()
		err := c.AddFavorite(ctx, 1234)
		if err == nil {
			t.Errorf("Expected error")
		}
	})
}
