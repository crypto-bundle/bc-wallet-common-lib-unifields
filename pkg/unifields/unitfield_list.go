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

	valpkg "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/val"
)

// UnitfieldList is the recommended immutable collection supporting both val.Unifield
// and ptr.UnifieldPtr. It embeds a [*valpkg.UnitfieldList] and delegates all collection
// operations to it. Use [NewUnitfieldList] to create a new instance.
type UnitfieldList struct {
	*valpkg.UnitfieldList
}

// NewUnitfieldList creates a new empty UnitfieldList collection.
// It wraps a [*valpkg.UnitfieldList] and exposes the full API from the parent package,
// so consumers never need to import the val/ sub-package directly.
func NewUnitfieldList() *UnitfieldList {
	return &UnitfieldList{UnitfieldList: valpkg.NewUnitfieldList()}
}

// --- Typed adders for UnitfieldList ---

// AddStr adds a string-typed Unifield with the given key and value.
func (u *UnitfieldList) AddStr(key string, value string) {
	u.Add(valpkg.String(key, value))
}

// AddInt adds an int-typed Unifield with the given key and value.
func (u *UnitfieldList) AddInt(key string, value int) {
	u.Add(valpkg.Int(key, value))
}

// AddInt8 adds an int8-typed Unifield with the given key and value.
func (u *UnitfieldList) AddInt8(key string, value int8) {
	u.Add(valpkg.Int8(key, value))
}

// AddInt16 adds an int16-typed Unifield with the given key and value.
func (u *UnitfieldList) AddInt16(key string, value int16) {
	u.Add(valpkg.Int16(key, value))
}

// AddInt32 adds an int32-typed Unifield with the given key and value.
func (u *UnitfieldList) AddInt32(key string, value int32) {
	u.Add(valpkg.Int32(key, value))
}

// AddInt64 adds an int64-typed Unifield with the given key and value.
func (u *UnitfieldList) AddInt64(key string, value int64) {
	u.Add(valpkg.Int64(key, value))
}

// AddUint adds a uint-typed Unifield with the given key and value.
func (u *UnitfieldList) AddUint(key string, value uint) {
	u.Add(valpkg.Uint(key, value))
}

// AddUint8 adds a uint8-typed Unifield with the given key and value.
func (u *UnitfieldList) AddUint8(key string, value uint8) {
	u.Add(valpkg.Uint8(key, value))
}

// AddUint16 adds a uint16-typed Unifield with the given key and value.
func (u *UnitfieldList) AddUint16(key string, value uint16) {
	u.Add(valpkg.Uint16(key, value))
}

// AddUint32 adds a uint32-typed Unifield with the given key and value.
func (u *UnitfieldList) AddUint32(key string, value uint32) {
	u.Add(valpkg.Uint32(key, value))
}

// AddUint64 adds a uint64-typed Unifield with the given key and value.
func (u *UnitfieldList) AddUint64(key string, value uint64) {
	u.Add(valpkg.Uint64(key, value))
}

// AddFloat32 adds a float32-typed Unifield with the given key and value.
func (u *UnitfieldList) AddFloat32(key string, value float32) {
	u.Add(valpkg.Float32(key, value))
}

// AddFloat64 adds a float64-typed Unifield with the given key and value.
func (u *UnitfieldList) AddFloat64(key string, value float64) {
	u.Add(valpkg.Float64(key, value))
}

// AddErr adds an error-typed Unifield with the given key and value.
func (u *UnitfieldList) AddErr(key string, value error) {
	u.Add(valpkg.Err(key, value))
}

// AddTime adds a time.Time-typed Unifield with the given key and value.
func (u *UnitfieldList) AddTime(key string, value time.Time) {
	u.Add(valpkg.Time(key, value))
}
