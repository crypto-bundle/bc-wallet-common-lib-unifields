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

package unifields

import (
	"errors"
	"testing"
	"time"
)

var errTest = errors.New("test error value")

func TestString(t *testing.T) {
	f := String("key", "val")
	if f.key != "key" {
		t.Errorf("expected key=key, got %s", f.key)
	}
	if f.type_ != valueString {
		t.Errorf("expected valueType valueString, got %v", f.type_)
	}
	var s string
	if err := f.MarshalToStr(&s); err != nil {
		t.Fatalf("MarshalToStr failed: %v", err)
	}
	if s != "val" {
		t.Errorf("expected val=val, got %s", s)
	}
	fz := String("k", "")
	var sz string
	if err := fz.MarshalToStr(&sz); err != nil {
		t.Fatalf("MarshalToStr(\"\") failed: %v", err)
	}
	if sz != "" {
		t.Errorf("expected empty, got %q", sz)
	}
}

func TestIntBasic(t *testing.T) {
	f := Int("k", 42)
	if f.type_ != valueInt {
		t.Errorf("type_=want valueInt, got %v", f.type_)
	}
	var got int
	if err := f.MarshalToInt(&got); err != nil {
		t.Fatalf("MarshalToInt failed: %v", err)
	}
	if got != 42 {
		t.Errorf("got %d, want 42", got)
	}
}

func TestInt8(t *testing.T) {
	f := Int8("k", 42)
	if f.type_ != valueInt8 {
		t.Errorf("type_=want valueInt8, got %v", f.type_)
	}
	var got int8
	if err := f.MarshalToInt8(&got); err != nil {
		t.Fatalf("MarshalToInt8 failed: %v", err)
	}
	if got != 42 {
		t.Errorf("got %d, want 42", got)
	}
}

func TestInt16(t *testing.T) {
	f := Int16("k", 12345)
	if f.type_ != valueInt16 {
		t.Errorf("type_=want valueInt16, got %v", f.type_)
	}
	var got int16
	if err := f.MarshalToInt16(&got); err != nil {
		t.Fatalf("MarshalToInt16 failed: %v", err)
	}
	if got != 12345 {
		t.Errorf("got %d, want 12345", got)
	}
}

func TestInt32(t *testing.T) {
	f := Int32("k", 123456789)
	if f.type_ != valueInt32 {
		t.Errorf("type_=want valueInt32, got %v", f.type_)
	}
	var got int32
	if err := f.MarshalToInt32(&got); err != nil {
		t.Fatalf("MarshalToInt32 failed: %v", err)
	}
	if got != 123456789 {
		t.Errorf("got %d, want 123456789", got)
	}
}

func TestInt64(t *testing.T) {
	f := Int64("k", 9223372036854775807)
	if f.type_ != valueInt64 {
		t.Errorf("type_=want valueInt64, got %v", f.type_)
	}
	var got int64
	if err := f.MarshalToInt64(&got); err != nil {
		t.Fatalf("MarshalToInt64 failed: %v", err)
	}
	if got != 9223372036854775807 {
		t.Errorf("got %d, want 9223372036854775807", got)
	}
}

func TestUint(t *testing.T) {
	f := Uint("k", 42)
	if f.type_ != valueUint {
		t.Errorf("type_=want valueUint, got %v", f.type_)
	}
	var got uint
	if err := f.MarshalToUint(&got); err != nil {
		t.Fatalf("MarshalToUint failed: %v", err)
	}
	if got != 42 {
		t.Errorf("got %d, want 42", got)
	}
}

func TestUint8(t *testing.T) {
	f := Uint8("k", 42)
	if f.type_ != valueUint8 {
		t.Errorf("type_=want valueUint8, got %v", f.type_)
	}
	var got uint8
	if err := f.MarshalToUint8(&got); err != nil {
		t.Fatalf("MarshalToUint8 failed: %v", err)
	}
	if got != 42 {
		t.Errorf("got %d, want 42", got)
	}
}

func TestUint16(t *testing.T) {
	f := Uint16("k", 42)
	if f.type_ != valueUint16 {
		t.Errorf("type_=want valueUint16, got %v", f.type_)
	}
	var got uint16
	if err := f.MarshalToUint16(&got); err != nil {
		t.Fatalf("MarshalToUint16 failed: %v", err)
	}
	if got != 42 {
		t.Errorf("got %d, want 42", got)
	}
}

func TestUint32(t *testing.T) {
	f := Uint32("k", 42)
	if f.type_ != valueUint32 {
		t.Errorf("type_=want valueUint32, got %v", f.type_)
	}
	var got uint32
	if err := f.MarshalToUint32(&got); err != nil {
		t.Fatalf("MarshalToUint32 failed: %v", err)
	}
	if got != 42 {
		t.Errorf("got %d, want 42", got)
	}
}

func TestUint64(t *testing.T) {
	f := Uint64("k", 18446744073709551615)
	if f.type_ != valueUint64 {
		t.Errorf("type_=want valueUint64, got %v", f.type_)
	}
	var got uint64
	if err := f.MarshalToUint64(&got); err != nil {
		t.Fatalf("MarshalToUint64 failed: %v", err)
	}
	if got != 18446744073709551615 {
		t.Errorf("got %d, want 18446744073709551615", got)
	}
}

func TestFloat32(t *testing.T) {
	f := Float32("k", 3.14)
	if f.type_ != valueFloat32 {
		t.Errorf("type_=want valueFloat32, got %v", f.type_)
	}
	var got float32
	if err := f.MarshalToFloat32(&got); err != nil {
		t.Fatalf("MarshalToFloat32 failed: %v", err)
	}
	const eps = 1e-6
	if absFloat(float64(got)-3.14) > eps {
		t.Errorf("got %v, want 3.14", got)
	}
}

func TestFloat64(t *testing.T) {
	f := Float64("k", 2.718)
	if f.type_ != valueFloat64 {
		t.Errorf("type_=want valueFloat64, got %v", f.type_)
	}
	var got float64
	if err := f.MarshalToFloat64(&got); err != nil {
		t.Fatalf("MarshalToFloat64 failed: %v", err)
	}
	if absFloat(got-2.718) > 1e-6 {
		t.Errorf("got %v, want 2.718", got)
	}
}

func TestErrField(t *testing.T) {
	f := Err("k", errTest)
	if f.type_ != valueError {
		t.Errorf("type_=want valueError, got %v", f.type_)
	}
	var got error
	if err := f.MarshalToError(&got); err != nil {
		t.Fatalf("MarshalToError failed: %v", err)
	}
	if got != errTest {
		t.Errorf("got %v, want %v", got, errTest)
	}
}

func TestTimeField(t *testing.T) {
	now := time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC)
	f := Time("k", now)
	if f.type_ != valueTime {
		t.Errorf("type_=want valueTime, got %v", f.type_)
	}
	var got time.Time
	if err := f.MarshalToTime(&got); err != nil {
		t.Fatalf("MarshalToTime failed: %v", err)
	}
	if !got.Equal(now) {
		t.Errorf("got %v, want %v", got, now)
	}
}

func TestClone(t *testing.T) {
	f := Int("key", 99)
	c := f.Clone()
	if c.key != f.key || c.i64 != f.i64 || c.type_ != f.type_ {
		t.Error("clone fields mismatch")
	}
	f.str = "mutated"
	if c.str != "" {
		t.Errorf("clone str mutated: expected empty, got %q", c.str)
	}
	if f.str != "mutated" {
		t.Errorf("original str not mutated: got %q", f.str)
	}
}

// --- MarshalTo negative tests ---

func TestMarshalToTypeMismatch(t *testing.T) {
	strF := String("k", "v")
	var i int
	if err := strF.MarshalToInt(&i); err == nil {
		t.Error("expected error for String -> *int")
	}

	floatF := Float64("k", 1.0)
	if err := floatF.MarshalToInt(&i); err == nil {
		t.Error("expected error for Float64 -> *int")
	}

	errF := Err("k", errTest)
	if err := errF.MarshalToInt(&i); err == nil {
		t.Error("expected error for Err -> *int")
	}

	timeF := Time("k", time.Now())
	if err := timeF.MarshalToInt(&i); err == nil {
		t.Error("expected error for Time -> *int")
	}
}

func TestMarshalToWrongSignedTarget(t *testing.T) {
	intF := Int64("k", 42)
	var u uint
	if err := intF.MarshalToUint(&u); err == nil {
		t.Error("expected error for Int64 -> *uint")
	}
}

func TestMarshalToWrongUnsignedTarget(t *testing.T) {
	uintF := Uint64("k", 42)
	var i int64
	if err := uintF.MarshalToInt64(&i); err == nil {
		t.Error("expected error for Uint64 -> *int64")
	}
}

func TestMarshalStrNilDst(t *testing.T) {
	f := String("k", "v")
	if err := f.MarshalToStr(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestMarshalIntNilDst(t *testing.T) {
	f := Int64("k", 42)
	if err := f.MarshalToInt(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestMarshalInt64NilDst(t *testing.T) {
	f := Int64("k", 42)
	if err := f.MarshalToInt64(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestMarshalToUintNilDst(t *testing.T) {
	f := Uint64("k", 42)
	if err := f.MarshalToUint(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestMarshalToUint64NilDst(t *testing.T) {
	f := Uint64("k", 42)
	if err := f.MarshalToUint64(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestMarshalToFloat32NilDst(t *testing.T) {
	f := Float32("k", 3.14)
	if err := f.MarshalToFloat32(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestMarshalToFloat64NilDst(t *testing.T) {
	f := Float64("k", 3.14)
	if err := f.MarshalToFloat64(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestMarshalToErrorNilDst(t *testing.T) {
	f := Err("k", errTest)
	if err := f.MarshalToError(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestMarshalToTimeNilDst(t *testing.T) {
	f := Time("k", time.Now())
	if err := f.MarshalToTime(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestKeyPreserved(t *testing.T) {
	f := String("my-key", "val")
	if f.key != "my-key" {
		t.Errorf("key expected 'my-key', got %q", f.key)
	}
	c := f.Clone()
	if c.key != "my-key" {
		t.Errorf("clone key expected 'my-key', got %q", c.key)
	}
}

func TestZeroValues(t *testing.T) {
	String("k", "")
	Int("k", 0)
	Int8("k", 0)
	Int16("k", 0)
	Int32("k", 0)
	Int64("k", 0)
	Uint("k", 0)
	Uint8("k", 0)
	Uint16("k", 0)
	Uint32("k", 0)
	Uint64("k", 0)
	Float32("k", 0)
	Float64("k", 0)
	Err("k", nil)
	Time("k", time.Time{})
}

func absFloat(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// --- Tests for Unifields collection ---

func TestNewUnifields(t *testing.T) {
	u := NewUnifields()
	if u == nil {
		t.Fatal("NewUnifields returned nil")
	}
	if len(u.items) != 0 {
		t.Errorf("expected empty items slice, got length %d", len(u.items))
	}
}

func TestUnifieldsAdd(t *testing.T) {
	u := NewUnifields()
	f := Int("counter", 42)
	u.Add(f)
	if len(u.items) != 1 {
		t.Errorf("expected 1 item, got %d", len(u.items))
	}
	// Verify value preserved through clone
	var got int
	if err := u.items[0].MarshalToInt(&got); err != nil {
		t.Fatalf("MarshalToInt failed: %v", err)
	}
	if got != 42 {
		t.Errorf("got %d, want 42", got)
	}
}

func TestUnifieldsAddImmutability(t *testing.T) {
	u := NewUnifields()
	f := String("msg", "hello")
	u.Add(f)
	// Mutate original — collection should be unaffected
	f.key = "modified"
	f.str = "world"
	if u.items[0].key != "msg" || u.items[0].str != "hello" {
		t.Error("collection was mutated by modifying original")
	}
}

func TestUnifieldsAddAll(t *testing.T) {
	u := NewUnifields()
	u.Add(Int("a", 1))
	u.AddStr("b", "two")
	u.AddInt64("c", 3)
	u.AddAll([]Unifield{Uint("d", 4), Err("e", errTest)})
	if len(u.items) != 5 {
		t.Errorf("expected 5 items, got %d", len(u.items))
	}
}

func TestUnifieldsAddAllNilEmpty(t *testing.T) {
	u := NewUnifields()
	initLen := len(u.items)
	u.AddAll(nil)
	if len(u.items) != initLen {
		t.Error("AddAll(nil) should be no-op")
	}
	u.AddAll([]Unifield{})
	if len(u.items) != initLen {
		t.Error("AddAll([]) should be no-op")
	}
}

func TestUnifieldsAddTyped(t *testing.T) {
	u := NewUnifields()
	u.AddStr("k", "v")
	u.AddInt("cnt", 99)
	u.AddInt8("b", 7)
	u.AddInt16("s", 1234)
	u.AddInt32("i", 55555)
	u.AddInt64("l", 999999)
	u.AddUint("u", 42)
	u.AddUint8("ub", 200)
	u.AddUint16("us", 30000)
	u.AddUint32("ui", 200000)
	u.AddUint64("ul", 5000000)
	u.AddFloat32("f32", 1.5)
	u.AddFloat64("f64", 2.718)
	u.AddErr("err", errTest)
	u.AddTime("now", time.Now())
	if len(u.items) != 15 {
		t.Errorf("expected 15 items, got %d", len(u.items))
	}
}

func TestUnifieldsPreservesKeys(t *testing.T) {
	u := NewUnifields()
	f := Int64("my_key", 123)
	u.Add(f)
	if u.items[0].key != "my_key" {
		t.Errorf("expected key 'my_key', got %q", u.items[0].key)
	}
}

func TestUnifieldsCloneSeparation(t *testing.T) {
	u := NewUnifields()
	u.AddStr("original", "value")
	// Create another unifield with same content
	f2 := String("original", "changed")
	u.Add(f2)
	// They should be independent copies
	f3 := u.items[0].Clone()
	f3.str = "independent"
	// Original in collection should be unchanged
	if u.items[0].str != "value" {
		t.Error("collection item was affected by external clone mutation")
	}
}

