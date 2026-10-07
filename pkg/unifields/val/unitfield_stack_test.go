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
	"testing"
	"time"

	unifolderv2 "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/unifielder"
)

// TestNewUnitfieldStack verifies that the constructor returns a non-nil pointer with Len()==0.
func TestNewUnitfieldStack(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	stack := NewUnitfieldStack()
	if stack == nil {
		t.Fatal("expected non-nil stack")
	}
	if stack.Len() != 0 {
		t.Errorf("expected Len() == 0, got %d", stack.Len())
	}
}

// TestPush verifies single and multiple push operations with LIFO ordering.
func TestPush(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	stack := NewUnitfieldStack()

	fld := String("a", "alpha")
	stack.Push(fld)
	if stack.Len() != 1 {
		t.Errorf("expected Len() == 1, got %d", stack.Len())
	}
	if got := stack.GetTop(); got == nil {
		t.Error("expected GetTop to return non-nil after Push")
	} else if u := got.(Unifield); u.key != "a" { //nolint:forcetypeassert
		t.Errorf("expected top key == 'a', got %q", u.key)
	}

	stack.Push(Int("b", 42))
	if stack.Len() != 2 {
		t.Errorf("expected Len() == 2, got %d", stack.Len())
	}
	if got := stack.GetTop(); got == nil {
		t.Error("expected GetTop to return non-nil after double Push")
	} else if u := got.(Unifield); u.key != "b" { //nolint:forcetypeassert
		t.Errorf("expected top key == 'b' (LIFO), got %q", u.key)
	}
}

// TestPushFields verifies variadic push preserves order and count.
func TestPushFields(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	stack := NewUnitfieldStack()

	fields := []unifolderv2.Unifielder{
		String("f1", "one"),
		Int("f2", 2),
		Float64("f3", 3.0),
	}
	stack.PushFields(fields...)
	if stack.Len() != 3 {
		t.Fatalf("expected Len() == 3, got %d", stack.Len())
	}
	if got := stack.GetTop(); got == nil {
		t.Error("expected GetTop to return non-nil")
	} else if u := got.(Unifield); u.key != "f3" { //nolint:forcetypeassert
		t.Errorf("expected top == 'f3' (last pushed), got %q", u.key)
	}
}

// TestPushFieldsNilSlice ensures nil/empty input is a no-op.
func TestPushFieldsNilSlice(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	stack := NewUnitfieldStack()
	stack.PushFields() // variadic nil
	if stack.Len() != 0 {
		t.Errorf("expected Len() == 0 after nil PushFields, got %d", stack.Len())
	}

	var fields []unifolderv2.Unifielder = nil
	stack.PushFields(fields...)
	if stack.Len() != 0 {
		t.Errorf("expected Len() == 0 after nil slice PushFields, got %d", stack.Len())
	}
}

// TestPop verifies positive path, empty stack behavior, and LIFO guarantee.
func TestPop(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	stack := NewUnitfieldStack()

	// Empty stack returns nil
	empty := stack.Pop()
	if empty != nil {
		t.Errorf("expected nil on empty pop, got %+v", empty)
	}
	if stack.Len() != 0 {
		t.Errorf("expected Len() == 0 after popping empty stack, got %d", stack.Len())
	}

	stack.Push(String("x", "first"))
	stack.Push(Int("y", 1))

	popped := stack.Pop()
	if popped == nil {
		t.Fatal("expected non-nil Pop result")
	}
	if u := popped.(Unifield); u.key != "y" { //nolint:forcetypeassert
		t.Errorf("expected top popped key == 'y' (LIFO), got %q", u.key)
	}
	if stack.Len() != 1 {
		t.Errorf("expected Len() == 1 after one pop, got %d", stack.Len())
	}

	popped2 := stack.Pop()
	if popped2 == nil {
		t.Fatal("expected non-nil second Pop result")
	}
	if u := popped2.(Unifield); u.key != "x" { //nolint:forcetypeassert
		t.Errorf("expected second pop key == 'x', got %q", u.key)
	}
	if stack.Len() != 0 {
		t.Errorf("expected Len() == 0 after two pops, got %d", stack.Len())
	}
}

// TestPopField verifies PopField behaves identically to Pop.
func TestPopField(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	stack := NewUnitfieldStack()
	stack.Push(String("a", "val"))

	result := stack.PopField()
	if result == nil {
		t.Fatal("expected non-nil PopField result")
	}
	if u := result.(Unifield); u.key != "a" { //nolint:forcetypeassert
		t.Errorf("expected key == 'a' from PopField, got %q", u.key)
	}
	if stack.Len() != 0 {
		t.Errorf("expected Len() == 0 after PopField, got %d", stack.Len())
	}

	// Empty stack returns nil
	empty := stack.PopField()
	if empty != nil {
		t.Errorf("expected nil on empty PopField, got %+v", empty)
	}
}

// TestPopN verifies batch pop with various edge cases.
func TestPopN(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	stack := NewUnitfieldStack()

	// Empty stack -> nil
	if got := stack.PopN(2); got != nil {
		t.Errorf("expected nil for empty stack PopN, got %+v", got)
	}

	// n <= 0 -> nil
	if got := stack.PopN(-1); got != nil {
		t.Errorf("expected nil for n== -1, got %+v", got)
	}
	if got := stack.PopN(0); got != nil {
		t.Errorf("expected nil for n== 0, got %+v", got)
	}

	stack.Push(String("a", "1"))
	stack.Push(Int("b", 2))
	stack.Push(Float64("c", 3.0))

	// Pop 2 from 3 elements — returns [b,c] in original array order
	popped := stack.PopN(2)
	if popped == nil {
		t.Fatal("expected non-nil PopN result")
	}
	if len(popped) != 2 {
		t.Fatalf("expected 2 popped, got %d", len(popped))
	}
	if b := popped[0].(Unifield); b.key != "b" { //nolint:forcetypeassert
		t.Errorf("expected [b,c] (slice order), got keys [%q,%q]", b.key, popped[1].(Unifield).key) //nolint:forcetypeassert
	}
	if c := popped[1].(Unifield); c.key != "c" { //nolint:forcetypeassert
		t.Errorf("expected second element key == 'c', got %q", c.key)
	}
	if stack.Len() != 1 {
		t.Errorf("expected remaining Len() == 1, got %d", stack.Len())
	}

	// Remaining element untouched
	top := stack.GetTop()
	if top == nil {
		t.Fatal("expected non-nil GetTop after partial PopN")
	}
	if u := top.(Unifield); u.key != "a" { //nolint:forcetypeassert
		t.Errorf("expected remaining top == 'a', got %q", u.key)
	}

	// Pop more than available -> clamp to all
	all := stack.PopN(5)
	if len(all) != 1 {
		t.Fatalf("expected 1 (clamped) popped when 5 requested from 1-element stack, got %d", len(all))
	}
	if stack.Len() != 0 {
		t.Errorf("expected stack empty after full PopN, got %d", stack.Len())
	}
}

// TestGetTop verifies peeking works without mutation.
func TestGetTop(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	stack := NewUnitfieldStack()

	// Empty stack -> nil
	empty := stack.GetTop()
	if empty != nil {
		t.Errorf("expected nil GetTop on empty stack, got %+v", empty)
	}

	stack.Push(String("a", "val"))
	got := stack.GetTop()
	if got == nil {
		t.Fatal("expected non-nil GetTop after Push")
	}
	if u := got.(Unifield); u.key != "a" { //nolint:forcetypeassert
		t.Errorf("expected GetTop == 'a', got %q", u.key)
	}
	if stack.Len() != 1 {
		t.Errorf("expected GetTop does not mutate, Len() still 1, got %d", stack.Len())
	}

	stack.Push(Int("b", 1))
	top := stack.GetTop()
	if top == nil {
		t.Fatal("expected non-nil GetTop")
	}
	if u := top.(Unifield); u.key != "b" { //nolint:forcetypeassert
		t.Errorf("expected top after double push == 'b', got %q", u.key)
	}
	if stack.Len() != 2 {
		t.Errorf("expected GetTop does not reduce length, got %d", stack.Len())
	}
}

// TestPeek verifies Peek is identical to GetTop.
func TestPeek(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	stack := NewUnitfieldStack()

	// Empty stack
	peeked := stack.Peek()
	if peeked != nil {
		t.Errorf("expected nil Peek on empty stack, got %+v", peeked)
	}

	stack.Push(String("key", "data"))

	first := stack.Peek()
	if first == nil {
		t.Fatal("expected non-nil Peek")
	}
	if u := first.(Unifield); u.key != "key" { //nolint:forcetypeassert
		t.Errorf("expected Peek == 'key', got %q", u.key)
	}
	second := stack.Peek()
	if second == nil {
		t.Fatal("expected non-nil second Peek")
	}
	if u := second.(Unifield); u.key != "key" { //nolint:forcetypeassert
		t.Errorf("expected Peek == 'key' again, got %q", u.key)
	}
	if stack.Len() != 1 {
		t.Errorf("expected Peek does not mutate, Len() == 1, got %d", stack.Len())
	}
}

// TestClear verifies Clear empties the stack completely.
func TestClear(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	stack := NewUnitfieldStack()
	stack.Push(String("a", "1"))
	stack.Push(Int("b", 2))
	stack.Push(Float64("c", 3.0))

	stack.Clear()
	if stack.Len() != 0 {
		t.Errorf("expected Len() == 0 after Clear, got %d", stack.Len())
	}

	// Pop on cleared stack -> nil
	popped := stack.Pop()
	if popped != nil {
		t.Errorf("expected nil Pop after Clear, got %+v", popped)
	}

	// Multiple Clears should be safe no-ops
	stack.Clear()
	stack.Clear()
	if stack.Len() != 0 {
		t.Error("expected Clear to be safe no-op on already-empty stack")
	}
}

// TestLen verifies Len tracks mutations correctly.
func TestLen(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	stack := NewUnitfieldStack()
	if stack.Len() != 0 {
		t.Errorf("expected initial Len() == 0, got %d", stack.Len())
	}

	stack.Push(String("a", "1"))
	if stack.Len() != 1 {
		t.Errorf("expected Len() == 1, got %d", stack.Len())
	}

	stack.PushFields(String("b", "2"), Float64("c", 3.0))
	if stack.Len() != 3 {
		t.Errorf("expected Len() == 3, got %d", stack.Len())
	}

	stack.Pop()
	if stack.Len() != 2 {
		t.Errorf("expected Len() == 2 after Pop, got %d", stack.Len())
	}

	stack.PopN(1)
	if stack.Len() != 1 {
		t.Errorf("expected Len() == 1 after PopN(1), got %d", stack.Len())
	}

	stack.Clear()
	if stack.Len() != 0 {
		t.Errorf("expected Len() == 0 after Clear, got %d", stack.Len())
	}
}

// TestMixedTypes verifies storing different typed Unifields and polymorphic val/ptr mixing.
func TestMixedTypes(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	stack := NewUnitfieldStack()
	stack.Push(String("str", "hello"))
	stack.Push(Int("int", 42))
	stack.Push(Float64("float", 3.14))
	stack.Push(Err("err", ErrTypeMismatch))
	stack.Push(Time("time", time.Now()))

	if stack.Len() != 5 {
		t.Fatalf("expected 5 elements, got %d", stack.Len())
	}

	for range 5 {
		popped := stack.Pop()
		if popped == nil {
			t.Fatalf("pop returned nil instead of stored element")
		}
		key := popped.(Unifield).key //nolint:forcetypeassert
		if key == "" {
			t.Fatalf("pop had empty key")
		}
	}
	if stack.Len() != 0 {
		t.Errorf("expected empty after 5 pops, got %d", stack.Len())
	}
}

// TestImmutability verifies that external mutations don't affect internal state.
func TestImmutability(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	stack := NewUnitfieldStack()
	initial := String("key", "value")
	stack.Push(initial)

	// Mutate the original — should NOT affect stored clone
	initial.key = "changed"

	top := stack.GetTop()
	if top == nil {
		t.Fatal("expected non-nil GetTop")
	}
	if u := top.(Unifield); u.key != "key" { //nolint:forcetypeassert
		t.Errorf("expected immutable key == 'key', got %q", u.key)
	}
}
