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
	unifolderv2 "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/unifielder"
)

// UnitfieldStack is a LIFO (Last-In-First-Out) stack backed by a slice of [unifolderv2.Unifielder].
// The last element pushed is considered the most recent and is returned first by Pop, PopN, GetTop, and Peek.
// All read/pop operations return the [unifolderv2.Unifielder] interface, enabling polymorphic storage
// of both [val.Unifield] and [ptr.UnifieldPtr]. Empty-stack behavior returns nil — no panics,
// no zero-value structs. Callers must check results before use.
//
// # Usage Example
//
//	stack := val.NewUnitfieldStack()
//	stack.Push(val.String("key", "value"))             // val.Unifield
//	stack.Push(ptr.Int("status", 200))                 // ptr.UnifieldPtr (polymorphic!)
//	top := stack.GetTop()                               // peek without removing
//	popped := stack.Pop()                               // remove and return top
//	all := stack.PopN(3)                                // pop up to 3 elements
//	stack.Clear()                                       // empty the stack
type UnitfieldStack struct {
	items []unifolderv2.Unifielder
}

// NewUnitfieldStack creates a new empty UnitfieldStack with zero-length backing slice.
func NewUnitfieldStack() *UnitfieldStack {
	return &UnitfieldStack{items: make([]unifolderv2.Unifielder, 0)}
}

// Push adds a single Unifielder onto the top of the stack.
// Accepts both val.Unifield and ptr.UnifieldPtr via the [unifolderv2.Unifielder] interface.
func (s *UnitfieldStack) Push(fld unifolderv2.Unifielder) {
	s.items = append(s.items, fld.Clone())
}

// PushFields adds multiple Unifielders onto the top of the stack.
// Each element is cloned before storage to maintain immutability — modifying
// a caller-supplied Unifielder after this call does not affect the stack's contents.
// If fields is nil or empty, this is a no-op.
func (s *UnitfieldStack) PushFields(fields ...unifolderv2.Unifielder) {
	if len(fields) == 0 {
		return
	}
	for _, f := range fields {
		s.items = append(s.items, f.Clone())
	}
}

// Pop removes and returns the top element from the stack.
// Returns nil when the stack is empty. Follows the [UnitfieldList] convention:
// single return value, nil on empty, never panics.
func (s *UnitfieldStack) Pop() unifolderv2.Unifielder { //nolint:ireturn //nolint:ireturn
	n := len(s.items)
	if n == 0 {
		return nil
	}
	idx := n - 1
	result := s.items[idx]
	s.items = s.items[:idx]
	return result
}

// PopField is an alias for [Pop]. Both names perform the same operation.
func (s *UnitfieldStack) PopField() unifolderv2.Unifielder { //nolint:ireturn //nolint:ireturn
	return s.Pop()
}

// PopN removes and returns up to count elements from the top of the stack.
// Returns nil when the stack is empty or count <= 0.
// When count exceeds the stack length, all elements are popped.
// Elements are returned in bottom-to-top order (slice order), not reversed LIFO.
func (s *UnitfieldStack) PopN(count int) []unifolderv2.Unifielder {
	if count <= 0 || len(s.items) == 0 {
		return nil
	}
	if count >= len(s.items) {
		result := s.items
		s.items = s.items[:0]
		return result
	}
	result := s.items[len(s.items)-count:]
	s.items = s.items[:len(s.items)-count]
	return result
}

// Clear removes all elements from the stack, resetting the internal slice to empty.
// Subsequent calls to Len() return 0. No effect when the stack is already empty.
func (s *UnitfieldStack) Clear() {
	s.items = s.items[:0]
}

// Len returns the number of elements currently stored in the stack.
func (s *UnitfieldStack) Len() int {
	return len(s.items)
}

// GetTop returns the top element of the stack without removing it.
// Returns nil when the stack is empty. Follows the [UnitfieldList] convention:
// single return value, nil on empty, never panics.
func (s *UnitfieldStack) GetTop() unifolderv2.Unifielder { //nolint:ireturn //nolint:ireturn
	n := len(s.items)
	if n == 0 {
		return nil
	}
	return s.items[n-1]
}

// Peek is an alias for [GetTop]. Both names peek at the top element without popping.
func (s *UnitfieldStack) Peek() unifolderv2.Unifielder { //nolint:ireturn //nolint:ireturn
	return s.GetTop()
}
