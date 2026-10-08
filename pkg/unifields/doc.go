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

// Package unifields provides an immutable collection supporting both val.Unifield and ptr.UnifieldPtr.
//
// # Architecture Overview
//
// The unifields library offers two implementations for holding typed values, each optimised
// for different performance trade-offs:
//
//	Value-based (pkg/unifields/val): zero allocations per field, ~56-byte struct.
//	Use when allocation-free hot paths are critical.
//
//	Pointer-based (pkg/unifields/ptr): one heap alloc per factory call, ~80-byte struct.
//	Use when pointer identity or deferred value binding is needed.
//
// Both implementations implement the [Unifielder] interface, allowing them to be mixed
// freely inside a single [UnitfieldList] collection. For new code, prefer [UnitfieldList]
// over the legacy [Unifields] type.
//
// # Creating Typed Values
//
// Each sub-package provides factory functions that create typed value containers.
// Factory functions take a key string and a typed value:
//
//	import (
//	    "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields"
//	    "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/val"
//	    "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/ptr"
//	)
//
//	// Value-based (zero alloc)
//	f1 := val.String("name", "alice")
//	f2 := val.Int("status_code", 200)
//
//	// Pointer-based (one alloc per call)
//	f3 := ptr.Int("user_id", 1234)
//	f4 := ptr.Err("error", someErr)
//
// # Reading values
//
// Each factory function has corresponding typed methods prefixed with MarshalTo. These methods
// copy the stored value into a caller-provided pointer:
//
//	var val int
//	if err := f2.MarshalToInt(&val); err != nil {
//	    // handle type mismatch
//	}
//
// If the container holds a different type than expected, MarshalTo returns a descriptive error.
// Passing a nil destination pointer also returns an error.
//
// # Cloning
//
// Clone returns a shallow copy implementing [Unifielder]. Modifying the original after cloning
// does not affect the clone's stored values (for val types; ptr types share the same underlying
// pointers).
//
// ## Collections
//
// ### UnitfieldList (recommended)
//
// [UnitfieldList] is the recommended collection with list-manipulation capabilities:
//
//	list := unifields.NewUnitfieldList()
//	list.AddStr("name", "alice")
//	list.Merge(otherList)                     // merge another list (items cloned)
//	tail := list.GetAfter(0)                  // items from index N onward
//	list.RemoveAfter(0)                       // drop items after index 0
//	list.RemoveBefore(2)                      // drop items before index 2
//	list.Clear()                              // reset to empty
//
// ### UnifieldStack (LIFO)
//
// [UnifieldStack] is a Last-In-First-Out stack over Unifield values. The top element (last pushed)
// is returned first by Pop and Peek operations. Empty-stack behavior returns zero-value Unifield{} or nil —
// following the same convention as UnitfieldList methods that return slices.
//
// Memory management: val.Unifield instances use sync.Pool internally for efficient pre-allocation.
// Push/Pop/Flush/Clear operations automatically manage pool membership — callers do NOT need
// to manually manage pool lifecycle. Only val.Unifield participates in pooling; ptr.UnifieldPtr
// instances become regular GC garbage when no longer referenced.
//
//	stack := unifields.NewUnifieldStack()
//	stack.Push(unifields.String("key", "value"))
//	top := stack.Peek()                       // peek without removing
//	popped := stack.Pop()                     // remove, clone, return original to pool
//	all := stack.PopN(3)                      // pop up to 3 elements as clones
//	flushed := stack.Flush()                  // pop ALL remaining as clones, empty stack
//	stack.Clear()                             // empty stack, return items to pool
//
// Pop() returns a **clone** of the top element to preserve immutability — mutations to the
// returned value do not affect any item stored in the stack or in the pool. Original items
// are returned to the internal pool after being cloned.
//
// ### Unifields (deprecated)
//
// The legacy [Unifields] type remains for backward compatibility:
//
//	cols := unifields.NewUnifields()
//	cols.Add(val.String("name", "alice"))       // val.Unifield via Unifielder
//	cols.Add(ptr.Int("count", 42))              // ptr.UnifieldPtr via Unifielder
//
// Available methods for both collections:
//   - NewUnitfieldList() / NewUnifields() — creates an empty collection
//   - Add(fld Unifielder) — adds a cloned Unifielder (supports both val and ptr)
//   - AddAll(flds []Unifielder) — bulk adds cloned Unifielders
//   - Len() — returns item count
//   - Items() — returns read-only copy of stored items
//   - AddStr, AddInt, AddInt8..AddInt64, AddUint..AddUint64, AddFloat32, AddFloat64, AddErr, AddTime — typed adders
//
// Additional methods on UnitfieldList (thread-safe via sync.RWMutex):
//   - Merge(source *UnitfieldList) — appends all items from source list (each element cloned); DO NOT pass same list as receiver and source (deadlock due to non-reentrant mutex)
//   - GetAfter(index uint) []Unifielder — returns copy of items from index N onward
//   - RemoveAfter(index uint) — keeps element at index, drops after; removed items returned to pool
//   - RemoveBefore(index uint) — keeps element at index, drops before; removed items returned to pool
//   - Clear() — resets list to empty; all items returned to pool
//
// Memory management: val.Unifield instances in collections use sync.Pool internally for efficient
// pre-allocation. Removal operations (RemoveAfter, RemoveBefore, Clear) automatically return recycled
// val.Unifield instances to the pool. ptr.UnifieldPtr instances become regular GC garbage when unreferenced.
//
// # Choosing Between val and ptr
//
// | Criteria         | val.Unifield              | ptr.UnifieldPtr          |
// |------------------|---------------------------|--------------------------|
// | Allocs per field | 0                         | 1 (heap copy)            |
// | Struct size      | ~56 bytes                 | ~80 bytes + heap pointers|
// | Best for         | Zero-cost hot paths       | Pointer identity needed  |
//
// See detailed benchmark results in `.agents/reports/benchmarks_val_unifield.md`.
//
// # Supported types
//
// Both implementations support the same 15 basic types: string, integers (int, int8..int64,
// uint, uint8..uint64), floats (float32, float64), error, and time.Time.
//
// ## AI Agents
//
// Repository: github.com/crypto-bundle/bc-wallet-common-lib-unifields
// Maintainer: @gudron (Alex V Kotelnikov) <gudron2s@gmail.com>
// License: MIT NON-AI
//
// ### Memory pool (`sync.Pool`)
//
// val.Unifield uses a single package-level sync.Pool internally for efficient pre-allocation.
// Each factory function obtains an Unifield from the pool, resets fields, sets values, and returns
// a value-copy to the caller. On collection removal operations (RemoveAfter, RemoveBefore, Clear),
// val.Unifield instances are returned to the pool via returnValToPool().
//
// | Operation       | Pool usage                                    | Clone? | Return to pool? |
// |-----------------|-----------------------------------------------|--------|-----------------|
// | Factory         | Get → resetToZero() → set fields → *Unifield   | No     | N/A             |
// | Clone()         | Plain struct copy (value semantics)            | Yes    | No              |
// | Add/Push        | Caller's item cloned into collection           | Yes    | On removal      |
// | Getters         | Already return copies or interface references  | N/A    | No              |
// | Remove/Clear    | N/A                                           | N/A    | Yes             |
// | Pop/PopN        | Items removed, originals pooled after clone    | Yes    | Yes             |
// | Flush           | All items cloned, originals pooled, stack empty| Yes    | Yes             |
//
// Only val.Unifield participates in pooling; ptr.UnifieldPtr is silently skipped by returnValToPool()
// and becomes regular GC garbage when unreferenced. Users must NOT manually manage pool lifecycle.
//
// ### File layout
//
// ```text
// pkg/unifields/
// ├── doc.go                        → Package-level godoc + AI reference
// ├── unified_field.go              → Type alias (Unifield = val.Unifield) + 15 factory wrappers
// ├── unitfield_list.go             → Type alias only (val.UnitfieldList)
// ├── unified_fields.go             → Deprecated Unifields collection + Unifielder interface alias
// ├── unified_fields_test.go        → Collection tests (polymorphic val/ptr)
// └── unifielder/
//
//	└── common.go               → Unifielder interface definition
//
// └─️ val/
//
//	├── unified_field.go          → Unifield type, factories, Clone(), MarshalTo*
//	├── unified_field_test.go     → Unifield unit tests
//	└─️ unified_field_benchmark_test.go → Benchmarks (factory, clone, marshal, collection)
//	├── unified_field_list.go     → UnitfieldList type + list manipulation methods
//	├── unified_field_list_test.go         → UnitfieldList unit tests
//	└─️ unified_field_list_benchmark_test.go → UnitfieldList benchmarks
//
// └─️ ptr/
//
//	├── unified_field_ptr.go      → UnifieldPtr struct, factories, Clone(), MarshalTo*
//	├── unified_field_ptr_test.go → UnifieldPtr unit tests
//	├─️ unified_field_ptr_benchmark_test.go → Ptr benchmarks
//
// ```
//
// ### Key types and interfaces
//
// | Type/Interface | Purpose                             | Location               |
// |----------------|-------------------------------------|------------------------|
// | Unifield       | Single typed value + key (value)    | pkg/unifields/val/     |
// | UnifieldPtr    | Single typed value + key (pointer)  | pkg/unifields/ptr/     |
// | Unifielder     | Interface for polymorphic storage   | pkg/unifields/unifielder/ |
// | UnitfieldList  | Mutable collection with list ops    | pkg/unifields/         |
// | Unifields      | Immutable collection (deprecated)   | pkg/unifields/         |
// | valueType      | Discriminator enum                  | internal to val/ & ptr/|
//
// ### Supported types & API table
//
// | Type | Factory (val/ptr) | Read method | Collection adder |
// |------|--------------------|-------------|------------------|
// | string | String(key, val) | MarshalToStr(*string) | AddStr |
// | int | Int(key, val) | MarshalToInt(*int) | AddInt |
// | int8 | Int8(key, val) | MarshalToInt8(*int8) | AddInt8 |
// | int16 | Int16(key, val) | MarshalToInt16(*int16) | AddInt16 |
// | int32 | Int32(key, val) | MarshalToInt32(*int32) | AddInt32 |
// | int64 | Int64(key, val) | MarshalToInt64(*int64) | AddInt64 |
// | uint | Uint(key, val) | MarshalToUint(*uint) | AddUint |
// | uint8 | Uint8(key, val) | MarshalToUint8(*uint8) | AddUint8 |
// | uint16 | Uint16(key, val) | MarshalToUint16(*uint16) | AddUint16 |
// | uint32 | Uint32(key, val) | MarshalToUint32(*uint32) | AddUint32 |
// | uint64 | Uint64(key, val) | MarshalToUint64(*uint64) | AddUint64 |
// | float32 | Float32(key, val) | MarshalToFloat32(*float32) | AddFloat32 |
// | float64 | Float64(key, val) | MarshalToFloat64(*float64) | AddFloat64 |
// | error | Err(key, val) | MarshalToError(*error) | AddErr |
// | time.Time | Time(key, val) | MarshalToTime(*time.Time) | AddTime |
//
// ### Extension points (adding a new basic type)
//
// 1. Add iota constant to `valueType` enum in `pkg/unifields/val/unified_field.go`
// 2. Add backing field to `Unifield` struct
// 3. Create factory function in val/ (`func NewType(key string, val Type) Unifield`)
// 4. Create `MarshalToType(dst *Type) error` receiver method in val/
// 5. Repeat steps 1-4 for ptr/ (`UnifieldPtr` variant)
// 6. Add Add<Type>() method to `pkg/unifields/unitfield_list.go` (delegates to val factory on UnitfieldList)
// 7. Add test cases to both `pkg/unifields/val/unified_field_test.go` and `pkg/unifields/ptr/unified_field_ptr_test.go`
// 8. Update `typeName()` switch in both val/ and ptr/ with new case label
// 9. Update helper predicates (`isSignedInt` / `isUnsignedInt`) if type belongs to range
//
// ### Build commands
//
// ```bash
// make lint            # golangci-lint
// make test            # go test -race ./...
// go vet ./...         # static analysis
// go test -bench=. -benchmem ./pkg/unifields/... # benchmarks
// ```
//
// See AGENTS.md at repo root for expanded style rules, CSG references, and task-tracking conventions.
package unifields
