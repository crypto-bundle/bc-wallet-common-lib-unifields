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

package val

import (
	"testing"
	"time"

	unifolderv2 "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/unifielder"
)

// BenchmarkNewUnitfieldStack measures constructor allocation overhead.
func BenchmarkNewUnitfieldStack(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = NewUnitfieldStack()
	}
}

// BenchmarkPushSingle measures single-element push throughput.
func BenchmarkPushSingle(b *testing.B) {
	b.ReportAllocs()
	fld := String("key", "value")
	for range b.N {
		s := NewUnitfieldStack()
		s.Push(fld)
	}
}

// BenchmarkPushLarge measures bulk push throughput for many elements.
func BenchmarkPushLarge(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		s := NewUnitfieldStack()
		s.PushFields(
			String("a", "alpha"), Int("b", 1), Float64("c", 1.0),
			String("d", "delta"), Uint("e", 5), Time("f", time.Now()),
		)
	}
}

// BenchmarkPushFieldsVariadic measures variadic PushFields throughput.
func BenchmarkPushFieldsVariadic(b *testing.B) {
	b.ReportAllocs()
	fields := []unifolderv2.Unifielder{
		String("one", "1"), Int("two", 2), Float64("three", 3.0),
		Uint("four", 4), String("five", "5"),
	}
	for range b.N {
		s := NewUnitfieldStack()
		s.PushFields(fields...)
	}
}

// BenchmarkPopSingle measures single-element pop throughput.
func BenchmarkPopSingle(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		s := NewUnitfieldStack()
		s.Push(String("x", "val"))
		_ = s.Pop()
	}
}

// BenchmarkPopNBatch measures batch pop throughput.
func BenchmarkPopNBatch(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		s := NewUnitfieldStack()
		s.PushFields(String("a", "1"), Int("b", 2), Float64("c", 3.0))
		_ = s.PopN(2)
	}
}

// BenchmarkGetTopPeek measures read-only peek throughput (no mutation).
func BenchmarkGetTopPeek(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		s := NewUnitfieldStack()
		s.Push(Int("n", 100))
		_ = s.GetTop()
	}
}

// BenchmarkClearThroughput measures Clear operation throughput.
func BenchmarkClearThroughput(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		s := NewUnitfieldStack()
		s.PushFields(String("a", "1"), Int("b", 2), Float64("c", 3.0))
		s.Clear()
	}
}

// BenchmarkLenAccess measures Len() call overhead.
func BenchmarkLenAccess(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		s := NewUnitfieldStack()
		s.PushFields(
			String("a", "1"), Int("b", 2), Float64("c", 3.0),
			Uint("d", 4), Err("e", ErrTypeMismatch),
		)
		_ = s.Len()
	}
}

// BenchmarkPopN_Large stresses PopN with a large stack.
func BenchmarkPopN_Large(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		s := NewUnitfieldStack()
		for j := 0; j < 1000; j++ {
			s.Push(Int("k", j))
		}
		_ = s.PopN(500)
	}
}
