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

package unifields

import (
	"errors"
	"fmt"
	"time"
)

// valueType identifies the stored datatype inside a Unifield instance.
type valueType uint8

const (
	valueEmpty   valueType = iota // no value stored
	valueString                   // str field holds string value
	valueInt                      // i64 field holds signed int value
	valueInt8                     // i64 field holds int8 value
	valueInt16                    // i64 field holds int16 value
	valueInt32                    // i64 field holds int32 value
	valueInt64                    // i64 field holds int64 value
	valueUint                     // u64 field holds unsigned int value
	valueUint8                    // u64 field holds uint8 value
	valueUint16                   // u64 field holds uint16 value
	valueUint32                   // u64 field holds uint32 value
	valueUint64                   // u64 field holds uint64 value
	valueFloat32                  // f64 field holds float32 value
	valueFloat64                  // f64 field holds float64 value
	valueError                    // err field holds error value
	valueTime                     // tm field holds time.Time value
)

// Static error values for common failure cases.
var (
	ErrDstNil        = errors.New("dst is nil")
	ErrTypeMismatch  = errors.New("type mismatch")
	ErrEmptyUnifield = errors.New("unifield is empty")
)

// Unifield holds exactly one typed value paired with a key identifier.
// Uses flat storage with a valueType discriminator tag — similar to zapcore.Field.
// All internal fields are value-types (except error which is an interface), enabling
// cheap copy semantics via Clone().
type Unifield struct {
	key   string    // identifier key for the field
	err   error     // error value storage
	tm    time.Time // time.Time value storage
	str   string    // string value storage
	i64   int64     // signed integer storage (int, int8..int64)
	u64   uint64    // unsigned integer storage (uint, uint8..uint64)
	f64   float64   // floating point storage (float32, float64)
	type_ valueType // discriminator tag
}

// String creates a new Unifield with key and string value.
func String(key string, val string) Unifield {
	return Unifield{key: key, str: val, type_: valueString}
}

// Int creates a new Unifield with key and int value.
func Int(key string, val int) Unifield {
	return Unifield{key: key, i64: int64(val), type_: valueInt}
}

// Int8 creates a new Unifield with key and int8 value.
func Int8(key string, val int8) Unifield {
	return Unifield{key: key, i64: int64(val), type_: valueInt8}
}

// Int16 creates a new Unifield with key and int16 value.
func Int16(key string, val int16) Unifield {
	return Unifield{key: key, i64: int64(val), type_: valueInt16}
}

// Int32 creates a new Unifield with key and int32 value.
func Int32(key string, val int32) Unifield {
	return Unifield{key: key, i64: int64(val), type_: valueInt32}
}

// Int64 creates a new Unifield with key and int64 value.
func Int64(key string, val int64) Unifield {
	return Unifield{key: key, i64: val, type_: valueInt64}
}

// Uint creates a new Unifield with key and uint value.
func Uint(key string, val uint) Unifield {
	return Unifield{key: key, u64: uint64(val), type_: valueUint}
}

// Uint8 creates a new Unifield with key and uint8 value.
func Uint8(key string, val uint8) Unifield {
	return Unifield{key: key, u64: uint64(val), type_: valueUint8}
}

// Uint16 creates a new Unifield with key and uint16 value.
func Uint16(key string, val uint16) Unifield {
	return Unifield{key: key, u64: uint64(val), type_: valueUint16}
}

// Uint32 creates a new Unifield with key and uint32 value.
func Uint32(key string, val uint32) Unifield {
	return Unifield{key: key, u64: uint64(val), type_: valueUint32}
}

// Uint64 creates a new Unifield with key and uint64 value.
func Uint64(key string, val uint64) Unifield {
	return Unifield{key: key, u64: val, type_: valueUint64}
}

// Float32 creates a new Unifield with key and float32 value.
func Float32(key string, val float32) Unifield {
	return Unifield{key: key, f64: float64(val), type_: valueFloat32}
}

// Float64 creates a new Unifield with key and float64 value.
func Float64(key string, val float64) Unifield {
	return Unifield{key: key, f64: val, type_: valueFloat64}
}

// Err creates a new Unifield with key and error value.
func Err(key string, val error) Unifield {
	return Unifield{key: key, err: val, type_: valueError}
}

// Time creates a new Unifield with key and time.Time value.
func Time(key string, val time.Time) Unifield {
	return Unifield{key: key, tm: val, type_: valueTime}
}

// Clone returns a shallow copy of the Unifield. Since all fields except error
// (an interface) are pure values, this preserves all stored values correctly.
func (u Unifield) Clone() Unifield {
	return u
}

// --- MarshalTo<Typed> deserialization methods ---

// MarshalToStr copies the stored string value into dst.
// Returns an error if the Unifield is empty or holds a different type.
func (u Unifield) MarshalToStr(dst *string) error {
	if dst == nil {
		return ErrDstNil
	}
	if u.type_ != valueString {
		return fmt.Errorf("%w: got %s, want str", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = u.str

	return nil
}

// MarshalToInt copies the stored signed integer value into dst as int.
func (u Unifield) MarshalToInt(dst *int) error {
	if dst == nil {
		return ErrDstNil
	}
	if !isSignedInt(u.type_) {
		return fmt.Errorf("%w: got %s, want int", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = int(u.i64)

	return nil
}

// MarshalToInt8 copies the stored signed integer value into dst as int8.
func (u Unifield) MarshalToInt8(dst *int8) error {
	if dst == nil {
		return ErrDstNil
	}
	if !isSignedInt(u.type_) {
		return fmt.Errorf("%w: got %s, want int8", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = int8(u.i64)

	return nil
}

// MarshalToInt16 copies the stored signed integer value into dst as int16.
func (u Unifield) MarshalToInt16(dst *int16) error {
	if dst == nil {
		return ErrDstNil
	}
	if !isSignedInt(u.type_) {
		return fmt.Errorf("%w: got %s, want int16", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = int16(u.i64)

	return nil
}

// MarshalToInt32 copies the stored signed integer value into dst as int32.
func (u Unifield) MarshalToInt32(dst *int32) error {
	if dst == nil {
		return ErrDstNil
	}
	if !isSignedInt(u.type_) {
		return fmt.Errorf("%w: got %s, want int32", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = int32(u.i64)

	return nil
}

// MarshalToInt64 copies the stored signed integer value into dst as int64.
func (u Unifield) MarshalToInt64(dst *int64) error {
	if dst == nil {
		return ErrDstNil
	}
	if !isSignedInt(u.type_) {
		return fmt.Errorf("%w: got %s, want int64", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = u.i64

	return nil
}

// MarshalToUint copies the stored unsigned integer value into dst as uint.
func (u Unifield) MarshalToUint(dst *uint) error {
	if dst == nil {
		return ErrDstNil
	}
	if !isUnsignedInt(u.type_) {
		return fmt.Errorf("%w: got %s, want uint", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = uint(u.u64)

	return nil
}

// MarshalToUint8 copies the stored unsigned integer value into dst as uint8.
func (u Unifield) MarshalToUint8(dst *uint8) error {
	if dst == nil {
		return ErrDstNil
	}
	if !isUnsignedInt(u.type_) {
		return fmt.Errorf("%w: got %s, want uint8", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = uint8(u.u64)

	return nil
}

// MarshalToUint16 copies the stored unsigned integer value into dst as uint16.
func (u Unifield) MarshalToUint16(dst *uint16) error {
	if dst == nil {
		return ErrDstNil
	}
	if !isUnsignedInt(u.type_) {
		return fmt.Errorf("%w: got %s, want uint16", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = uint16(u.u64)

	return nil
}

// MarshalToUint32 copies the stored unsigned integer value into dst as uint32.
func (u Unifield) MarshalToUint32(dst *uint32) error {
	if dst == nil {
		return ErrDstNil
	}
	if !isUnsignedInt(u.type_) {
		return fmt.Errorf("%w: got %s, want uint32", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = uint32(u.u64)

	return nil
}

// MarshalToUint64 copies the stored unsigned integer value into dst as uint64.
func (u Unifield) MarshalToUint64(dst *uint64) error {
	if dst == nil {
		return ErrDstNil
	}
	if !isUnsignedInt(u.type_) {
		return fmt.Errorf("%w: got %s, want uint64", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = u.u64

	return nil
}

// MarshalToFloat32 copies the stored floating-point value into dst as float32.
// Accepts both float32 and float64 sources since they share the same backing field.
func (u Unifield) MarshalToFloat32(dst *float32) error {
	if dst == nil {
		return ErrDstNil
	}
	if u.type_ != valueFloat32 && u.type_ != valueFloat64 {
		return fmt.Errorf("%w: got %s, want float32", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = float32(u.f64)

	return nil
}

// MarshalToFloat64 copies the stored floating-point value into dst as float64.
// Accepts both float32 and float64 sources since they share the same backing field.
func (u Unifield) MarshalToFloat64(dst *float64) error {
	if dst == nil {
		return ErrDstNil
	}
	if u.type_ != valueFloat32 && u.type_ != valueFloat64 {
		return fmt.Errorf("%w: got %s, want float64", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = u.f64

	return nil
}

// MarshalToError copies the stored error value into dst.
func (u Unifield) MarshalToError(dst *error) error {
	if dst == nil {
		return ErrDstNil
	}
	if u.type_ != valueError {
		return fmt.Errorf("%w: got %s, want error", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = u.err

	return nil
}

// MarshalToTime copies the stored time.Time value into dst.
func (u Unifield) MarshalToTime(dst *time.Time) error {
	if dst == nil {
		return ErrDstNil
	}
	if u.type_ != valueTime {
		return fmt.Errorf("%w: got %s, want time.Time", ErrTypeMismatch, typeName(u.type_))
	}
	*dst = u.tm

	return nil
}

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
