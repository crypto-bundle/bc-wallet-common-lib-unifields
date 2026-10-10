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

	unifolderv2 "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/unifielder"
)

func BenchmarkUnitfieldList_New(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = NewUnitfieldList()
	}
}

func BenchmarkUnitfieldList_Add_SingleVal(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		u := NewUnitfieldList()
		u.Add(String("value"))
	}
}

func BenchmarkUnitfieldList_Add_SingleInt(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		u := NewUnitfieldList()
		u.Add(Int(42))
	}
}

func BenchmarkUnitfieldList_AddAll_10(b *testing.B) {
	b.ReportAllocs()
	items := []unifolderv2.Unifielder{
		String("1"), Int(2), String("3"),
		Int(4), String("5"), Uint(6),
		Uint64(7), Float32(8.1), Float64(9.9),
		String("last"),
	}
	for range b.N {
		u := NewUnitfieldList()
		u.AddAll(items)
	}
}

func BenchmarkUnitfieldList_AddAll_1000(b *testing.B) {
	b.ReportAllocs()
	items := make([]unifolderv2.Unifielder, 1000)
	for j := range 1000 {
		items[j] = String(string(rune(j)))
	}
	for range b.N {
		u := NewUnitfieldList()
		u.AddAll(items)
	}
}

func BenchmarkUnitfieldList_Merge_TwoMediumLists(b *testing.B) {
	b.ReportAllocs()
	srcA := newUnitfieldListValue(50, "a")
	srcB := newUnitfieldListValue(50, "b")
	for range b.N {
		dst := NewUnitfieldList()
		dst.Merge(srcA)
		dst.Merge(srcB)
		_ = dst.Len()
	}
}

func BenchmarkUnitfieldList_GetAfter_HalfSlice(b *testing.B) {
	b.ReportAllocs()
	list := newUnitfieldListValue(100, "m")
	half := list.Len() / 2
	for range b.N {
		result := list.GetAfter(uint(half))
		_ = len(result)
	}
}

func BenchmarkUnitfieldList_RemoveAfter_TailDrop(b *testing.B) {
	b.ReportAllocs()
	const size = 100
	splitIdx := size / 2
	for range b.N {
		tmp := newUnitfieldListValue(size, "o")
		tmp.RemoveAfter(uint(splitIdx))
	}
}

func BenchmarkUnitfieldList_RemoveBefore_HeadDrop(b *testing.B) {
	b.ReportAllocs()
	const size = 100
	splitIdx := size / 2
	for range b.N {
		tmp := newUnitfieldListValue(size, "q")
		tmp.RemoveBefore(uint(splitIdx))
	}
}

func BenchmarkUnitfieldList_Clear_FastReset(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		u := NewUnitfieldList()
		for range 200 {
			u.Add(String(string(rune('a'+len(u.items)%26))))
		}
		u.Clear()
	}
}

// newUnitfieldListValue creates a UnitfieldList with count items using the given prefix.
func newUnitfieldListValue(count int, prefix string) *UnitfieldList {
	u := NewUnitfieldList()
	for range count {
		u.Add(String(prefix))
	}
	return u
}
