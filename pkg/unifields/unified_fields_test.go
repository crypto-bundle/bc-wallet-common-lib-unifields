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
	"testing"
	"time"
)

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
	f2 := String("original", "changed")
	u.Add(f2)
	f3 := u.items[0].Clone()
	f3.str = "independent"
	if u.items[0].str != "value" {
		t.Error("collection item was affected by external clone mutation")
	}
}
