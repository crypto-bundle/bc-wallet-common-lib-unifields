// MIT NON-AI License
//
// Copyright (c) 2022-2026 Aleksei Kotelnikov(gudron2s@gmail.com)
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

import "time"

// Unifields is an immutable collection of Unifield values stored as pointers internally.
// External users cannot mutate internal state — all Add operations store clones,
// ensuring that modifying a returned Unifield does not affect the collection's contents.
type Unifields struct {
	items []*Unifield
}

// NewUnifields creates a new empty Unifields collection.
func NewUnifields() *Unifields {
	return &Unifields{items: make([]*Unifield, 0)}
}

// Add appends a clone of the given Unifield to the collection.
func (u *Unifields) Add(fld Unifield) {
	c := fld.Clone()
	u.items = append(u.items, &c)
}

// AddAll appends clones of all given Unifields to the collection.
// If flds is empty or nil, this is a no-op.
func (u *Unifields) AddAll(flds []Unifield) {
	if len(flds) == 0 {
		return
	}
	cloned := make([]*Unifield, len(flds))
	for i := range flds {
		c := flds[i].Clone()
		cloned[i] = &c
	}
	u.items = append(u.items, cloned...)
}

// AddStr adds a string-typed Unifield with the given key and value.
func (u *Unifields) AddStr(key string, val string) {
	u.Add(String(key, val))
}

// AddInt adds an int-typed Unifield with the given key and value.
func (u *Unifields) AddInt(key string, val int) {
	u.Add(Int(key, val))
}

// AddInt8 adds an int8-typed Unifield with the given key and value.
func (u *Unifields) AddInt8(key string, val int8) {
	u.Add(Int8(key, val))
}

// AddInt16 adds an int16-typed Unifield with the given key and value.
func (u *Unifields) AddInt16(key string, val int16) {
	u.Add(Int16(key, val))
}

// AddInt32 adds an int32-typed Unifield with the given key and value.
func (u *Unifields) AddInt32(key string, val int32) {
	u.Add(Int32(key, val))
}

// AddInt64 adds an int64-typed Unifield with the given key and value.
func (u *Unifields) AddInt64(key string, val int64) {
	u.Add(Int64(key, val))
}

// AddUint adds a uint-typed Unifield with the given key and value.
func (u *Unifields) AddUint(key string, val uint) {
	u.Add(Uint(key, val))
}

// AddUint8 adds a uint8-typed Unifield with the given key and value.
func (u *Unifields) AddUint8(key string, val uint8) {
	u.Add(Uint8(key, val))
}

// AddUint16 adds a uint16-typed Unifield with the given key and value.
func (u *Unifields) AddUint16(key string, val uint16) {
	u.Add(Uint16(key, val))
}

// AddUint32 adds a uint32-typed Unifield with the given key and value.
func (u *Unifields) AddUint32(key string, val uint32) {
	u.Add(Uint32(key, val))
}

// AddUint64 adds a uint64-typed Unifield with the given key and value.
func (u *Unifields) AddUint64(key string, val uint64) {
	u.Add(Uint64(key, val))
}

// AddFloat32 adds a float32-typed Unifield with the given key and value.
func (u *Unifields) AddFloat32(key string, val float32) {
	u.Add(Float32(key, val))
}

// AddFloat64 adds a float64-typed Unifield with the given key and value.
func (u *Unifields) AddFloat64(key string, val float64) {
	u.Add(Float64(key, val))
}

// AddErr adds an error-typed Unifield with the given key and value.
func (u *Unifields) AddErr(key string, val error) {
	u.Add(Err(key, val))
}

// AddTime adds a time.Time-typed Unifield with the given key and value.
func (u *Unifields) AddTime(key string, val time.Time) {
	u.Add(Time(key, val))
}
