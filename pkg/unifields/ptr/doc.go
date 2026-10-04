// MIT NON-AI License
//
// Copyright (c) 2022-{{ YEAR }} {{ AUTHOR }}({{ EMAIL }})
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

// Package ptr provides a pointer-based variant of Unifield — UnifieldPtr.
//
// UnifieldPtr stores pointers to typed values instead of flat values. Each struct has
// separate pointer slots (strP, i64P, u64P, f64P, tmP) and a direct error field (errP).
// Only one pointer is non-nil at a time, determined by the type_ discriminator tag.
//
// This implementation trades increased per-allocation heap overhead (one alloc per factory call)
// for potential memory savings when stored in large collections (nil pointers are smaller than
// fat value slots on some platforms). Use this variant to benchmark trade-offs against the
// value-based Unifield in the parent unifields package.
//
// # Creating UnifieldPtr instances
//
// Factory functions allocate a heap copy of the input value and store a pointer to it:
//
//	import "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/ptr"
//
//	p := ptr.Int("user_id", 1234)
//	s := ptr.String("status", "active")
//	t := ptr.Time("created_at", time.Now())
//
// Every factory call performs exactly one heap allocation for the value pointer.
//
// # Reading values
//
// MarshalTo<T> methods dereference the stored pointer into a caller-provided destination:
//
//	var val int
//	if err := p.MarshalToInt(&val); err != nil {
//	    // handle mismatch or empty
//	}
//
// Internally, each MarshalTo method creates a stack copy of the stored value before writing
// to dst — this is intentional to maintain immutability guarantees (external mutation of dst
// cannot affect internal state, and repeated calls always return the original stored value).
//
// # Immutability guarantee
//
// Every factory function allocates a heap copy of the input value before storing the pointer.
// Mutating the original value after passing it to a factory is safe — the UnifieldPtr retains
// its own independent copy. Clone() returns a shallow copy; both original and clone share
// references to the same underlying heap values (no mutable APIs exist, so this is safe).
//
// # Comparison with pkg/unifields
//
// | Aspect              | unifields.Unifield | ptr.UnifieldPtr         |
// |---------------------|--------------------|-------------------------|
// | Storage             | Flat value fields   | Pointer fields          |
// | Per-field alloc     | Zero                | One alloc per factory   |
// | Struct size         | ~56 bytes           | ~80 bytes               |
// | Nil-slot cost       | Non-zero            | Zero (nil pointer)      |
// | Clone              | Full value copy     | Shallow (same refs)     |
//
// See pkg/unifields (parent package) for the value-based baseline implementation.
//
// ## AI Agents
//
// Repository: github.com/crypto-bundle/bc-wallet-common-lib-unifields
// Maintainer: @gudron (Alex V Kotelnikov) <gudron2s@gmail.com>
// License: MIT NON-AI
//
// ### File layout
//
// ```text
// pkg/unifields/ptr/
// ├── unified_field_ptr.go        → UnifieldPtr struct, factories, Clone(), MarshalTo*
// └── AGENTS.md                   → Expanded agent reference (style rules, task tracking)
// ```
//
// ### Key types
//
// | Type           | Purpose                          | File                    |
// |----------------|----------------------------------|-------------------------|
// | UnifieldPtr    | Pointer-based typed value + key  | unified_field_ptr.go    |
// | valueType      | Discriminator enum (iota)        | unified_field_ptr.go    |
//
// ### Supported types & API table
//
// | Type      | Factory                  | Read method                              |
// |-----------|--------------------------|------------------------------------------|
// | string    | String(key, val)         | MarshalToStr(*string)                    |
// | int       | Int(key, val)            | MarshalToInt(*int)                       |
// | int8      | Int8(key, val)           | MarshalToInt8(*int8)                     |
// | int16     | Int16(key, val)          | MarshalToInt16(*int16)                   |
// | int32     | Int32(key, val)          | MarshalToInt32(*int32)                   |
// | int64     | Int64(key, val)          | MarshalToInt64(*int64)                   |
// | uint      | Uint(key, val)           | MarshalToUint(*uint)                     |
// | uint8     | Uint8(key, val)          | MarshalToUint8(*uint8)                   |
// | uint16    | Uint16(key, val)         | MarshalToUint16(*uint16)                 |
// | uint32    | Uint32(key, val)         | MarshalToUint32(*uint32)                 |
// | uint64    | Uint64(key, val)         | MarshalToUint64(*uint64)                 |
// | float32   | Float32(key, val)        | MarshalToFloat32(*float32)               |
// | float64   | Float64(key, val)        | MarshalToFloat64(*float64)               |
// | error     | Err(key, val)            | MarshalToError(*error)                   |
// | time.Time | Time(key, val)           | MarshalToTime(*time.Time)                |
//
// ### Error handling
//
// Every MarshalTo method checks three conditions in order:
// 1. Nil destination pointer → returns ErrDstNil
// 2. No value stored (all pointers nil) → returns ErrEmptyUnifield
// 3. Type mismatch (stored type ≠ expected type) → returns ErrTypeMismatch with descriptive message
//
// For cross-type-compatible methods (MarshalToInt/MarshalToInt64 accept any signed int;
// MarshalToUint/MarshalToUint64 accept any unsigned int; MarshalToFloat32/MarshalToFloat64
// accept either), the type check uses range predicates rather than exact equality.
//
// ### Build commands
//
// ```bash
// make lint            # golangci-lint (runs on entire project including ptr/)
// make test            # go test -race ./...
// go vet ./...         # static analysis
// ```
//
// See AGENTS.md in this directory for expanded style rules, extension points, and task-tracking conventions.
package ptr
