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
	"github.com/spaolacci/murmur3"
)

// Strategy a strategy to translate byte array to k bit indexes.
type Strategy interface {
	// Indexes gets indexes.
	Indexes(data []byte, m uint, k uint) []uint
}

// MURMUR128MITZ64 is a hashing strategy using murmur3
type MURMUR128MITZ64 struct {
}

// Indexes gets indexes.
func (strategy *MURMUR128MITZ64) Indexes(data []byte, m uint, k uint) []uint {
	indexes := make([]uint, k)
	h1, h2 := murmur3.Sum128(data)
	combined := int64(h1)

	for i := uint(0); i < k; i++ {
		indexes[i] = uint((combined & int64(0x7fffffffffffffff)) % int64(m))
		combined += int64(h2)
	}
	return indexes
}
