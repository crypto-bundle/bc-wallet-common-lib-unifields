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

package unifields

import (
	"errors"
	"testing"

	"github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/ptr"
	valpkg "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/val"
)

var errTest = errors.New("test error value")

func TestNewUnifields(t *testing.T) {
	u := NewUnifields()
	if u == nil {
		t.Fatal("NewUnifields returned nil")
	}
	if u.Len() != 0 {
		t.Errorf("expected empty items slice, got length %d", u.Len())
	}
}

func TestUnifieldsAddVal(t *testing.T) {
	u := NewUnifields()
	f := valpkg.Int(42)
	u.Add(f)
	if u.Len() != 1 {
		t.Errorf("expected 1 item, got %d", u.Len())
	}
	var got int
	if err := f.MarshalToInt(&got); err != nil {
		t.Fatalf("MarshalToInt failed: %v", err)
	}
	if got != 42 {
		t.Errorf("got %d, want 42", got)
	}
}

func TestUnifieldsAddPtr(t *testing.T) {
	u := NewUnifields()
	f := ptr.String("hello world")
	u.Add(f)
	if u.Len() != 1 {
		t.Errorf("expected 1 item, got %d", u.Len())
	}
	var got string
	if err := f.MarshalToStr(&got); err != nil {
		t.Fatalf("MarshalToStr failed: %v", err)
	}
	if got != "hello world" {
		t.Errorf("got %q, want %q", got, "hello world")
	}
}

func TestUnifieldsAddMixedTypes(t *testing.T) {
	u := NewUnifields()
	u.Add(valpkg.Int(1))
	u.Add(ptr.String("two"))
	u.Add(valpkg.Float64(3.0))
	if u.Len() != 3 {
		t.Errorf("expected 3 items, got %d", u.Len())
	}
}

func TestUnifieldsAddAll(t *testing.T) {
	u := NewUnifields()
	u.Add(valpkg.Int(1))
	u.Add(valpkg.String("two"))
	u.Add(valpkg.Int64(3))
	errField := ptr.Err(errTest)
	u.AddAll([]Unifielder{valpkg.Uint(4), errField})
	if u.Len() != 5 {
		t.Errorf("expected 5 items, got %d", u.Len())
	}
}

func TestUnifieldsAddAllNilEmpty(t *testing.T) {
	u := NewUnifields()
	initLen := u.Len()
	u.AddAll(nil)
	if u.Len() != initLen {
		t.Error("AddAll(nil) should be no-op")
	}
	u.AddAll([]Unifielder{})
	if u.Len() != initLen {
		t.Error("AddAll([]) should be no-op")
	}
}

func TestUnifieldsItemsReturnsCopy(t *testing.T) {
	u := NewUnifields()
	items := u.Items()
	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}
}
