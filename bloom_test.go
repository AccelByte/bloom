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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBloomFilter_M_K(t *testing.T) {
	f := NewWithFPP(1000, 0.01)
	assert.Equal(t, uint(9600), f.M())
	assert.Equal(t, uint(7), f.K())

	f = NewWithFPP(10000, 0.03)
	assert.Equal(t, uint(73024), f.M())
	assert.Equal(t, uint(5), f.K())

	f = NewWithFPP(1024, 0.01)
	assert.Equal(t, uint(9856), f.M())
	assert.Equal(t, uint(7), f.K())
}

func TestBloomFilter(t *testing.T) {
	s1 := "this_is_a_test_string"
	s2 := "this_is_another_test_string"
	s3 := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"
	s4 := "dyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"

	f := NewWithFPP(100, 0.01)
	f.Put([]byte(s1)).Put([]byte(s3))
	assert.True(t, f.MightContain([]byte(s1)))
	assert.False(t, f.MightContain([]byte(s2)))
	assert.True(t, f.MightContain([]byte(s3)))
	assert.False(t, f.MightContain([]byte(s4)))
}

func TestEmptyFilter(t *testing.T) {
	f := NewWithFPP(0, 0.01)
	assert.True(t, f.m > 0)
}

func TestZeroFPP(t *testing.T) {
	f := NewWithFPP(1, 0)
	assert.True(t, f.m > 0)
}

func TestBloomFilter_JSON(t *testing.T) {
	s1 := "this_is_a_test_string"
	s2 := "this_is_another_test_string"
	s3 := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"
	s4 := "dyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"

	f := NewWithFPP(100, 0.01)
	f.Put([]byte(s1)).Put([]byte(s3))
	exported, _ := f.MarshalJSON()
	bloomFilterJSON := &FilterJSON{}
	json.Unmarshal(exported, bloomFilterJSON)
	newF := From(bloomFilterJSON.B, bloomFilterJSON.K)
	assert.True(t, newF.MightContain([]byte(s1)))
	assert.False(t, newF.MightContain([]byte(s2)))
	assert.True(t, newF.MightContain([]byte(s3)))
	assert.False(t, newF.MightContain([]byte(s4)))
}
