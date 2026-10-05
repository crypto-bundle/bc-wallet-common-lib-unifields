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

package val

import (
	"errors"
	"testing"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/ptr"
	unifolderv2 "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/unifielder"
)

var errBenchmark = errors.New("benchmark error")

// ===== BenchmarkFactory* — benchmark factory functions =====

func BenchmarkString(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = String("key", "value")
	}
}

func BenchmarkInt(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = Int("key", 42)
	}
}

func BenchmarkInt8(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = Int8("key", 127)
	}
}

func BenchmarkInt16(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = Int16("key", 32767)
	}
}

func BenchmarkInt32(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = Int32("key", 999999999)
	}
}

func BenchmarkInt64(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = Int64("key", 9223372036854775807)
	}
}

func BenchmarkUint(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = Uint("key", 12345)
	}
}

func BenchmarkUint8(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = Uint8("key", 255)
	}
}

func BenchmarkUint16(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = Uint16("key", 65535)
	}
}

func BenchmarkUint32(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = Uint32("key", 4294967295)
	}
}

func BenchmarkUint64(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = Uint64("key", 18446744073709551615)
	}
}

func BenchmarkFloat32(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = Float32("key", 3.14159)
	}
}

func BenchmarkFloat64(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = Float64("key", 2.718281828)
	}
}

func BenchmarkErr(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = Err("key", errBenchmark)
	}
}

func BenchmarkTime(b *testing.B) {
	t := time.Date(2025, 6, 15, 10, 30, 45, 123456789, time.UTC)
	b.ReportAllocs()
	for range b.N {
		_ = Time("key", t)
	}
}

// ===== BenchmarkClone* — benchmark Clone method =====

func BenchmarkCloneString(b *testing.B) {
	u := String("key", "value")
	b.ReportAllocs()
	for range b.N {
		_ = u.Clone()
	}
}

func BenchmarkCloneInt(b *testing.B) {
	u := Int("key", 42)
	b.ReportAllocs()
	for range b.N {
		_ = u.Clone()
	}
}

func BenchmarkCloneFloat64(b *testing.B) {
	u := Float64("key", 2.718)
	b.ReportAllocs()
	for range b.N {
		_ = u.Clone()
	}
}

func BenchmarkCloneTime(b *testing.B) {
	t := time.Date(2025, 6, 15, 10, 30, 45, 0, time.UTC)
	u := Time("key", t)
	b.ReportAllocs()
	for range b.N {
		_ = u.Clone()
	}
}

// ===== BenchmarkMarshalTo* — benchmark MarshalTo methods =====

func BenchmarkMarshalToStr(b *testing.B) {
	u := String("key", "value")
	var s string
	b.ReportAllocs()
	for range b.N {
		_ = u.MarshalToStr(&s)
	}
}

func BenchmarkMarshalToInt(b *testing.B) {
	u := Int("key", 42)
	var i int
	b.ReportAllocs()
	for range b.N {
		_ = u.MarshalToInt(&i)
	}
}

func BenchmarkMarshalToInt8(b *testing.B) {
	u := Int8("key", 127)
	var i int8
	b.ReportAllocs()
	for range b.N {
		_ = u.MarshalToInt8(&i)
	}
}

func BenchmarkMarshalToInt16(b *testing.B) {
	u := Int16("key", 32767)
	var i int16
	b.ReportAllocs()
	for range b.N {
		_ = u.MarshalToInt16(&i)
	}
}

func BenchmarkMarshalToInt32(b *testing.B) {
	u := Int32("key", 999999999)
	var i int32
	b.ReportAllocs()
	for range b.N {
		_ = u.MarshalToInt32(&i)
	}
}

func BenchmarkMarshalToInt64(b *testing.B) {
	u := Int64("key", 9223372036854775807)
	var i int64
	b.ReportAllocs()
	for range b.N {
		_ = u.MarshalToInt64(&i)
	}
}

func BenchmarkMarshalToUint(b *testing.B) {
	u := Uint("key", 12345)
	var ui uint
	b.ReportAllocs()
	for range b.N {
		_ = u.MarshalToUint(&ui)
	}
}

func BenchmarkMarshalToUint8(b *testing.B) {
	u := Uint8("key", 255)
	var ui uint8
	b.ReportAllocs()
	for range b.N {
		_ = u.MarshalToUint8(&ui)
	}
}

func BenchmarkMarshalToUint16(b *testing.B) {
	u := Uint16("key", 65535)
	var ui uint16
	b.ReportAllocs()
	for range b.N {
		_ = u.MarshalToUint16(&ui)
	}
}

func BenchmarkMarshalToUint32(b *testing.B) {
	u := Uint32("key", 4294967295)
	var ui uint32
	b.ReportAllocs()
	for range b.N {
		_ = u.MarshalToUint32(&ui)
	}
}

func BenchmarkMarshalToUint64(b *testing.B) {
	u := Uint64("key", 18446744073709551615)
	var ui uint64
	b.ReportAllocs()
	for range b.N {
		_ = u.MarshalToUint64(&ui)
	}
}

func BenchmarkMarshalToFloat32(b *testing.B) {
	u := Float32("key", 3.14159)
	var f float32
	b.ReportAllocs()
	for range b.N {
		_ = u.MarshalToFloat32(&f)
	}
}

func BenchmarkMarshalToFloat64(b *testing.B) {
	u := Float64("key", 2.718281828)
	var f float64
	b.ReportAllocs()
	for range b.N {
		_ = u.MarshalToFloat64(&f)
	}
}

func BenchmarkMarshalToError(b *testing.B) {
	u := Err("key", errBenchmark)
	var e error
	b.ReportAllocs()
	for range b.N {
		_ = u.MarshalToError(&e)
	}
}

func BenchmarkMarshalToTime(b *testing.B) {
	t := time.Date(2025, 6, 15, 10, 30, 45, 123456789, time.UTC)
	u := Time("key", t)
	var tm time.Time
	b.ReportAllocs()
	for range b.N {
		_ = u.MarshalToTime(&tm)
	}
}

// ===== BenchmarkCollectionAdd* — benchmark adding items to Unifields collection =====

func BenchmarkCollectionAddVal_String(b *testing.B) {
	coll := NewBenchmarkCollection()
	b.ReportAllocs()
	for range b.N {
		coll.Add(String("key", "value"))
	}
}

func BenchmarkCollectionAddPtr_String(b *testing.B) {
	coll := NewBenchmarkCollection()
	b.ReportAllocs()
	for range b.N {
		coll.Add(ptr.String("key", "value"))
	}
}

func BenchmarkCollectionAddVal_Int(b *testing.B) {
	coll := NewBenchmarkCollection()
	b.ReportAllocs()
	for range b.N {
		coll.Add(Int("key", 42))
	}
}

func BenchmarkCollectionAddPtr_Int(b *testing.B) {
	coll := NewBenchmarkCollection()
	b.ReportAllocs()
	for range b.N {
		coll.Add(ptr.Int("key", 42))
	}
}

func BenchmarkCollectionAddVal_Float64(b *testing.B) {
	coll := NewBenchmarkCollection()
	b.ReportAllocs()
	for range b.N {
		coll.Add(Float64("key", 2.718))
	}
}

func BenchmarkCollectionAddPtr_Float64(b *testing.B) {
	coll := NewBenchmarkCollection()
	b.ReportAllocs()
	for range b.N {
		coll.Add(ptr.Float64("key", 2.718))
	}
}

func BenchmarkCollectionAddVal_Time(b *testing.B) {
	t := time.Date(2025, 6, 15, 10, 30, 45, 0, time.UTC)
	coll := NewBenchmarkCollection()
	b.ReportAllocs()
	for range b.N {
		coll.Add(Time("key", t))
	}
}

func BenchmarkCollectionAddPtr_Time(b *testing.B) {
	t := time.Date(2025, 6, 15, 10, 30, 45, 0, time.UTC)
	coll := NewBenchmarkCollection()
	b.ReportAllocs()
	for range b.N {
		coll.Add(ptr.Time("key", t))
	}
}

// helper type for benchmark collections -- mirrors unifields.Unifields internally.
type benchmarkCollection struct {
	items []unifolderv2.Unifielder
}

func NewBenchmarkCollection() *benchmarkCollection {
	return &benchmarkCollection{items: make([]unifolderv2.Unifielder, 0, 64)}
}

func (c *benchmarkCollection) Add(f unifolderv2.Unifielder) {
	c.items = append(c.items, f)
}
