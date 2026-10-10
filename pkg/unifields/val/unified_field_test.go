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
	"errors"
	"testing"
	"time"
)

var errTest = errors.New("test error value")

func TestString(t *testing.T) {
	f := String("val")
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
	fz := String("")
	var sz string
	if err := fz.MarshalToStr(&sz); err != nil {
		t.Fatalf("MarshalToStr(\"\") failed: %v", err)
	}
	if sz != "" {
		t.Errorf("expected empty, got %q", sz)
	}
}

func TestIntBasic(t *testing.T) {
	f := Int(42)
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
	f := Int8(42)
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
	f := Int16(12345)
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
	f := Int32(123456789)
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
	f := Int64(9223372036854775807)
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
	f := Uint(42)
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
	f := Uint8(42)
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
	f := Uint16(42)
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
	f := Uint32(42)
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
	f := Uint64(18446744073709551615)
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
	f := Float32(3.14)
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
	f := Float64(2.718)
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
	f := Err(errTest)
	if f.type_ != valueError {
		t.Errorf("type_=want valueError, got %v", f.type_)
	}
	var got error
	if err := f.MarshalToError(&got); err != nil {
		t.Fatalf("MarshalToError failed: %v", err)
	}
	if !errors.Is(got, errTest) {
		t.Errorf("got %v, want %v", got, errTest)
	}
}

func TestTimeField(t *testing.T) {
	now := time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC)
	f := Time(now)
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
	f := Int(99)
	c := f.Clone()
	clonedVal := c.(Unifield)
	if clonedVal.i64 != f.i64 || clonedVal.type_ != f.type_ {
		t.Error("clone fields mismatch")
	}
}

// --- MarshalTo negative tests ---

func TestMarshalToTypeMismatch(t *testing.T) {
	strF := String("v")
	var i int
	if err := strF.MarshalToInt(&i); err == nil {
		t.Error("expected error for String -> *int")
	}

	floatF := Float64(1.0)
	if err := floatF.MarshalToInt(&i); err == nil {
		t.Error("expected error for Float64 -> *int")
	}

	errF := Err(errTest)
	if err := errF.MarshalToInt(&i); err == nil {
		t.Error("expected error for Err -> *int")
	}

	timeF := Time(time.Now())
	if err := timeF.MarshalToInt(&i); err == nil {
		t.Error("expected error for Time -> *int")
	}
}

func TestMarshalToWrongSignedTarget(t *testing.T) {
	intF := Int64(42)
	var u uint
	if err := intF.MarshalToUint(&u); err == nil {
		t.Error("expected error for Int64 -> *uint")
	}
}

func TestMarshalToWrongUnsignedTarget(t *testing.T) {
	uintF := Uint64(42)
	var i int64
	if err := uintF.MarshalToInt64(&i); err == nil {
		t.Error("expected error for Uint64 -> *int64")
	}
}

func TestMarshalStrNilDst(t *testing.T) {
	f := String("v")
	if err := f.MarshalToStr(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestMarshalIntNilDst(t *testing.T) {
	f := Int64(42)
	if err := f.MarshalToInt(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestMarshalInt64NilDst(t *testing.T) {
	f := Int64(42)
	if err := f.MarshalToInt64(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestMarshalToUintNilDst(t *testing.T) {
	f := Uint64(42)
	if err := f.MarshalToUint(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestMarshalToUint64NilDst(t *testing.T) {
	f := Uint64(42)
	if err := f.MarshalToUint64(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestMarshalToFloat32NilDst(t *testing.T) {
	f := Float32(3.14)
	if err := f.MarshalToFloat32(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestMarshalToFloat64NilDst(t *testing.T) {
	f := Float64(3.14)
	if err := f.MarshalToFloat64(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestMarshalToErrorNilDst(t *testing.T) {
	f := Err(errTest)
	if err := f.MarshalToError(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestMarshalToTimeNilDst(t *testing.T) {
	f := Time(time.Now())
	if err := f.MarshalToTime(nil); err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestZeroValues(t *testing.T) {
	var _ Unifield
	_ = String("")
	_ = Int(0)
	_ = Int8(0)
	_ = Int16(0)
	_ = Int32(0)
	_ = Int64(0)
	_ = Uint(0)
	_ = Uint8(0)
	_ = Uint16(0)
	_ = Uint32(0)
	_ = Uint64(0)
	_ = Float32(0)
	_ = Float64(0)
	_ = Err(nil)
	_ = Time(time.Time{})
}

func absFloat(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
