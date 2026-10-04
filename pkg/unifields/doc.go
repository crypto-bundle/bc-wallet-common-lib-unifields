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

// Package unifields provides a typed value container similar to zapcore.Field.
//
// An Unifield holds exactly one typed value paired with a string key identifier.
// It supports the following basic types: string, integers (int, int8..int64, uint, uint8..uint64),
// floating-point numbers (float32, float64), error, and time.Time. Internally, each Unifield uses
// flat storage — a single struct with separate fields for each type and a valueType discriminator tag.
// This design avoids heap allocation per field and enables cheap copy semantics via Clone().
//
// # Creating Unifields
//
// Factory functions create Unifield instances. Each function takes a key string and a typed value:
//
//	import "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields"
//
//	// Integer values
//	f1 := unifields.Int("status_code", 200)
//	f2 := unifields.Int64("count", 42)
//	f3 := unifields.Uint64("id", 18446744073709551615)
//
//	// String and error
//	f4 := unifields.String("message", "ok")
//	f5 := unifields.Err("error", someError)
//
//	// Time
//	f6 := unifields.Time("timestamp", time.Now())
//
// # Reading values
//
// Each factory function has corresponding typed methods prefixed with MarshalTo. These methods
// copy the stored value into a caller-provided pointer:
//
//	var val int
//	if err := f1.MarshalToInt(&val); err != nil {
//	    // handle type mismatch
//	}
//
// If the Unifield holds a different type than expected, MarshalTo returns a descriptive error.
// Passing a nil destination pointer also returns an error.
//
// # Cloning
//
// Clone returns a shallow copy of the Unifield. Since all stored fields are value-types
// (except error which is an interface), modifying the original after cloning does not affect
// the clone's stored values.
//
// # Unifields collection
//
// The Unifolds type provides an immutable collection wrapper around []*Unifold. All Add operations
// store clones of input values, ensuring external mutation cannot affect the collection's contents.
//
//	import "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields"
//
//	cols := unifields.NewUnifolds()
//	cols.AddInt("user_id", 1234)
//	cols.AddStr("status", "active")
//	cols.AddTime("created_at", time.Now())
//
// Available methods:
//   - NewUnifolds() — creates an empty collection
//   - Add(fld Unifold) — adds a cloned Unifold
//   - AddAll(flds []Unifold) — bulk adds cloned Unifolds
//   - AddStr, AddInt, AddInt8..AddInt64, AddUint..AddUint64, AddFloat32, AddFloat64, AddErr, AddTime — typed adders
//
// # ptr package
//
// The `ptr` sub-package (`github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/ptr`)
// provides a pointer-based variant of Unifield called `UnifieldPtr`. It stores pointers to typed values
// instead of flat values. Each struct has separate pointer slots (`strP`, `i64P`, `u64P`, `f64P`, `tmP`)
// and a direct error field (`errP`). Only one pointer is non-nil at a time, determined by the `type_`
// discriminator tag.
//
// This variant trades increased per-allocation heap overhead (one alloc per factory call) for potential
// memory savings when stored in large collections (nil pointers are smaller than fat value slots on some
// platforms). Use it to benchmark trade-offs against the value-based `Unifield`:
//
//	import "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/ptr"
//
//	p := ptr.Int("user_id", 1234)
//	var val int
//	p.MarshalToInt(&val) // val == 1234
//
// Factory functions allocate a heap copy of the input value so external mutation cannot affect the stored
// value. Clone() returns a shallow copy — both original and clone share references to the same underlying
// heap values (no mutable APIs exist, so this is safe). See `pkg/unifields/ptr/doc.go` for full details.
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
// pkg/unifields/
// ├── doc.go                      → Package-level godoc + AI reference
// ├── unified_field.go            → Unifold struct, factories, Clone(), MarshalTo*
// ├── unified_field_test.go       → Unifold tests (factories, MarshalTo, Clone, errors)
// ├── unified_fields.go           → Unifolds collection (New, Add, AddAll, typed adders)
// └── unified_fields_test.go      → Collection tests (immutability, clone separation)
// ```
//
// ### Key types
//
// | Type          | Purpose                          | File                |
// |---------------|----------------------------------|---------------------|
// | Unifold       | Single typed value + key         | unified_field.go    |
// | Unifolds      | Immutable collection of Unifold  | unified_fields.go   |
// | valueType     | Discriminator enum (iota)        | unified_field.go    |
//
// ### Supported types & API table
//
// | Type | Factory | Read method | Collection adder |
// |------|---------|-------------|------------------|
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
// 1. Add iota constant to `valueType` enum in `unified_field.go`
// 2. Add backing field to `Unifold` struct (e.g. `i32 int32`)
// 3. Create factory function (`func Int32(key string, val int32) Unifold`)
// 4. Create `MarshalToInt32(dst *int32) error` receiver method
// 5. Create `AddInt32(key string, val int32)` collection adder
// 6. Add test to `unified_field_test.go` + `unified_fields_test.go`
// 7. Update `typeName()` switch with new case label
// 8. Update helper predicates (`isSignedInt` / `isUnsignedInt`) if type belongs to range
//
// ### Build commands
//
// ```bash
// make lint            # golangci-lint
// make test            # go test -race ./...
// go vet ./...         # static analysis
// ```
//
// See AGENTS.md at repo root for expanded style rules, CSG references, and task-tracking conventions.
package unifields
