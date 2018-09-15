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
	"math"

	"github.com/willf/bitset"
)

const defaultFPP = 1.e-5

// Filter is wrapper for bloom filter.
// Cheat sheet:
//
// m: total bits
// n: expected insertions
// b: m/n, bits per insertion
// p: expected false positive probability
// k: number of hash functions
//
// 1) Optimal k = b * ln2
// 2) p = (1 - e ^ (-kn/m))^k
// 3) For optimal k: p = 2 ^ (-k) ~= 0.6185^b
// 4) For optimal k: m = -nlnp / ((ln2) ^ 2)
type Filter struct {
	m    uint
	k    uint
	bits *bitset.BitSet
	s    Strategy
}

func pad(m uint) uint {
	r := m % 64
	d := m / 64
	if r == 0 {
		return m
	}
	return 64 * (d + 1)
}

func max(x, y uint) uint {
	if x > y {
		return x
	}
	return y
}

// New creates a new bloom filter with default FPP
// n: expected insertions
func New(n uint) *Filter {
	m, k := EstimateParameters(n, defaultFPP)
	return &Filter{pad(m), k, bitset.New(m), &MURMUR128MITZ64{}}
}

// NewWithFPP creates a new bloom filter.
// n: expected insertions
// p: expected false positive probability
func NewWithFPP(n uint, p float64) *Filter {
	m, k := EstimateParameters(n, p)
	return &Filter{pad(m), k, bitset.New(m), &MURMUR128MITZ64{}}
}

// NewWithStrategy creates a new bloom filter with strategy.
// n: expected insertions
// p: expected false positive probability
// s: index strategy
func NewWithStrategy(n uint, p float64, s Strategy) *Filter {
	m, k := EstimateParameters(n, p)
	return &Filter{pad(m), k, bitset.New(m), s}
}

// From populates a bloom filter.
// bits: long array represents bits
// k: number of hash functions
func From(bits []uint64, k uint) *Filter {
	m := uint(len(bits) * 64)
	return &Filter{m, k, bitset.From(bits), &MURMUR128MITZ64{}}
}

// FromWithStrategy populates a bloom filter with strategy.
// bits: long array represents bits
// k: number of hash functions
// s: index strategy
func FromWithStrategy(bits []uint64, k uint, s Strategy) *Filter {
	m := uint(len(bits) * 64)
	return &Filter{m, k, bitset.From(bits), s}
}

// EstimateParameters estimates requirements for m and k.
// n: expected insertions
// p: expected false positive probability
func EstimateParameters(n uint, p float64) (m uint, k uint) {
	n = max(n, 1)
	m = uint(-1 * float64(n) * math.Log(p) / math.Pow(math.Log(2), 2))
	k = uint(math.Max(float64(1), math.Floor(0.5+float64(m)/float64(n)*math.Log(2))))
	return
}

// M returns the capacity of a Bloom filter.
func (f *Filter) M() uint {
	return f.m
}

// K returns the number of hash functions used in the Filter.
func (f *Filter) K() uint {
	return f.k
}

// B returns the filter's bitset in []uint64 format.
func (f *Filter) B() []uint64 {
	return f.bits.Bytes()
}

// Put data to the Bloom Filter. Returns the filter (allows chaining)
func (f *Filter) Put(data []byte) *Filter {
	indexes := f.s.Indexes(data, f.m, f.k)
	for _, index := range indexes {
		f.bits.Set(index)
	}
	return f
}

// MightContain query data existence.
func (f *Filter) MightContain(data []byte) bool {
	indexes := f.s.Indexes(data, f.m, f.k)
	for _, index := range indexes {
		if !f.bits.Test(index) {
			return false
		}
	}
	return true
}

// FilterJSON class wrapper for marshaling/unmarshaling Filter struct.
type FilterJSON struct {
	M uint     `json:"m"`
	K uint     `json:"k"`
	B []uint64 `json:"bits"`
}

// MarshalJSON implements json.Marshaler interface.
func (f *Filter) MarshalJSON() ([]byte, error) {
	return json.Marshal(FilterJSON{f.m, f.k, f.bits.Bytes()})
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (f *Filter) UnmarshalJSON(data []byte) error {
	var fJSON FilterJSON
	err := json.Unmarshal(data, &fJSON)
	if err != nil {
		return err
	}
	f.m = fJSON.M
	f.k = fJSON.K
	f.bits = bitset.From(fJSON.B)
	f.s = &MURMUR128MITZ64{}
	return nil
}

// SetStrategy sets strategy. Do not set it manually except only you unmarshal a bloom filter that using different strategy.
func (f *Filter) SetStrategy(s Strategy) {
	f.s = s
}
