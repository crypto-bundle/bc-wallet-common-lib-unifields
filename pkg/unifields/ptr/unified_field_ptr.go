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
// See pkg/unifields (parent package) for the value-based baseline implementation.
package ptr

import (
	"errors"
	"fmt"
	"time"
)

// valueType identifies the stored datatype inside an UnifieldPtr instance.
type valueType uint8

const (
	valueEmpty   valueType = iota // no value stored
	valueString                   // strP holds pointer to string value
	valueInt                      // i64P holds pointer to int64 (int value)
	valueInt8                     // i64P holds pointer to int64 (int8 value)
	valueInt16                    // i64P holds pointer to int64 (int16 value)
	valueInt32                    // i64P holds pointer to int64 (int32 value)
	valueInt64                    // i64P holds pointer to int64 (int64 value)
	valueUint                     // u64P holds pointer to uint64 (uint value)
	valueUint8                    // u64P holds pointer to uint64 (uint8 value)
	valueUint16                   // u64P holds pointer to uint64 (uint16 value)
	valueUint32                   // u64P holds pointer to uint64 (uint32 value)
	valueUint64                   // u64P holds pointer to uint64 (uint64 value)
	valueFloat32                  // f64P holds pointer to float64 (float32 value)
	valueFloat64                  // f64P holds pointer to float64 (float64 value)
	valueError                    // errP holds error value directly
	valueTime                     // tmP holds pointer to time.Time value
)

// typeName returns a human-readable name for a valueType.
func typeName(t valueType) string {
	switch t {
	case valueEmpty:
		return "empty"
	case valueString:
		return "str"
	case valueInt:
		return "int"
	case valueInt8:
		return "int8"
	case valueInt16:
		return "int16"
	case valueInt32:
		return "int32"
	case valueInt64:
		return "int64"
	case valueUint:
		return "uint"
	case valueUint8:
		return "uint8"
	case valueUint16:
		return "uint16"
	case valueUint32:
		return "uint32"
	case valueUint64:
		return "uint64"
	case valueFloat32:
		return "float32"
	case valueFloat64:
		return "float64"
	case valueError:
		return "error"
	case valueTime:
		return "time.Time"
	default:
		return "unknown"
	}
}

// isSignedInt reports whether the valueType represents a signed integer type (not unsigned).
func isSignedInt(t valueType) bool {
	return t >= valueInt && t < valueUint
}

// isUnsignedInt reports whether the valueType represents an unsigned integer type (not signed).
func isUnsignedInt(t valueType) bool {
	return t >= valueUint && t <= valueUint64
}

// Static error values for common failure cases.
var (
	ErrDstNil        = errors.New("dst is nil")
	ErrTypeMismatch  = errors.New("type mismatch")
	ErrEmptyUnifield = errors.New("unifield is empty")
)

// UnifieldPtr holds exactly one typed value via pointer, paired with a key identifier.
//
// All internal pointer fields start nil; only one is non-nil after construction,
// as indicated by the type_ discriminator tag. The errP field stores error directly
// (not as a pointer) because error is already an interface type — storing *error would
// require double-indirection (*interface{}) which adds no value.
//
// Immutability guarantee: every factory function copies the input value onto the heap
// before storing the pointer. Callers can safely mutate their original values after
// passing them to a factory — the UnifieldPtr retains its own independent copy.
//
// Field layout: key + error interface consume 32 bytes minimum; five pointers follow
// in natural order. Due to alignment constraints with string header (16B) and error
// interface (16B), the struct totals 80 bytes — optimized as much as possible given
// Go's memory model requirements.
type UnifieldPtr struct { //nolint:govet // alignment constrained by string+error+5*pointer types
	key   string     // identifier key for the field
	errP  error      // error value stored directly (error is already an interface)
	i64P  *int64     // pointer to signed int value (nil if not set)
	u64P  *uint64    // pointer to unsigned int value (nil if not set)
	f64P  *float64   // pointer to float value (nil if not set)
	tmP   *time.Time // pointer to time.Time value (nil if not set)
	strP  *string    // pointer to string value (nil if not set)
	type_ valueType  // discriminator tag — exactly one pointer is non-nil (or errP is set)
}

// Clone returns a shallow copy of the UnifieldPtr. Since all pointer fields are
// copied by value (the pointers themselves), both the original and clone share
// references to the same underlying heap values. Modifications to those heap
// values through one clone will be visible through the other — however, the
// UnifieldPtr struct itself never exposes mutable APIs, so this is safe in practice.
func (u UnifieldPtr) Clone() UnifieldPtr {
	return u
}

// MarshalToStr copies the stored string into dst.
func (u UnifieldPtr) MarshalToStr(dst *string) error {
	if dst == nil {
		return ErrDstNil
	}
	if u.strP == nil {
		return ErrEmptyUnifield
	}
	if u.type_ != valueString {
		return fmt.Errorf("%w: got %s, want str", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = *u.strP
	return nil
}

// MarshalToInt copies the stored int into dst.
func (u UnifieldPtr) MarshalToInt(dst *int) error {
	if dst == nil {
		return ErrDstNil
	}
	if u.i64P == nil {
		return ErrEmptyUnifield
	}
	if !isSignedInt(u.type_) {
		return fmt.Errorf("%w: got %s, want int", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = int(*u.i64P)
	return nil
}

// MarshalToInt8 copies the stored int8 into dst.
func (u UnifieldPtr) MarshalToInt8(dst *int8) error {
	if dst == nil {
		return ErrDstNil
	}
	if u.i64P == nil {
		return ErrEmptyUnifield
	}
	if u.type_ != valueInt8 {
		return fmt.Errorf("%w: got %s, want int8", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = int8(*u.i64P)
	return nil
}

// MarshalToInt16 copies the stored int16 into dst.
func (u UnifieldPtr) MarshalToInt16(dst *int16) error {
	if dst == nil {
		return ErrDstNil
	}
	if u.i64P == nil {
		return ErrEmptyUnifield
	}
	if u.type_ != valueInt16 {
		return fmt.Errorf("%w: got %s, want int16", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = int16(*u.i64P)
	return nil
}

// MarshalToInt32 copies the stored int32 into dst.
func (u UnifieldPtr) MarshalToInt32(dst *int32) error {
	if dst == nil {
		return ErrDstNil
	}
	if u.i64P == nil {
		return ErrEmptyUnifield
	}
	if u.type_ != valueInt32 {
		return fmt.Errorf("%w: got %s, want int32", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = int32(*u.i64P)
	return nil
}

// MarshalToInt64 copies the stored int64 into dst.
func (u UnifieldPtr) MarshalToInt64(dst *int64) error {
	if dst == nil {
		return ErrDstNil
	}
	if u.i64P == nil {
		return ErrEmptyUnifield
	}
	if !isSignedInt(u.type_) {
		return fmt.Errorf("%w: got %s, want int64", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = *u.i64P
	return nil
}

// MarshalToUint copies the stored uint into dst.
func (u UnifieldPtr) MarshalToUint(dst *uint) error {
	if dst == nil {
		return ErrDstNil
	}
	if u.u64P == nil {
		return ErrEmptyUnifield
	}
	if u.type_ != valueUint {
		return fmt.Errorf("%w: got %s, want uint", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = uint(*u.u64P)
	return nil
}

// MarshalToUint8 copies the stored uint8 into dst.
func (u UnifieldPtr) MarshalToUint8(dst *uint8) error {
	if dst == nil {
		return ErrDstNil
	}
	if u.u64P == nil {
		return ErrEmptyUnifield
	}
	if u.type_ != valueUint8 {
		return fmt.Errorf("%w: got %s, want uint8", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = uint8(*u.u64P)
	return nil
}

// MarshalToUint16 copies the stored uint16 into dst.
func (u UnifieldPtr) MarshalToUint16(dst *uint16) error {
	if dst == nil {
		return ErrDstNil
	}
	if u.u64P == nil {
		return ErrEmptyUnifield
	}
	if u.type_ != valueUint16 {
		return fmt.Errorf("%w: got %s, want uint16", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = uint16(*u.u64P)
	return nil
}

// MarshalToUint32 copies the stored uint32 into dst.
func (u UnifieldPtr) MarshalToUint32(dst *uint32) error {
	if dst == nil {
		return ErrDstNil
	}
	if u.u64P == nil {
		return ErrEmptyUnifield
	}
	if u.type_ != valueUint32 {
		return fmt.Errorf("%w: got %s, want uint32", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = uint32(*u.u64P)
	return nil
}

// MarshalToUint64 copies the stored uint64 into dst.
func (u UnifieldPtr) MarshalToUint64(dst *uint64) error {
	if dst == nil {
		return ErrDstNil
	}
	if u.u64P == nil {
		return ErrEmptyUnifield
	}
	if !isUnsignedInt(u.type_) {
		return fmt.Errorf("%w: got %s, want uint64", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = *u.u64P
	return nil
}

// MarshalToFloat32 copies the stored float32 into dst.
func (u UnifieldPtr) MarshalToFloat32(dst *float32) error {
	if dst == nil {
		return ErrDstNil
	}
	if u.f64P == nil {
		return ErrEmptyUnifield
	}
	if u.type_ != valueFloat32 && u.type_ != valueFloat64 {
		return fmt.Errorf("%w: got %s, want float32", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = float32(*u.f64P)
	return nil
}

// MarshalToFloat64 copies the stored float64 into dst.
func (u UnifieldPtr) MarshalToFloat64(dst *float64) error {
	if dst == nil {
		return ErrDstNil
	}
	if u.f64P == nil {
		return ErrEmptyUnifield
	}
	if u.type_ != valueFloat32 && u.type_ != valueFloat64 {
		return fmt.Errorf("%w: got %s, want float64", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = *u.f64P
	return nil
}

// MarshalToError copies the stored error into dst.
func (u UnifieldPtr) MarshalToError(dst *error) error {
	if dst == nil {
		return ErrDstNil
	}
	if u.errP == nil {
		return ErrEmptyUnifield
	}
	if u.type_ != valueError {
		return fmt.Errorf("%w: got %s, want error", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = u.errP
	return nil
}

// MarshalToTime copies the stored time.Time into dst.
func (u UnifieldPtr) MarshalToTime(dst *time.Time) error {
	if dst == nil {
		return ErrDstNil
	}
	if u.tmP == nil {
		return ErrEmptyUnifield
	}
	if u.type_ != valueTime {
		return fmt.Errorf("%w: got %s, want time.Time", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = *u.tmP
	return nil
}

// --- Factory functions ---

// String creates a new UnifieldPtr with key and string value.
// Allocates a heap copy of val — mutating the original string after this call is safe.
func String(key string, val string) UnifieldPtr {
	v := val
	return UnifieldPtr{key: key, strP: &v, type_: valueString}
}

// Int creates a new UnifieldPtr with key and int value.
// Allocates a heap copy of val (as int64) — mutating the original int after this call is safe.
func Int(key string, val int) UnifieldPtr {
	v := int64(val)
	return UnifieldPtr{key: key, i64P: &v, type_: valueInt}
}

// Int8 creates a new UnifieldPtr with key and int8 value.
// Allocates a heap copy of val (as int64) — mutating the original int8 after this call is safe.
func Int8(key string, val int8) UnifieldPtr {
	v := int64(val)
	return UnifieldPtr{key: key, i64P: &v, type_: valueInt8}
}

// Int16 creates a new UnifieldPtr with key and int16 value.
// Allocates a heap copy of val (as int64) — mutating the original int16 after this call is safe.
func Int16(key string, val int16) UnifieldPtr {
	v := int64(val)
	return UnifieldPtr{key: key, i64P: &v, type_: valueInt16}
}

// Int32 creates a new UnifieldPtr with key and int32 value.
// Allocates a heap copy of val (as int64) — mutating the original int32 after this call is safe.
func Int32(key string, val int32) UnifieldPtr {
	v := int64(val)
	return UnifieldPtr{key: key, i64P: &v, type_: valueInt32}
}

// Int64 creates a new UnifieldPtr with key and int64 value.
// Allocates a heap copy of val — mutating the original int64 after this call is safe.
func Int64(key string, val int64) UnifieldPtr {
	v := val
	return UnifieldPtr{key: key, i64P: &v, type_: valueInt64}
}

// Uint creates a new UnifieldPtr with key and uint value.
// Allocates a heap copy of val (as uint64) — mutating the original uint after this call is safe.
func Uint(key string, val uint) UnifieldPtr {
	v := uint64(val)
	return UnifieldPtr{key: key, u64P: &v, type_: valueUint}
}

// Uint8 creates a new UnifieldPtr with key and uint8 value.
// Allocates a heap copy of val (as uint64) — mutating the original uint8 after this call is safe.
func Uint8(key string, val uint8) UnifieldPtr {
	v := uint64(val)
	return UnifieldPtr{key: key, u64P: &v, type_: valueUint8}
}

// Uint16 creates a new UnifieldPtr with key and uint16 value.
// Allocates a heap copy of val (as uint64) — mutating the original uint16 after this call is safe.
func Uint16(key string, val uint16) UnifieldPtr {
	v := uint64(val)
	return UnifieldPtr{key: key, u64P: &v, type_: valueUint16}
}

// Uint32 creates a new UnifieldPtr with key and uint32 value.
// Allocates a heap copy of val (as uint64) — mutating the original uint32 after this call is safe.
func Uint32(key string, val uint32) UnifieldPtr {
	v := uint64(val)
	return UnifieldPtr{key: key, u64P: &v, type_: valueUint32}
}

// Uint64 creates a new UnifieldPtr with key and uint64 value.
// Allocates a heap copy of val — mutating the original uint64 after this call is safe.
func Uint64(key string, val uint64) UnifieldPtr {
	v := val
	return UnifieldPtr{key: key, u64P: &v, type_: valueUint64}
}

// Float32 creates a new UnifieldPtr with key and float32 value.
// Allocates a heap copy of val (as float64) — mutating the original float32 after this call is safe.
func Float32(key string, val float32) UnifieldPtr {
	v := float64(val)
	return UnifieldPtr{key: key, f64P: &v, type_: valueFloat32}
}

// Float64 creates a new UnifieldPtr with key and float64 value.
// Allocates a heap copy of val — mutating the original float64 after this call is safe.
func Float64(key string, val float64) UnifieldPtr {
	v := val
	return UnifieldPtr{key: key, f64P: &v, type_: valueFloat64}
}

// Err creates a new UnifieldPtr with key and error value.
// Stores the error interface directly — since error is already an interface, no pointer indirection is needed.
// Passing nil as val produces an empty UnifieldPtr.
func Err(key string, val error) UnifieldPtr {
	if val == nil {
		return UnifieldPtr{key: key, type_: valueEmpty}
	}
	return UnifieldPtr{key: key, errP: val, type_: valueError}
}

// Time creates a new UnifieldPtr with key and time.Time value.
// Allocates a heap copy of val — mutating the original time.Time after this call is safe.
func Time(key string, val time.Time) UnifieldPtr {
	v := val
	return UnifieldPtr{key: key, tmP: &v, type_: valueTime}
}
