/*
 * Copyright 2018 AccelByte Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package bloom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMURMUR128MITZ64 the expected data are exported from guava test results.
func TestMURMUR128MITZ64(t *testing.T) {
	expected1 := []uint{3942, 6555, 9168, 2181, 4794, 7407, 420}
	s1 := "this_is_a_test_string"
	strategy := &MURMUR128MITZ64{}
	indexes := strategy.Indexes([]byte(s1), 9600, 7)
	assert.Equal(t, expected1, indexes)
}
