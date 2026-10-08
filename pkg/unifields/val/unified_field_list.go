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
	"time"

	unifolderv2 "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/unifielder"
)

// UnitfieldList is an immutable collection supporting both val.Unifield and ptr.UnifieldPtr
// via the shared [Unifielder] interface. External users cannot mutate internal state — all Add
// operations store clones, ensuring that modifying a returned Unifielder does not affect
// the collection's contents.
type UnitfieldList struct {
	items []unifolderv2.Unifielder
}

// NewUnitfieldList creates a new empty UnitfieldList with zero-length items slice.
func NewUnitfieldList() *UnitfieldList {
	return &UnitfieldList{items: make([]unifolderv2.Unifielder, 0)}
}

// Add appends a clone of the given Unifielder to the collection.
// Accepts both val.Unifield and ptr.UnifieldPtr implementations.
func (u *UnitfieldList) Add(fld unifolderv2.Unifielder) {
	c := fld.Clone()
	u.items = append(u.items, c)
}

// AddAll appends clones of all given Unifielders to the collection.
// If flds is empty or nil, this is a no-op.
func (u *UnitfieldList) AddAll(flds []unifolderv2.Unifielder) {
	if len(flds) == 0 {
		return
	}
	cloned := make([]unifolderv2.Unifielder, len(flds))
	for i := range flds {
		c := flds[i].Clone()
		cloned[i] = c
	}
	u.items = append(u.items, cloned...)
}

// Merge appends clones of all items from the source list into this list.
// Each item is cloned individually to match the invariant of Add/AddAll —
// stored references are independent of the source list. Keys may not be unique
// across the merged result. No key-collision semantics (deduplication) is applied.
func (u *UnitfieldList) Merge(list UnitfieldList) {
	if len(list.items) == 0 {
		return
	}
	for _, item := range list.items {
		u.items = append(u.items, item.Clone())
	}
}

// GetAfter returns a read-only copy of items starting at the given index (inclusive).
// Returns an empty slice when index >= len(items); never panics on out-of-range indices.
func (u *UnitfieldList) GetAfter(index uint) []unifolderv2.Unifielder {
	n := len(u.items)
	idx := int(index)
	if idx < 0 || idx >= n {
		return nil
	}
	out := make([]unifolderv2.Unifielder, n-idx)
	copy(out, u.items[idx:])
	return out
}

// RemoveAfter removes all items after the element at the given index (keeping index itself).
// For example: RemoveAfter(0) on [a,b,c] keeps only [a].
// Out-of-range indices (index >= len) are a no-op; never panics.
func (u *UnitfieldList) RemoveAfter(index uint) {
	n := len(u.items)
	idx := int(index)
	if idx < 0 || idx >= n {
		return
	}
	u.items = u.items[:idx+1]
}

// RemoveBefore removes all items before the element at the given index (keeping index itself).
// For example: RemoveBefore(2) on [a,b,c,d] keeps only [c,d].
// When index is 0 or out-of-range, this is a no-op. Never panics.
func (u *UnitfieldList) RemoveBefore(index uint) {
	n := len(u.items)
	idx := int(index)
	if idx <= 0 || idx >= n {
		return
	}
	u.items = u.items[idx:]
}

// Clear removes all items from the list, resetting the internal slice to empty.
// Subsequent calls to Len() return 0.
func (u *UnitfieldList) Clear() {
	u.items = u.items[:0]
}

// Len returns the number of items currently stored in the list.
func (u *UnitfieldList) Len() int {
	return len(u.items)
}

// Items returns a read-only copy of all stored Unifielders.
// Modifying the returned slice has no effect on the list's internal state.
func (u *UnitfieldList) Items() []unifolderv2.Unifielder {
	out := make([]unifolderv2.Unifielder, len(u.items))
	copy(out, u.items)
	return out
}

// AddStr adds a string-typed Unifield with the given key and value.
func (u *UnitfieldList) AddStr(key string, value string) {
	u.Add(String(key, value))
}

// AddInt adds an int-typed Unifield with the given key and value.
func (u *UnitfieldList) AddInt(key string, value int) {
	u.Add(Int(key, value))
}

// AddInt8 adds an int8-typed Unifield with the given key and value.
func (u *UnitfieldList) AddInt8(key string, value int8) {
	u.Add(Int8(key, value))
}

// AddInt16 adds an int16-typed Unifield with the given key and value.
func (u *UnitfieldList) AddInt16(key string, value int16) {
	u.Add(Int16(key, value))
}

// AddInt32 adds an int32-typed Unifield with the given key and value.
func (u *UnitfieldList) AddInt32(key string, value int32) {
	u.Add(Int32(key, value))
}

// AddInt64 adds an int64-typed Unifield with the given key and value.
func (u *UnitfieldList) AddInt64(key string, value int64) {
	u.Add(Int64(key, value))
}

// AddUint adds a uint-typed Unifield with the given key and value.
func (u *UnitfieldList) AddUint(key string, value uint) {
	u.Add(Uint(key, value))
}

// AddUint8 adds a uint8-typed Unifield with the given key and value.
func (u *UnitfieldList) AddUint8(key string, value uint8) {
	u.Add(Uint8(key, value))
}

// AddUint16 adds a uint16-typed Unifield with the given key and value.
func (u *UnitfieldList) AddUint16(key string, value uint16) {
	u.Add(Uint16(key, value))
}

// AddUint32 adds a uint32-typed Unifield with the given key and value.
func (u *UnitfieldList) AddUint32(key string, value uint32) {
	u.Add(Uint32(key, value))
}

// AddUint64 adds a uint64-typed Unifield with the given key and value.
func (u *UnitfieldList) AddUint64(key string, value uint64) {
	u.Add(Uint64(key, value))
}

// AddFloat32 adds a float32-typed Unifield with the given key and value.
func (u *UnitfieldList) AddFloat32(key string, value float32) {
	u.Add(Float32(key, value))
}

// AddFloat64 adds a float64-typed Unifield with the given key and value.
func (u *UnitfieldList) AddFloat64(key string, value float64) {
	u.Add(Float64(key, value))
}

// AddErr adds an error-typed Unifield with the given key and value.
func (u *UnitfieldList) AddErr(key string, value error) {
	u.Add(Err(key, value))
}

// AddTime adds a time.Time-typed Unifield with the given key and value.
func (u *UnitfieldList) AddTime(key string, value time.Time) {
	u.Add(Time(key, value))
}
