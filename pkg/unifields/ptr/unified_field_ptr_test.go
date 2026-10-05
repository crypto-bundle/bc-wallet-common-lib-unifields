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
// 3. Any person or organization found to be in violation of these restrictions will be held liable
// for any damages resulting from such use.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package ptr

import (
	"errors"
	"testing"
	"time"
)

var errTest = errors.New("test error value")

// --- Factory tests: string ---

func TestStringFactory(t *testing.T) {
	t.Parallel()

	key := "test_key"
	val := "hello"

	u := String(key, val)

	if u.key != key {
		t.Errorf("key = %q, want %q", u.key, key)
	}
	if u.type_ != valueString {
		t.Errorf("type_ = %v, want valueString", u.type_)
	}
	if u.strP == nil {
		t.Fatal("strP is nil")
	}
	if *u.strP != val {
		t.Errorf("*strP = %q, want %q", *u.strP, val)
	}
}

func TestStringImmutability(t *testing.T) {
	t.Parallel()

	val := "original"
	u := String("k", val)

	// Mutate original after factory call.
	val = " mutated" //nolint:ineffassign,wastedassign // verifies immutability: mutation must NOT affect stored value

	if *u.strP != "original" {
		t.Errorf("*strP = %q, want \"original\" — factory must copy input", *u.strP)
	}
}

func TestStringZeroValue(t *testing.T) {
	t.Parallel()

	u := String("k", "")

	if u.type_ != valueString {
		t.Errorf("type_ = %v, want valueString", u.type_)
	}
	if u.strP == nil || *u.strP != "" {
		t.Errorf("*strP = %v, want empty string", u.strP)
	}
}

// --- Factory tests: signed integers ---

func checkIntStorage(t *testing.T, u UnifieldPtr, expected int64, wantType valueType) {
	t.Helper()

	if u.key != "k" {
		t.Errorf("key = %q, want \"k\"", u.key)
	}
	if u.type_ != wantType {
		t.Errorf("type_ = %v, want %v", u.type_, wantType)
	}
	if u.i64P == nil {
		t.Fatal("i64P is nil")
	}
	if *u.i64P != expected {
		t.Errorf("*i64P = %d, want %d", *u.i64P, expected)
	}
	if u.u64P != nil || u.f64P != nil || u.strP != nil || u.tmP != nil {
		t.Error("unexpected non-nil pointer fields")
	}
}

func TestIntFactory(t *testing.T) {
	t.Parallel()

	const val = 42
	u := Int("k", val)

	checkIntStorage(t, u, int64(val), valueInt)
}

func TestInt8Factory(t *testing.T) {
	t.Parallel()

	const val = int8(127)
	u := Int8("k", val)

	checkIntStorage(t, u, int64(val), valueInt8)
}

func TestInt16Factory(t *testing.T) {
	t.Parallel()

	const val = int16(-100)
	u := Int16("k", val)

	checkIntStorage(t, u, int64(val), valueInt16)
}

func TestInt32Factory(t *testing.T) {
	t.Parallel()

	const val = int32(999999)
	u := Int32("k", val)

	checkIntStorage(t, u, int64(val), valueInt32)
}

func TestInt64Factory(t *testing.T) {
	t.Parallel()

	const val = int64(-9223372036854775808) // min int64
	u := Int64("k", val)

	checkIntStorage(t, u, val, valueInt64)
}

func TestIntImmutability(t *testing.T) {
	t.Parallel()

	val := int(100)
	u := Int("k", val)
	val = 200 //nolint:ineffassign,wastedassign // verifies immutability: mutation must NOT affect stored value

	if *u.i64P != 100 {
		t.Errorf("*i64P = %d, want 100 — factory must copy input", *u.i64P)
	}
}

func TestIntZeroValue(t *testing.T) {
	t.Parallel()

	u := Int("k", 0)

	if u.type_ != valueInt {
		t.Errorf("type_ = %v, want valueInt", u.type_)
	}
	if u.i64P == nil || *u.i64P != 0 {
		t.Errorf("*i64P = %v, want 0", u.i64P)
	}
}

// --- Factory tests: unsigned integers ---

func checkUintStorage(t *testing.T, u UnifieldPtr, expected uint64, wantType valueType) {
	t.Helper()

	if u.key != "k" {
		t.Errorf("key = %q, want \"k\"", u.key)
	}
	if u.type_ != wantType {
		t.Errorf("type_ = %v, want %v", u.type_, wantType)
	}
	if u.u64P == nil {
		t.Fatal("u64P is nil")
	}
	if *u.u64P != expected {
		t.Errorf("*u64P = %d, want %d", *u.u64P, expected)
	}
	if u.i64P != nil || u.f64P != nil || u.strP != nil || u.tmP != nil {
		t.Error("unexpected non-nil pointer fields")
	}
}

func TestUintFactory(t *testing.T) {
	t.Parallel()

	const val = uint(12345)
	u := Uint("k", val)

	checkUintStorage(t, u, uint64(val), valueUint)
}

func TestUint8Factory(t *testing.T) {
	t.Parallel()

	const val = uint8(255)
	u := Uint8("k", val)

	checkUintStorage(t, u, uint64(val), valueUint8)
}

func TestUint16Factory(t *testing.T) {
	t.Parallel()

	const val = uint16(65535)
	u := Uint16("k", val)

	checkUintStorage(t, u, uint64(val), valueUint16)
}

func TestUint32Factory(t *testing.T) {
	t.Parallel()

	const val = uint32(4294967295)
	u := Uint32("k", val)

	checkUintStorage(t, u, uint64(val), valueUint32)
}

func TestUint64Factory(t *testing.T) {
	t.Parallel()

	const val = uint64(18446744073709551615) // max uint64
	u := Uint64("k", val)

	checkUintStorage(t, u, val, valueUint64)
}

func TestUintImmutability(t *testing.T) {
	t.Parallel()

	val := uint(999)
	u := Uint("k", val)
	val = 777 //nolint:ineffassign,wastedassign // verifies immutability: mutation must NOT affect stored value

	if *u.u64P != 999 {
		t.Errorf("*u64P = %d, want 999 — factory must copy input", *u.u64P)
	}
}

func TestUintZeroValue(t *testing.T) {
	t.Parallel()

	u := Uint("k", 0)

	if u.type_ != valueUint {
		t.Errorf("type_ = %v, want valueUint", u.type_)
	}
	if u.u64P == nil || *u.u64P != 0 {
		t.Errorf("*u64P = %v, want 0", u.u64P)
	}
}

// --- Factory tests: floats ---

func TestFloat32Factory(t *testing.T) {
	t.Parallel()

	val := float32(3.14159)
	u := Float32("k", val)

	if u.key != "k" {
		t.Errorf("key = %q, want \"k\"", u.key)
	}
	if u.type_ != valueFloat32 {
		t.Errorf("type_ = %v, want valueFloat32", u.type_)
	}
	if u.f64P == nil {
		t.Fatal("f64P is nil")
	}
	if float32(*u.f64P) != val {
		t.Errorf("float32(*f64P) = %f, want %f", float32(*u.f64P), val)
	}
	if u.i64P != nil || u.u64P != nil || u.strP != nil || u.tmP != nil {
		t.Error("unexpected non-nil pointer fields")
	}
}

func TestFloat64Factory(t *testing.T) {
	t.Parallel()

	val := 2.99792458e8 // speed of light
	u := Float64("k", val)

	if u.key != "k" {
		t.Errorf("key = %q, want \"k\"", u.key)
	}
	if u.type_ != valueFloat64 {
		t.Errorf("type_ = %v, want valueFloat64", u.type_)
	}
	if u.f64P == nil {
		t.Fatal("f64P is nil")
	}
	if *u.f64P != val {
		t.Errorf("*f64P = %g, want %g", *u.f64P, val)
	}
	if u.i64P != nil || u.u64P != nil || u.strP != nil || u.tmP != nil {
		t.Error("unexpected non-nil pointer fields")
	}
}

func TestFloat32Immutability(t *testing.T) {
	t.Parallel()

	val := float32(1.5)
	u := Float32("k", val)
	val = 9.9 //nolint:ineffassign,wastedassign // verifies immutability: mutation must NOT affect stored value

	if float32(*u.f64P) != 1.5 {
		t.Errorf("float32(*f64P) = %f, want 1.5 — factory must copy input", float32(*u.f64P))
	}
}

func TestFloat32ZeroValue(t *testing.T) {
	t.Parallel()

	u := Float32("k", 0)

	if u.type_ != valueFloat32 {
		t.Errorf("type_ = %v, want valueFloat32", u.type_)
	}
	if u.f64P == nil || float32(*u.f64P) != 0 {
		t.Errorf("*f64P = %v, want 0", u.f64P)
	}
}

// --- Factory tests: error ---

func TestErrFactory(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("test error")
	u := Err("k", wantErr)

	if u.key != "k" {
		t.Errorf("key = %q, want \"k\"", u.key)
	}
	if u.type_ != valueError {
		t.Errorf("type_ = %v, want valueError", u.type_)
	}
	if u.errP == nil {
		t.Error("errP is nil")
	}
	if u.errP.Error() != wantErr.Error() {
		t.Errorf("errP.Error() = %q, want %q", u.errP.Error(), wantErr.Error())
	}
	// Other pointers must be nil.
	if u.strP != nil || u.tmP != nil || u.i64P != nil || u.u64P != nil || u.f64P != nil {
		t.Error("unexpected non-nil pointer fields on error UnifieldPtr")
	}
}

func TestErrNilIsEmpty(t *testing.T) {
	t.Parallel()

	u := Err("k", nil)

	if u.type_ != valueEmpty {
		t.Errorf("type_ = %v, want valueEmpty for nil error", u.type_)
	}
}

// --- Factory tests: time ---

func TestTimeFactory(t *testing.T) {
	t.Parallel()

	want := time.Date(2025, 6, 15, 10, 30, 45, 123456789, time.UTC)
	u := Time("k", want)

	if u.key != "k" {
		t.Errorf("key = %q, want \"k\"", u.key)
	}
	if u.type_ != valueTime {
		t.Errorf("type_ = %v, want valueTime", u.type_)
	}
	if u.tmP == nil {
		t.Fatal("tmP is nil")
	}
	if !u.tmP.Equal(want) {
		t.Errorf("*tmP = %v, want %v", u.tmP, want)
	}
	if u.strP != nil || u.i64P != nil || u.u64P != nil || u.f64P != nil {
		t.Error("unexpected non-nil pointer fields on time UnifieldPtr")
	}
}

func TestTimeImmutability(t *testing.T) {
	t.Parallel()

	want := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	u := Time("k", want)

	// Mutate original: Go's time.Time is a struct so reassignment works.
	wants := want.Add(time.Hour)
	if !u.tmP.Equal(want) {
		t.Errorf("*tmP = %v, want %v — factory must copy input", u.tmP, want)
	}
	if u.tmP.Equal(wants) {
		t.Error("factory captured caller's reference instead of copying")
	}
}

func TestTimeZeroValue(t *testing.T) {
	t.Parallel()

	u := Time("k", time.Time{})

	if u.type_ != valueTime {
		t.Errorf("type_ = %v, want valueTime", u.type_)
	}
	if u.tmP == nil {
		t.Fatal("tmP is nil")
	}
	if !u.tmP.IsZero() {
		t.Errorf("*tmP is not zero: %v", u.tmP)
	}
}

// --- Clone separation test ---

func TestCloneSeparation(t *testing.T) {
	t.Parallel()

	ptrSrc := Int("k", 100)
	ptrCloned := ptrSrc.Clone()

	// Both point to same heap int64, but the struct itself is copied.
	// Since we don't expose mutable APIs, Clone separation here means
	// the struct-level fields (key, type_) are independent copies.
	clonedVal := ptrCloned.(UnifieldPtr) //nolint:forcetypeassert // Clone() always returns concrete UnifieldPtr wrapping itself
	if ptrSrc.key != clonedVal.key {
		t.Errorf("cloned key = %q, want %q", clonedVal.key, ptrSrc.key)
	}
	if ptrSrc.type_ != clonedVal.type_ {
		t.Errorf("cloned type_ = %v, want %v", clonedVal.type_, ptrSrc.type_)
	}
}

// ================================================================
// MarshalTo<T> positive tests — one per supported type
// ================================================================

func TestMarshalToStr(t *testing.T) {
	t.Parallel()

	u := String("k", "hello")
	var got string
	if err := u.MarshalToStr(&got); err != nil {
		t.Fatalf("MarshalToStr failed: %v", err)
	}
	if got != "hello" {
		t.Errorf("got %q, want %q", got, "hello")
	}
}

func TestMarshalToInt(t *testing.T) {
	t.Parallel()

	u := Int("k", 42)
	var got int
	if err := u.MarshalToInt(&got); err != nil {
		t.Fatalf("MarshalToInt failed: %v", err)
	}
	if got != 42 {
		t.Errorf("got %d, want 42", got)
	}
}

func TestMarshalToInt8(t *testing.T) {
	t.Parallel()

	u := Int8("k", 127)
	var got int8
	if err := u.MarshalToInt8(&got); err != nil {
		t.Fatalf("MarshalToInt8 failed: %v", err)
	}
	if got != 127 {
		t.Errorf("got %d, want 127", got)
	}
}

func TestMarshalToInt16(t *testing.T) {
	t.Parallel()

	u := Int16("k", -32768)
	var got int16
	if err := u.MarshalToInt16(&got); err != nil {
		t.Fatalf("MarshalToInt16 failed: %v", err)
	}
	if got != -32768 {
		t.Errorf("got %d, want -32768", got)
	}
}

func TestMarshalToInt32(t *testing.T) {
	t.Parallel()

	u := Int32("k", 123456789)
	var got int32
	if err := u.MarshalToInt32(&got); err != nil {
		t.Fatalf("MarshalToInt32 failed: %v", err)
	}
	if got != 123456789 {
		t.Errorf("got %d, want 123456789", got)
	}
}

func TestMarshalToInt64(t *testing.T) {
	t.Parallel()

	u := Int64("k", 9223372036854775807)
	var got int64
	if err := u.MarshalToInt64(&got); err != nil {
		t.Fatalf("MarshalToInt64 failed: %v", err)
	}
	if got != 9223372036854775807 {
		t.Errorf("got %d, want 9223372036854775807", got)
	}
}

func TestMarshalToUint(t *testing.T) {
	t.Parallel()

	u := Uint("k", 12345)
	var got uint
	if err := u.MarshalToUint(&got); err != nil {
		t.Fatalf("MarshalToUint failed: %v", err)
	}
	if got != 12345 {
		t.Errorf("got %d, want 12345", got)
	}
}

func TestMarshalToUint8(t *testing.T) {
	t.Parallel()

	u := Uint8("k", 255)
	var got uint8
	if err := u.MarshalToUint8(&got); err != nil {
		t.Fatalf("MarshalToUint8 failed: %v", err)
	}
	if got != 255 {
		t.Errorf("got %d, want 255", got)
	}
}

func TestMarshalToUint16(t *testing.T) {
	t.Parallel()

	u := Uint16("k", 65535)
	var got uint16
	if err := u.MarshalToUint16(&got); err != nil {
		t.Fatalf("MarshalToUint16 failed: %v", err)
	}
	if got != 65535 {
		t.Errorf("got %d, want 65535", got)
	}
}

func TestMarshalToUint32(t *testing.T) {
	t.Parallel()

	u := Uint32("k", 4294967295)
	var got uint32
	if err := u.MarshalToUint32(&got); err != nil {
		t.Fatalf("MarshalToUint32 failed: %v", err)
	}
	if got != 4294967295 {
		t.Errorf("got %d, want 4294967295", got)
	}
}

func TestMarshalToUint64(t *testing.T) {
	t.Parallel()

	u := Uint64("k", 18446744073709551615)
	var got uint64
	if err := u.MarshalToUint64(&got); err != nil {
		t.Fatalf("MarshalToUint64 failed: %v", err)
	}
	if got != 18446744073709551615 {
		t.Errorf("got %d, want 18446744073709551615", got)
	}
}

func TestMarshalToFloat32(t *testing.T) {
	t.Parallel()

	u := Float32("k", 3.14159)
	var got float32
	if err := u.MarshalToFloat32(&got); err != nil {
		t.Fatalf("MarshalToFloat32 failed: %v", err)
	}
	const eps = 1e-6
	if absDiff(float64(got), 3.14159) > eps {
		t.Errorf("got %v, want 3.14159", got)
	}
}

func TestMarshalToFloat64(t *testing.T) {
	t.Parallel()

	u := Float64("k", 2.718281828)
	var got float64
	if err := u.MarshalToFloat64(&got); err != nil {
		t.Fatalf("MarshalToFloat64 failed: %v", err)
	}
	if absDiff(got, 2.718281828) > 1e-6 {
		t.Errorf("got %v, want 2.718281828", got)
	}
}

func TestMarshalToError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("test error")
	u := Err("k", wantErr)
	var got error
	if err := u.MarshalToError(&got); err != nil {
		t.Fatalf("MarshalToError failed: %v", err)
	}
	if !errors.Is(got, wantErr) {
		t.Errorf("got %v, want %v", got, wantErr)
	}
}

func TestMarshalToTime(t *testing.T) {
	t.Parallel()

	want := time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC)
	u := Time("k", want)
	var got time.Time
	if err := u.MarshalToTime(&got); err != nil {
		t.Fatalf("MarshalToTime failed: %v", err)
	}
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// ================================================================
// MarshalTo<T> negative tests — type mismatch per supported type
// ================================================================

func TestMismatchStringToInt(t *testing.T) {
	t.Parallel()
	u := String("k", "v")
	var i int
	if err := u.MarshalToInt(&i); err == nil {
		t.Error("expected error for String → *int")
	}
}

func TestMismatchStringToInt8(t *testing.T) {
	t.Parallel()
	u := String("k", "v")
	var i int8
	if err := u.MarshalToInt8(&i); err == nil {
		t.Error("expected error for String → *int8")
	}
}

func TestMismatchStringToInt16(t *testing.T) {
	t.Parallel()
	u := String("k", "v")
	var i int16
	if err := u.MarshalToInt16(&i); err == nil {
		t.Error("expected error for String → *int16")
	}
}

func TestMismatchStringToInt32(t *testing.T) {
	t.Parallel()
	u := String("k", "v")
	var i int32
	if err := u.MarshalToInt32(&i); err == nil {
		t.Error("expected error for String → *int32")
	}
}

func TestMismatchStringToInt64(t *testing.T) {
	t.Parallel()
	u := String("k", "v")
	var i int64
	if err := u.MarshalToInt64(&i); err == nil {
		t.Error("expected error for String → *int64")
	}
}

func TestMismatchStringToUint(t *testing.T) {
	t.Parallel()
	u := String("k", "v")
	var u64 uint
	if err := u.MarshalToUint(&u64); err == nil {
		t.Error("expected error for String → *uint")
	}
}

func TestMismatchStringToUint8(t *testing.T) {
	t.Parallel()
	u := String("k", "v")
	var u8 uint8
	if err := u.MarshalToUint8(&u8); err == nil {
		t.Error("expected error for String → *uint8")
	}
}

func TestMismatchStringToUint16(t *testing.T) {
	t.Parallel()
	u := String("k", "v")
	var u16 uint16
	if err := u.MarshalToUint16(&u16); err == nil {
		t.Error("expected error for String → *uint16")
	}
}

func TestMismatchStringToUint32(t *testing.T) {
	t.Parallel()
	u := String("k", "v")
	var u32 uint32
	if err := u.MarshalToUint32(&u32); err == nil {
		t.Error("expected error for String → *uint32")
	}
}

func TestMismatchStringToUint64(t *testing.T) {
	t.Parallel()
	u := String("k", "v")
	var u64 uint64
	if err := u.MarshalToUint64(&u64); err == nil {
		t.Error("expected error for String → *uint64")
	}
}

func TestMismatchStringToFloat32(t *testing.T) {
	t.Parallel()
	u := String("k", "v")
	var f float32
	if err := u.MarshalToFloat32(&f); err == nil {
		t.Error("expected error for String → *float32")
	}
}

func TestMismatchStringToFloat64(t *testing.T) {
	t.Parallel()
	u := String("k", "v")
	var f float64
	if err := u.MarshalToFloat64(&f); err == nil {
		t.Error("expected error for String → *float64")
	}
}

func TestMismatchStringToError(t *testing.T) {
	t.Parallel()
	u := String("k", "v")
	var e error
	if err := u.MarshalToError(&e); err == nil {
		t.Error("expected error for String → *error")
	}
}

func TestMismatchStringToTime(t *testing.T) {
	t.Parallel()
	u := String("k", "v")
	var gotTime time.Time
	if err := u.MarshalToTime(&gotTime); err == nil {
		t.Error("expected error for String → *time.Time")
	}
}

func TestMismatchIntToString(t *testing.T) {
	t.Parallel()
	u := Int("k", 1)
	var s string
	if err := u.MarshalToStr(&s); err == nil {
		t.Error("expected error for Int → *string")
	}
}

func TestMismatchInt8ToInt(t *testing.T) {
	t.Parallel()
	u := Int8("k", 1)
	// Int8 shares i64P backing with int — MarshalToInt accepts all signed ints (permissive group matching).
	var i int
	if err := u.MarshalToInt(&i); err != nil {
		t.Fatalf("expected success for Int8 → *int (both backed by i64P), got error: %v", err)
	}
	if i != 1 {
		t.Errorf("got %d, want 1", i)
	}
}

func TestMismatchIntToInt64(t *testing.T) {
	t.Parallel()
	u := Int("k", 1)
	// Int shares i64P backing with int64 — MarshalToInt64 accepts all signed ints (permissive group matching).
	var i int64
	if err := u.MarshalToInt64(&i); err != nil {
		t.Fatalf("expected success for Int → *int64 (both backed by i64P), got error: %v", err)
	}
	if i != 1 {
		t.Errorf("got %d, want 1", i)
	}
}

func TestMismatchUintToInt64(t *testing.T) {
	t.Parallel()
	u := Uint("k", 1)
	var i int64
	if err := u.MarshalToInt64(&i); err == nil {
		t.Error("expected error for Uint → *int64 (signed/unsigned)")
	}
}

func TestMismatchFloat64ToInt(t *testing.T) {
	t.Parallel()
	u := Float64("k", 1.0)
	var i int
	if err := u.MarshalToInt(&i); err == nil {
		t.Error("expected error for Float64 → *int")
	}
}

func TestMismatchErrorToInt(t *testing.T) {
	t.Parallel()
	u := Err("k", errTest)
	var i int
	if err := u.MarshalToInt(&i); err == nil {
		t.Error("expected error for Err → *int")
	}
}

func TestMismatchTimeToInt(t *testing.T) {
	t.Parallel()
	u := Time("k", time.Now())
	var i int
	if err := u.MarshalToInt(&i); err == nil {
		t.Error("expected error for Time → *int")
	}
}

func TestMismatchInt64ToUint(t *testing.T) {
	t.Parallel()
	u := Int64("k", 42)
	var ui uint64
	if err := u.MarshalToUint64(&ui); err == nil {
		t.Error("expected error for Int64 → *uint64")
	}
}

func TestMismatchUint64ToInt64(t *testing.T) {
	t.Parallel()
	u := Uint64("k", 42)
	var i int64
	if err := u.MarshalToInt64(&i); err == nil {
		t.Error("expected error for Uint64 → *int64")
	}
}

func TestMismatchFloat32ToFloat64DifferentTag(t *testing.T) {
	t.Parallel()
	u := Float32("k", 1.5)
	// Float32 and Float64 share f64P backing — MarshalToFloat64 accepts all floats (permissive group matching).
	var f float64
	if err := u.MarshalToFloat64(&f); err != nil {
		t.Fatalf("expected success for Float32 → *float64 (both backed by f64P), got error: %v", err)
	}
	const eps = 1e-6
	if absDiff(f, 1.5) > eps {
		t.Errorf("got %v, want 1.5", f)
	}
}

// ================================================================
// MarshalTo<T> negative tests — nil destination
// ================================================================

func TestNilDstMarshalToStr(t *testing.T) {
	t.Parallel()
	u := String("k", "v")
	if err := u.MarshalToStr(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestNilDstMarshalToInt(t *testing.T) {
	t.Parallel()
	u := Int("k", 42)
	if err := u.MarshalToInt(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestNilDstMarshalToInt8(t *testing.T) {
	t.Parallel()
	u := Int8("k", 1)
	if err := u.MarshalToInt8(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestNilDstMarshalToInt16(t *testing.T) {
	t.Parallel()
	u := Int16("k", 1)
	if err := u.MarshalToInt16(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestNilDstMarshalToInt32(t *testing.T) {
	t.Parallel()
	u := Int32("k", 1)
	if err := u.MarshalToInt32(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestNilDstMarshalToInt64(t *testing.T) {
	t.Parallel()
	u := Int64("k", 1)
	if err := u.MarshalToInt64(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestNilDstMarshalToUint(t *testing.T) {
	t.Parallel()
	u := Uint("k", 1)
	if err := u.MarshalToUint(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestNilDstMarshalToUint8(t *testing.T) {
	t.Parallel()
	u := Uint8("k", 1)
	if err := u.MarshalToUint8(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestNilDstMarshalToUint16(t *testing.T) {
	t.Parallel()
	u := Uint16("k", 1)
	if err := u.MarshalToUint16(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestNilDstMarshalToUint32(t *testing.T) {
	t.Parallel()
	u := Uint32("k", 1)
	if err := u.MarshalToUint32(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestNilDstMarshalToUint64(t *testing.T) {
	t.Parallel()
	u := Uint64("k", 1)
	if err := u.MarshalToUint64(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestNilDstMarshalToFloat32(t *testing.T) {
	t.Parallel()
	u := Float32("k", 1.0)
	if err := u.MarshalToFloat32(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestNilDstMarshalToFloat64(t *testing.T) {
	t.Parallel()
	u := Float64("k", 1.0)
	if err := u.MarshalToFloat64(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestNilDstMarshalToError(t *testing.T) {
	t.Parallel()
	u := Err("k", errTest)
	if err := u.MarshalToError(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestNilDstMarshalToTime(t *testing.T) {
	t.Parallel()
	u := Time("k", time.Now())
	if err := u.MarshalToTime(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

// ================================================================
// Empty UnifieldPtr access tests
// ================================================================

func TestEmptyAccessMarshalToStr(t *testing.T) {
	t.Parallel()
	u := UnifieldPtr{key: "k"} //nolint:exhaustruct_v5 // testing empty UnifieldPtr access errors
	var s string
	if err := u.MarshalToStr(&s); err == nil {
		t.Error("expected error for empty → MarshalToStr")
	} else if !errors.Is(err, ErrEmptyUnifield) {
		t.Errorf("expected ErrEmptyUnifield, got %v", err)
	}
}

func TestEmptyAccessMarshalToInt(t *testing.T) {
	t.Parallel()
	u := UnifieldPtr{key: "k"} //nolint:exhaustruct_v5 // testing empty UnifieldPtr access errors
	var i int
	if err := u.MarshalToInt(&i); err == nil {
		t.Error("expected error for empty → MarshalToInt")
	}
}

func TestEmptyAccessMarshalToUint64(t *testing.T) {
	t.Parallel()
	u := UnifieldPtr{key: "k"} //nolint:exhaustruct_v5 // testing empty UnifieldPtr access errors
	var ui uint64
	if err := u.MarshalToUint64(&ui); err == nil {
		t.Error("expected error for empty → MarshalToUint64")
	}
}

func TestEmptyAccessMarshalToFloat64(t *testing.T) {
	t.Parallel()
	u := UnifieldPtr{key: "k"} //nolint:exhaustruct_v5 // testing empty UnifieldPtr access errors
	var f float64
	if err := u.MarshalToFloat64(&f); err == nil {
		t.Error("expected error for empty → MarshalToFloat64")
	}
}

func TestEmptyAccessMarshalToTime(t *testing.T) {
	t.Parallel()
	u := UnifieldPtr{key: "k"} //nolint:exhaustruct_v5 // testing empty UnifieldPtr access errors
	var tm time.Time
	if err := u.MarshalToTime(&tm); err == nil {
		t.Error("expected error for empty → MarshalToTime")
	}
}

func TestNilErrorAccessMarshalToError(t *testing.T) {
	t.Parallel()
	u := Err("k", nil) // nil error → empty
	var e error
	if err := u.MarshalToError(&e); err == nil {
		t.Error("expected error for nil error → MarshalToError")
	}
}

// ================================================================
// Additional edge cases
// ================================================================

func TestKeyPreserved(t *testing.T) {
	t.Parallel()

	u := String("my-key", "val")
	if u.key != "my-key" {
		t.Errorf("key expected 'my-key', got %q", u.key)
	}
	c := u.Clone()
	clonedVal := c.(UnifieldPtr) //nolint:forcetypeassert // Clone() always returns concrete UnifieldPtr wrapping itself
	if clonedVal.key != "my-key" {
		t.Errorf("clone key expected 'my-key', got %q", clonedVal.key)
	}
}

func TestZeroValues(t *testing.T) {
	t.Parallel()

	_ = String("k", "")
	_ = Int("k", 0)
	_ = Int8("k", 0)
	_ = Int16("k", 0)
	_ = Int32("k", 0)
	_ = Int64("k", 0)
	_ = Uint("k", 0)
	_ = Uint8("k", 0)
	_ = Uint16("k", 0)
	_ = Uint32("k", 0)
	_ = Uint64("k", 0)
	_ = Float32("k", 0)
	_ = Float64("k", 0)
	_ = Err("k", nil)
	_ = Time("k", time.Time{})
}

func TestErrKeyStored(t *testing.T) {
	t.Parallel()

	u := Err("my-error-key", errTest)
	if u.key != "my-error-key" {
		t.Errorf("key = %q, want %q", u.key, "my-error-key")
	}
	if !errors.Is(u.errP, errTest) {
		t.Errorf("errP = %v, want %v", u.errP, errTest)
	}
}

func TestTimeKeyAndPointerOnly(t *testing.T) {
	t.Parallel()

	want := time.Date(2024, 1, 1, 0, 0, 0, 0, time.FixedZone("test", 3600))
	u := Time("tz-key", want)

	if u.key != "tz-key" {
		t.Errorf("key = %q, want %q", u.key, "tz-key")
	}
	if u.tmP == nil || !u.tmP.Equal(want) {
		t.Fatal("tmP does not match expected time")
	}
	if u.strP != nil || u.i64P != nil || u.u64P != nil || u.f64P != nil {
		t.Error("unexpected non-nil pointers on time UnifieldPtr")
	}
}

func absDiff(a, b float64) float64 {
	diff := a - b
	if diff < 0 {
		return -diff
	}
	return diff
}
