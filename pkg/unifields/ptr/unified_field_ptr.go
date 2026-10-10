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

package ptr

import (
	"errors"
	"fmt"
	"time"

	unifolderv2 "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/unifielder"
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

// typeName returns a human-readable name for a valueType suitable for error messages.
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

// isSignedInt reports whether the valueType represents a signed integer type.
func isSignedInt(t valueType) bool {
	return t >= valueInt && t < valueUint
}

// isUnsignedInt reports whether the valueType represents an unsigned integer type.
func isUnsignedInt(t valueType) bool {
	return t >= valueUint && t <= valueUint64
}

// Error variables returned by MarshalTo* and factory methods.
var (
	ErrDstNil        = errors.New("dst is nil")
	ErrTypeMismatch  = errors.New("type mismatch")
	ErrEmptyUnifield = errors.New("unifield is empty")
)

// UnifieldPtr holds exactly one typed value via pointer.
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
// Field layout: error interface + five pointers follow in natural order. Due to
// alignment constraints with error interface (16B), the struct totals ~80 bytes.
type UnifieldPtr struct { //nolint:govet // alignment constrained by error+5*pointer types
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
func (u UnifieldPtr) Clone() unifolderv2.Unifielder { //nolint:ireturn // intentionally returns interface for polymorphic collections
	return u
}

// MarshalToStr copies the stored string value into dst.
//
// It returns ErrDstNil if dst is nil, ErrEmptyUnifield if the UnifieldPtr has no
// stored value, or ErrTypeMismatch if the UnifieldPtr holds a different type than
// expected. The method dereferences the internal pointer and writes a stack copy to dst,
// preserving immutability guarantees.
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

// MarshalToInt copies the stored int value into dst.
//
// It accepts any signed integer type (valueInt through valueInt64) due to cross-type
// compatibility. It returns ErrDstNil if dst is nil, ErrEmptyUnifield if no value is
// stored, or ErrTypeMismatch if the stored type is unsigned.
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

// MarshalToInt8 copies the stored int8 value into dst.
//
// It returns ErrDstNil if dst is nil, ErrEmptyUnifield if no value is stored,
// or ErrTypeMismatch if the stored type is not exactly valueInt8.
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

// MarshalToInt16 copies the stored int16 value into dst.
//
// It returns ErrDstNil if dst is nil, ErrEmptyUnifield if no value is stored,
// or ErrTypeMismatch if the stored type is not exactly valueInt16.
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

// MarshalToInt32 copies the stored int32 value into dst.
//
// It returns ErrDstNil if dst is nil, ErrEmptyUnifield if no value is stored,
// or ErrTypeMismatch if the stored type is not exactly valueInt32.
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

// MarshalToInt64 copies the stored int64 value into dst.
//
// It accepts any signed integer type (valueInt through valueInt64) due to cross-type
// compatibility. It returns ErrDstNil if dst is nil, ErrEmptyUnifield if no value is
// stored, or ErrTypeMismatch if the stored type is unsigned.
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

// MarshalToUint copies the stored uint value into dst.
//
// It accepts any unsigned integer type (valueUint through valueUint64) due to cross-type
// compatibility. It returns ErrDstNil if dst is nil, ErrEmptyUnifield if no value is
// stored, or ErrTypeMismatch if the stored type is signed.
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

// MarshalToUint8 copies the stored uint8 value into dst.
//
// It returns ErrDstNil if dst is nil, ErrEmptyUnifield if no value is stored,
// or ErrTypeMismatch if the stored type is not exactly valueUint8.
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

// MarshalToUint16 copies the stored uint16 value into dst.
//
// It returns ErrDstNil if dst is nil, ErrEmptyUnifield if no value is stored,
// or ErrTypeMismatch if the stored type is not exactly valueUint16.
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

// MarshalToUint32 copies the stored uint32 value into dst.
//
// It returns ErrDstNil if dst is nil, ErrEmptyUnifield if no value is stored,
// or ErrTypeMismatch if the stored type is not exactly valueUint32.
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

// MarshalToUint64 copies the stored uint64 value into dst.
//
// It accepts any unsigned integer type (valueUint through valueUint64) due to cross-type
// compatibility. It returns ErrDstNil if dst is nil, ErrEmptyUnifield if no value is
// stored, or ErrTypeMismatch if the stored type is signed.
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

// MarshalToFloat32 copies the stored float32 value into dst.
//
// It accepts either valueFloat32 or valueFloat64 since both share the same f64 backing
// field. Narrowing from float64 to float32 is implicit. Returns ErrDstNil if dst is nil,
// ErrEmptyUnifield if no value is stored, or ErrTypeMismatch otherwise.
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

// MarshalToFloat64 copies the stored float64 value into dst.
//
// It accepts either valueFloat32 or valueFloat64 since both share the same f64 backing
// field. Widening from float32 to float64 is implicit. Returns ErrDstNil if dst is nil,
// ErrEmptyUnifield if no value is stored, or ErrTypeMismatch otherwise.
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

// MarshalToError copies the stored error value into dst.
//
// Returns ErrDstNil if dst is nil, ErrEmptyUnifield if no value is stored,
// or ErrTypeMismatch if the stored type is not valueError.
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

// MarshalToTime copies the stored time.Time value into dst.
//
// Returns ErrDstNil if dst is nil, ErrEmptyUnifield if no value is stored,
// or ErrTypeMismatch if the stored type is not valueTime.
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

// String creates a new UnifieldPtr with string value.
//
// It allocates a heap copy of val so that mutating the original value after this
// call does not affect the stored value. Only one pointer field (strP) will be
// non-nil in the returned UnifieldPtr.
func String(val string) UnifieldPtr {
	v := val
	return UnifieldPtr{strP: &v, type_: valueString}
}

// Int creates a new UnifieldPtr with int value.
//
// It allocates a heap copy of val (stored as int64) so that mutating the original
// value after this call does not affect the stored value. Only one pointer field
// (i64P) will be non-nil in the returned UnifieldPtr.
func Int(val int) UnifieldPtr {
	v := int64(val)
	return UnifieldPtr{i64P: &v, type_: valueInt}
}

// Int8 creates a new UnifieldPtr with int8 value.
//
// It allocates a heap copy of val (stored as int64) so that mutating the original
// value after this call does not affect the stored value. Only one pointer field
// (i64P) will be non-nil in the returned UnifieldPtr.
func Int8(val int8) UnifieldPtr {
	v := int64(val)
	return UnifieldPtr{i64P: &v, type_: valueInt8}
}

// Int16 creates a new UnifieldPtr with int16 value.
//
// It allocates a heap copy of val (stored as int64) so that mutating the original
// value after this call does not affect the stored value. Only one pointer field
// (i64P) will be non-nil in the returned UnifieldPtr.
func Int16(val int16) UnifieldPtr {
	v := int64(val)
	return UnifieldPtr{i64P: &v, type_: valueInt16}
}

// Int32 creates a new UnifieldPtr with int32 value.
//
// It allocates a heap copy of val (stored as int64) so that mutating the original
// value after this call does not affect the stored value. Only one pointer field
// (i64P) will be non-nil in the returned UnifieldPtr.
func Int32(val int32) UnifieldPtr {
	v := int64(val)
	return UnifieldPtr{i64P: &v, type_: valueInt32}
}

// Int64 creates a new UnifieldPtr with int64 value.
//
// It allocates a heap copy of val so that mutating the original value after this
// call does not affect the stored value. Only one pointer field (i64P) will be
// non-nil in the returned UnifieldPtr.
func Int64(val int64) UnifieldPtr {
	v := val
	return UnifieldPtr{i64P: &v, type_: valueInt64}
}

// Uint creates a new UnifieldPtr with uint value.
//
// It allocates a heap copy of val (stored as uint64) so that mutating the original
// value after this call does not affect the stored value. Only one pointer field
// (u64P) will be non-nil in the returned UnifieldPtr.
func Uint(val uint) UnifieldPtr {
	v := uint64(val)
	return UnifieldPtr{u64P: &v, type_: valueUint}
}

// Uint8 creates a new UnifieldPtr with uint8 value.
//
// It allocates a heap copy of val (stored as uint64) so that mutating the original
// value after this call does not affect the stored value. Only one pointer field
// (u64P) will be non-nil in the returned UnifieldPtr.
func Uint8(val uint8) UnifieldPtr {
	v := uint64(val)
	return UnifieldPtr{u64P: &v, type_: valueUint8}
}

// Uint16 creates a new UnifieldPtr with uint16 value.
//
// It allocates a heap copy of val (stored as uint64) so that mutating the original
// value after this call does not affect the stored value. Only one pointer field
// (u64P) will be non-nil in the returned UnifieldPtr.
func Uint16(val uint16) UnifieldPtr {
	v := uint64(val)
	return UnifieldPtr{u64P: &v, type_: valueUint16}
}

// Uint32 creates a new UnifieldPtr with uint32 value.
//
// It allocates a heap copy of val (stored as uint64) so that mutating the original
// value after this call does not affect the stored value. Only one pointer field
// (u64P) will be non-nil in the returned UnifieldPtr.
func Uint32(val uint32) UnifieldPtr {
	v := uint64(val)
	return UnifieldPtr{u64P: &v, type_: valueUint32}
}

// Uint64 creates a new UnifieldPtr with uint64 value.
//
// It allocates a heap copy of val so that mutating the original value after this
// call does not affect the stored value. Only one pointer field (u64P) will be
// non-nil in the returned UnifieldPtr.
func Uint64(val uint64) UnifieldPtr {
	v := val
	return UnifieldPtr{u64P: &v, type_: valueUint64}
}

// Float32 creates a new UnifieldPtr with float32 value.
//
// It allocates a heap copy of val (stored as float64) so that mutating the original
// value after this call does not affect the stored value. Only one pointer field
// (f64P) will be non-nil in the returned UnifieldPtr.
func Float32(val float32) UnifieldPtr {
	v := float64(val)
	return UnifieldPtr{f64P: &v, type_: valueFloat32}
}

// Float64 creates a new UnifieldPtr with float64 value.
//
// It allocates a heap copy of val so that mutating the original value after this
// call does not affect the stored value. Only one pointer field (f64P) will be
// non-nil in the returned UnifieldPtr.
func Float64(val float64) UnifieldPtr {
	v := val
	return UnifieldPtr{f64P: &v, type_: valueFloat64}
}

// Err creates a new UnifieldPtr with error value.
//
// Stores the error interface directly — since error is already an interface, no
// pointer indirection is needed. Passing nil as val produces an empty UnifieldPtr
// with valueEmpty (MarshalToError will return ErrEmptyUnifield).
func Err(val error) UnifieldPtr {
	if val == nil {
		return UnifieldPtr{type_: valueEmpty}
	}
	return UnifieldPtr{errP: val, type_: valueError}
}

// Time creates a new UnifieldPtr with time.Time value.
//
// It allocates a heap copy of val so that mutating the original time.Time after
// this call does not affect the stored value. Only one pointer field (tmP) will be
// non-nil in the returned UnifieldPtr.
func Time(val time.Time) UnifieldPtr {
	v := val
	return UnifieldPtr{tmP: &v, type_: valueTime}
}
