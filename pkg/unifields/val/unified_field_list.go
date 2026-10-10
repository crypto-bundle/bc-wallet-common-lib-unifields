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
	"sync"
	"time"

	unifolderv2 "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/unifielder"
)

// UnitfieldList is an immutable collection supporting both val.Unifield and ptr.UnifieldPtr
// via the shared [Unifielder] interface. External users cannot mutate internal state — all Add
// operations store clones, ensuring that modifying a returned Unifielder does not affect
// the collection's contents.
type UnitfieldList struct {
	items []unifolderv2.Unifielder
	mu    sync.RWMutex
}

// NewUnitfieldList creates a new empty UnitfieldList with zero-length items slice.
func NewUnitfieldList() *UnitfieldList {
	return &UnitfieldList{items: make([]unifolderv2.Unifielder, 0)}
}

// Add appends a clone of the given Unifielder to the collection.
// Accepts both val.Unifield and ptr.UnifieldPtr implementations.
func (u *UnitfieldList) Add(fld unifolderv2.Unifielder) {
	u.mu.Lock()
	defer u.mu.Unlock()
	c := fld.Clone()
	u.items = append(u.items, c)
}

// AddAll appends clones of all given Unifielders to the collection.
// If flds is empty or nil, this is a no-op.
func (u *UnitfieldList) AddAll(flds []unifolderv2.Unifielder) {
	u.mu.Lock()
	defer u.mu.Unlock()
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

// Merge appends clones of all items from source into this list.
// Each item is cloned individually to match the invariant of Add/AddAll —
// stored references are independent of the source list. Keys may not be unique
// across the merged result. No key-collision semantics (deduplication) is applied.
// Merge(nil) is a safe no-op that returns immediately without modifying either list.
func (u *UnitfieldList) Merge(source *UnitfieldList) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if source != nil {
		source.mu.RLock()
		defer source.mu.RUnlock()
		u.mergeNoLock(source)
	}
}

// mergeNoLock appends clones of all items from source into this list.
// Caller must hold mu.Lock on u and mu.RLock on source.
func (u *UnitfieldList) mergeNoLock(source *UnitfieldList) {
	for _, item := range source.items {
		u.items = append(u.items, item.Clone())
	}
}

// GetAfter returns a read-only copy of items starting at the given index (inclusive).
// Returns an empty slice when index >= len(items); never panics on out-of-range indices.
func (u *UnitfieldList) GetAfter(index uint) []unifolderv2.Unifielder {
	u.mu.RLock()
	defer u.mu.RUnlock()
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
// Items removed are returned to the pool.
func (u *UnitfieldList) RemoveAfter(index uint) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.removeAfterNoLock(int(index))
}

// removeAfterNoLock removes all items after the element at the given index (keeping index itself).
// For example: RemoveAfter(0) on [a,b,c] keeps only [a].
// Out-of-range indices (index >= len) are a no-op; never panics.
// Items removed are returned to the pool. Not safe externally.
func (u *UnitfieldList) removeAfterNoLock(index int) {
	n := len(u.items)
	idx := index
	if idx < 0 || idx >= n {
		return
	}
	for i := idx + 1; i < n; i++ {
		returnValToPool(u.items[i])
	}
	u.items = u.items[:idx+1]
}

// RemoveBefore removes all items before the element at the given index (keeping index itself).
// For example: RemoveBefore(2) on [a,b,c,d] keeps only [c,d].
// When index is 0 or out-of-range, this is a no-op. Never panics.
// Items removed are returned to the pool.
func (u *UnitfieldList) RemoveBefore(index uint) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.removeBeforeNoLock(int(index))
}

// removeBeforeNoLock removes all items before the element at the given index (keeping index itself).
// For example: RemoveBefore(2) on [a,b,c,d] keeps only [c,d].
// When index is 0 or out-of-range, this is a no-op. Never panics.
// Items removed are returned to the pool. Not safe externally.
func (u *UnitfieldList) removeBeforeNoLock(index int) {
	n := len(u.items)
	idx := index
	if idx <= 0 || idx >= n {
		return
	}
	for _, item := range u.items[:idx] {
		returnValToPool(item)
	}
	u.items = u.items[idx:]
}

// clearNoLock clears all items from the list and returns them to the pool.
// Caller must hold u.mu.Lock(). Not safe externally.
func (u *UnitfieldList) clearNoLock() {
	for _, item := range u.items {
		returnValToPool(item)
	}
	u.items = u.items[:0]
}

// Clear removes all items from the list, resetting the internal slice to empty.
// Subsequent calls to Len() return 0.
// Items are returned to the pool.
func (u *UnitfieldList) Clear() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.clearNoLock()
}

// Len returns the number of items currently stored in the list.
func (u *UnitfieldList) Len() int {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return len(u.items)
}

// Items returns a defensive copy of all stored Unifielders.
// Modifying the returned slice has no effect on the list's internal state.
func (u *UnitfieldList) Items() []unifolderv2.Unifielder {
	u.mu.RLock()
	defer u.mu.RUnlock()
	out := make([]unifolderv2.Unifielder, len(u.items))
	copy(out, u.items)
	return out
}

// AddStr adds a string-typed Unifield to the list.
func (u *UnitfieldList) AddStr(value string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.items = append(u.items, String(value))
}

// AddInt adds an int-typed Unifield to the list.
func (u *UnitfieldList) AddInt(value int) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.items = append(u.items, Int(value))
}

// AddInt8 adds an int8-typed Unifield to the list.
func (u *UnitfieldList) AddInt8(value int8) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.items = append(u.items, Int8(value))
}

// AddInt16 adds an int16-typed Unifield to the list.
func (u *UnitfieldList) AddInt16(value int16) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.items = append(u.items, Int16(value))
}

// AddInt32 adds an int32-typed Unifield to the list.
func (u *UnitfieldList) AddInt32(value int32) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.items = append(u.items, Int32(value))
}

// AddInt64 adds an int64-typed Unifield to the list.
func (u *UnitfieldList) AddInt64(value int64) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.items = append(u.items, Int64(value))
}

// AddUint adds a uint-typed Unifield to the list.
func (u *UnitfieldList) AddUint(value uint) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.items = append(u.items, Uint(value))
}

// AddUint8 adds a uint8-typed Unifield to the list.
func (u *UnitfieldList) AddUint8(value uint8) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.items = append(u.items, Uint8(value))
}

// AddUint16 adds a uint16-typed Unifield to the list.
func (u *UnitfieldList) AddUint16(value uint16) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.items = append(u.items, Uint16(value))
}

// AddUint32 adds a uint32-typed Unifield to the list.
func (u *UnitfieldList) AddUint32(value uint32) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.items = append(u.items, Uint32(value))
}

// AddUint64 adds a uint64-typed Unifield to the list.
func (u *UnitfieldList) AddUint64(value uint64) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.items = append(u.items, Uint64(value))
}

// AddFloat32 adds a float32-typed Unifield to the list.
func (u *UnitfieldList) AddFloat32(value float32) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.items = append(u.items, Float32(value))
}

// AddFloat64 adds a float64-typed Unifield to the list.
func (u *UnitfieldList) AddFloat64(value float64) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.items = append(u.items, Float64(value))
}

// AddErr adds an error-typed Unifield to the list.
func (u *UnitfieldList) AddErr(value error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.items = append(u.items, Err(value))
}

// AddTime adds a time.Time-typed Unifield to the list.
func (u *UnitfieldList) AddTime(value time.Time) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.items = append(u.items, Time(value))
}
