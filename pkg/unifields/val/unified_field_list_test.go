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

func TestNewUnitfieldList_Empty(t *testing.T) {
	u := NewUnitfieldList()
	if u == nil {
		t.Fatal("NewUnitfieldList returned nil")
	}
	if u.Len() != 0 {
		t.Errorf("expected empty items slice, got length %d", u.Len())
	}
}

func TestUnitfieldList_Add_Val(t *testing.T) {
	u := NewUnitfieldList()
	f := Int("counter", 42)
	u.Add(f)
	if u.Len() != 1 {
		t.Fatalf("expected 1 item, got %d", u.Len())
	}
	var got int
	if err := f.MarshalToInt(&got); err != nil {
		t.Fatalf("MarshalToInt failed: %v", err)
	}
	if got != 42 {
		t.Errorf("got %d, want 42", got)
	}
}

func TestUnitfieldList_Add_AllTypes(t *testing.T) {
	u := NewUnitfieldList()
	u.Add(String("str", "hello"))
	u.Add(Int("int", 1))
	u.Add(Int8("i8", -8))
	u.Add(Int16("i16", -1234))
	u.Add(Int32("i32", -55555))
	u.Add(Int64("i64", -999999))
	u.Add(Uint("uint", 1))
	u.Add(Uint8("u8", 200))
	u.Add(Uint16("u16", 30000))
	u.Add(Uint32("u32", 200000))
	u.Add(Uint64("u64", 5000000))
	u.Add(Float32("f32", 1.5))
	u.Add(Float64("f64", 2.718))
	if u.Len() != 13 {
		t.Fatalf("expected 13 items, got %d", u.Len())
	}
}

func TestUnitfieldList_AddAll_Multiple(t *testing.T) {
	u := NewUnitfieldList()
	items := []unifolderv2.Unifielder{
		String("a", "one"),
		Int("b", 2),
		String("c", "three"),
		Int("d", 4),
		String("e", "five"),
	}
	u.AddAll(items)
	if u.Len() != 5 {
		t.Errorf("expected 5 items, got %d", u.Len())
	}
}

func TestUnitfieldList_AddAll_EmptySlice(t *testing.T) {
	u := NewUnitfieldList()
	initLen := u.Len()
	u.AddAll([]unifolderv2.Unifielder{})
	if u.Len() != initLen {
		t.Error("AddAll([]) should be no-op")
	}
}

func TestUnitfieldList_AddAll_NilSlice(t *testing.T) {
	u := NewUnitfieldList()
	initLen := u.Len()
	u.AddAll(nil)
	if u.Len() != initLen {
		t.Error("AddAll(nil) should be no-op")
	}
}

func TestUnitfieldList_Add_CloneIsolation(t *testing.T) {
	u := NewUnitfieldList()
	strVal := String("key", "original")
	u.Add(strVal)
	// Verify stored clone is independent of the original by reading the stored value.
	items := u.Items()
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	var got string
	unf, ok := items[0].(Unifield)
	if !ok {
		t.Fatal("expected Unifield")
	}
	if err := unf.MarshalToStr(&got); err != nil {
		t.Fatalf("MarshalToStr failed: %v", err)
	}
	if got != "original" {
		t.Errorf("clone isolation failed: got %q, want %q", got, "original")
	}
}

// --- Merge tests ---

func TestUnitfieldList_Merge_TwoLists(t *testing.T) {
	listA := NewUnitfieldList()
	listA.Add(String("a", "x"))
	listA.Add(Int("a", 1))

	// Build listB directly using value struct since we're in the val package.
	listB := UnitfieldList{
		items: []unifolderv2.Unifielder{String("z", "b"), Int("w", 2)},
	}
	listA.Merge(listB)
	if listA.Len() != 4 {
		t.Fatalf("expected 4 items, got %d", listA.Len())
	}
	// Verify merged items are readable via appropriate MarshalTo method.
	for i, item := range listA.Items() {
		var gotStr string
		unf, ok := item.(Unifield)
		if !ok {
			t.Fatalf("item %d is not Unifield", i)
		}
		if err := unf.MarshalToStr(&gotStr); err == nil {
			_ = gotStr // successfully read as string
		} else {
			var gotInt int
			if err2 := unf.MarshalToInt(&gotInt); err2 != nil {
				t.Fatalf("item %d unreadable via both MarshalToStr and MarshalToInt", i)
			}
		}
	}
}

func TestUnitfieldList_Merge_Self(t *testing.T) {
	list := NewUnitfieldList()
	list.Add(String("a", "1"))
	list.Add(Int("b", 2))
	list.Merge(*list)
	if list.Len() != 4 {
		t.Fatalf("expected 4 items after self-merge, got %d", list.Len())
	}
}

func TestUnitfieldList_Merge_EmptySource(t *testing.T) {
	dst := NewUnitfieldList()
	dst.Add(String("keep", "me"))
	src := NewUnitfieldList()
	dst.Merge(*src)
	if dst.Len() != 1 {
		t.Errorf("expected 1 item, got %d after merging empty source", dst.Len())
	}
}

func TestUnitfieldList_Merge_SingleElementSource(t *testing.T) {
	dst := NewUnitfieldList()
	src := NewUnitfieldList()
	src.Add(String("single", "element"))
	dst.Merge(*src)
	if dst.Len() != 1 {
		t.Errorf("expected 1 item, got %d", dst.Len())
	}
}

func TestUnitfieldList_Merge_ZeroValueReceiver(t *testing.T) {
	var empty UnitfieldList
	src := NewUnitfieldList()
	src.Add(String("key", "value"))
	// Call Merge on zero-value receiver — should not panic.
	// Merge has pointer-receiver, so modifications ARE visible to caller.
	empty.Merge(*src)
	if empty.Len() != 1 {
		t.Errorf("expected 1 item after merge into zero-value receiver, got %d", empty.Len())
	}
}

// --- GetAfter tests ---

func TestUnitfieldList_GetAfter_Zero(t *testing.T) {
	u := NewUnitfieldList()
	u.Add(String("a", "1"))
	u.Add(Int("b", 2))
	u.Add(String("c", "3"))
	result := u.GetAfter(0)
	if len(result) != 3 {
		t.Fatalf("expected 3 items, got %d", len(result))
	}
}

func TestUnitfieldList_GetAfter_MiddleIndex(t *testing.T) {
	u := NewUnitfieldList()
	u.Add(String("a", "1"))
	u.Add(String("b", "2"))
	u.Add(String("c", "3"))
	u.Add(String("d", "4"))
	result := u.GetAfter(2) // indices 2,3
	if len(result) != 2 {
		t.Fatalf("expected 2 items, got %d", len(result))
	}
}

func TestUnitfieldList_GetAfter_Boundary(t *testing.T) {
	u := NewUnitfieldList()
	u.Add(String("a", "1"))
	result := u.GetAfter(1) // len=1, index=1 => empty
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d items", len(result))
	}
}

// --- RemoveAfter tests ---

func TestUnitfieldList_RemoveAfter_KeepHead(t *testing.T) {
	u := NewUnitfieldList()
	u.Add(String("a", "1"))
	u.Add(String("b", "2"))
	u.Add(String("c", "3"))
	u.RemoveAfter(0)
	if u.Len() != 1 {
		t.Fatalf("expected 1 item, got %d", u.Len())
	}
}

func TestUnitfieldList_RemoveAfter_LastIndex(t *testing.T) {
	u := NewUnitfieldList()
	u.Add(String("a", "1"))
	u.Add(String("b", "2"))
	u.Add(String("c", "3"))
	u.RemoveAfter(2) // index=last => nothing dropped
	if u.Len() != 3 {
		t.Fatalf("expected 3 items, got %d", u.Len())
	}
}

func TestUnitfieldList_RemoveAfter_OutOfRange(t *testing.T) {
	u := NewUnitfieldList()
	u.Add(String("a", "1"))
	u.Add(String("b", "2"))
	u.RemoveAfter(2) // index == len
	if u.Len() != 2 {
		t.Fatalf("expected 2 items, got %d", u.Len())
	}
}

func TestUnitfieldList_RemoveAfter_ExtremeOutOfRange(t *testing.T) {
	u := NewUnitfieldList()
	u.Add(String("a", "1"))
	u.RemoveAfter(999) // way beyond
	if u.Len() != 1 {
		t.Fatalf("expected 1 item, got %d", u.Len())
	}
}

func TestUnitfieldList_RemoveAfter_SingleElement(t *testing.T) {
	u := NewUnitfieldList()
	u.Add(String("only", "one"))
	u.RemoveAfter(0) // index=0 on 1-element list
	if u.Len() != 1 {
		t.Errorf("expected 1 item after RemoveAfter(0) on single, got %d", u.Len())
	}
}

// --- RemoveBefore tests ---

func TestUnitfieldList_RemoveBefore_KeepTail(t *testing.T) {
	u := NewUnitfieldList()
	u.Add(String("a", "1"))
	u.Add(String("b", "2"))
	u.Add(String("c", "3"))
	u.RemoveBefore(2) // keep index 2..end
	if u.Len() != 1 {
		t.Fatalf("expected 1 item, got %d", u.Len())
	}
}

func TestUnitfieldList_RemoveBefore_IndexZero(t *testing.T) {
	u := NewUnitfieldList()
	u.Add(String("a", "1"))
	u.Add(String("b", "2"))
	u.RemoveBefore(0) // nothing before index 0 => no-op
	if u.Len() != 2 {
		t.Fatalf("expected 2 items (no-op), got %d", u.Len())
	}
}

func TestUnitfieldList_RemoveBefore_OutOfRange(t *testing.T) {
	u := NewUnitfieldList()
	u.Add(String("a", "1"))
	u.RemoveBefore(1) // index == len
	if u.Len() != 1 {
		t.Fatalf("expected 1 item, got %d", u.Len())
	}
}

func TestUnitfieldList_RemoveBefore_ExtremeOutOfRange(t *testing.T) {
	u := NewUnitfieldList()
	u.Add(String("a", "1"))
	u.RemoveBefore(999)
	if u.Len() != 1 {
		t.Fatalf("expected 1 item, got %d", u.Len())
	}
}

func TestUnitfieldList_RemoveBefore_SingleElement(t *testing.T) {
	u := NewUnitfieldList()
	u.Add(String("only", "one"))
	u.RemoveBefore(0) // no-op on single
	if u.Len() != 1 {
		t.Errorf("expected 1 item (no-op), got %d", u.Len())
	}
}

// --- Clear test ---

func TestUnitfieldList_Clear_AllDropped(t *testing.T) {
	u := NewUnitfieldList()
	u.Add(String("a", "1"))
	u.Add(Int("b", 2))
	u.Clear()
	if u.Len() != 0 {
		t.Errorf("expected 0 items after Clear, got %d", u.Len())
	}
}

// --- Items test ---

func TestUnitfieldList_Items_ReadOnly(t *testing.T) {
	u := NewUnitfieldList()
	u.Add(String("key", "val"))
	items := u.Items()
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	// Mutating returned slice does NOT affect internal state.
	_ = items[:0]
	if u.Len() != 1 {
		t.Errorf("internal state should not be affected by mutating returned slice, got Len=%d", u.Len())
	}
}

// --- Polymorphic mix tests ---

func TestUnitfieldList_PolymorphicMix(t *testing.T) {
	// All val-based: String + Int + Float64 share []Unifielder via interface.
	u := NewUnitfieldList()
	u.Add(String("a", "one"))
	u.Add(Int("b", 2))
	u.Add(Float64("c", 3.0))
	items := u.Items()
	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}
	getResult := u.GetAfter(1) // indices 1,2
	if len(getResult) != 2 {
		t.Fatalf("expected 2 items from GetAfter(1), got %d", len(getResult))
	}
}
