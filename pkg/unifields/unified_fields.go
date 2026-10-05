// MIT NON-AI License
//
// Copyright (c) 2024-2026 Aleksei Kotelnikov(gudron2s@gmail.com)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
//
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
//
// In addition, the following restrictions apply:
//
// 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
// including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
// modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
//
// 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
// including but not limited to artificial intelligence, natural language processing, or data mining.
//
// 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
// for any damages resulting from such use.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package unifields

import (
	"time"

	unifolderv2 "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/unifielder"
	valpkg "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/val"
)

// Unifielder is the shared interface implemented by both val.Unifield and ptr.UnifieldPtr.
type Unifielder = unifolderv2.Unifielder

// Unifields is an immutable collection supporting both val.Unifield and ptr.UnifieldPtr.
// External users cannot mutate internal state — all Add operations store clones,
// ensuring that modifying a returned Unifield does not affect the collection's contents.
type Unifields struct {
	items []Unifielder
}

// NewUnifields creates a new empty Unifields collection.
func NewUnifields() *Unifields {
	return &Unifields{items: make([]Unifielder, 0)}
}

// Add appends a clone of the given Unifielder to the collection.
// Accepts both val.Unifield and ptr.UnifieldPtr implementations.
func (u *Unifields) Add(fld Unifielder) {
	c := fld.Clone()
	u.items = append(u.items, c)
}

// AddAll appends clones of all given Unifielders to the collection.
// If flds is empty or nil, this is a no-op.
func (u *Unifields) AddAll(flds []Unifielder) {
	if len(flds) == 0 {
		return
	}
	cloned := make([]Unifielder, len(flds))
	for i := range flds {
		c := flds[i].Clone()
		cloned[i] = c
	}
	u.items = append(u.items, cloned...)
}

// AddStr adds a string-typed Unifield with the given key and value.
func (u *Unifields) AddStr(key string, value string) {
	u.Add(valpkg.String(key, value))
}

// AddInt adds an int-typed Unifield with the given key and value.
func (u *Unifields) AddInt(key string, value int) {
	u.Add(valpkg.Int(key, value))
}

// AddInt8 adds an int8-typed Unifield with the given key and value.
func (u *Unifields) AddInt8(key string, value int8) {
	u.Add(valpkg.Int8(key, value))
}

// AddInt16 adds an int16-typed Unifield with the given key and value.
func (u *Unifields) AddInt16(key string, value int16) {
	u.Add(valpkg.Int16(key, value))
}

// AddInt32 adds an int32-typed Unifield with the given key and value.
func (u *Unifields) AddInt32(key string, value int32) {
	u.Add(valpkg.Int32(key, value))
}

// AddInt64 adds an int64-typed Unifield with the given key and value.
func (u *Unifields) AddInt64(key string, value int64) {
	u.Add(valpkg.Int64(key, value))
}

// AddUint adds a uint-typed Unifield with the given key and value.
func (u *Unifields) AddUint(key string, value uint) {
	u.Add(valpkg.Uint(key, value))
}

// AddUint8 adds a uint8-typed Unifield with the given key and value.
func (u *Unifields) AddUint8(key string, value uint8) {
	u.Add(valpkg.Uint8(key, value))
}

// AddUint16 adds an int16-typed Unifield with the given key and value.
func (u *Unifields) AddUint16(key string, value uint16) {
	u.Add(valpkg.Uint16(key, value))
}

// AddUint32 adds an int32-typed Unifield with the given key and value.
func (u *Unifields) AddUint32(key string, value uint32) {
	u.Add(valpkg.Uint32(key, value))
}

// AddUint64 adds a uint64-typed Unifield with the given key and value.
func (u *Unifields) AddUint64(key string, value uint64) {
	u.Add(valpkg.Uint64(key, value))
}

// AddFloat32 adds a float32-typed Unifield with the given key and value.
func (u *Unifields) AddFloat32(key string, value float32) {
	u.Add(valpkg.Float32(key, value))
}

// AddFloat64 adds a float64-typed Unifield with the given key and value.
func (u *Unifields) AddFloat64(key string, value float64) {
	u.Add(valpkg.Float64(key, value))
}

// AddErr adds an error-typed Unifield with the given key and value.
func (u *Unifields) AddErr(key string, value error) {
	u.Add(valpkg.Err(key, value))
}

// AddTime adds a time.Time-typed Unifield with the given key and value.
func (u *Unifields) AddTime(key string, value time.Time) {
	u.Add(valpkg.Time(key, value))
}

// Len returns the number of items in the collection.
func (u *Unifields) Len() int {
	return len(u.items)
}

// Items returns a read-only slice of all stored Unifielders.
// Modifying returned items has no effect on the collection's internal state.
func (u *Unifields) Items() []Unifielder {
	out := make([]Unifielder, len(u.items))
	copy(out, u.items)
	return out
}
