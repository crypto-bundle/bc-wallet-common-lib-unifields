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
	"sync"

	unifolderv2 "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/unifielder"
)

// UnifieldStack is a LIFO (Last-In-First-Out) stack backed by a slice of [unifolderv2.Unifielder].
// The last element pushed is considered the most recent and is returned first by Pop, PopN, GetTop, and Peek.
// All read/pop operations return the [unifolderv2.Unifielder] interface, enabling polymorphic storage
// of both [val.Unifield] and [ptr.UnifieldPtr]. Empty-stack behavior returns nil — no panics,
// no zero-value structs. Callers must check results before use.
//
// Thread safety: UnifieldStack is safe for concurrent access by multiple goroutines.
// Write methods (Push, PushFields, Pop, PopN, Clear) acquire exclusive Lock;
// read-only methods (Len, GetTop, Peek) acquire shared RLock for concurrent readers.
//
// # Usage Example
//
//	stack := val.NewUnifieldStack()
//	stack.Push(val.String("key", "value"))             // val.Unifield
//	stack.Push(ptr.Int("status", 200))                 // ptr.UnifieldPtr (polymorphic!)
//	top := stack.GetTop()                               // peek without removing
//	popped := stack.Pop()                               // remove and return top
//	all := stack.PopN(3)                                // pop up to 3 elements
//	stack.Clear()                                       // empty the stack
type UnifieldStack struct {
	items []unifolderv2.Unifielder
	mu    sync.RWMutex
}

// NewUnifieldStack creates a new empty UnifieldStack with zero-length backing slice.
func NewUnifieldStack() *UnifieldStack {
	return &UnifieldStack{
		mu:    sync.RWMutex{},
		items: make([]unifolderv2.Unifielder, 0),
	}
}

// Push adds a single Unifielder onto the top of the stack.
// Accepts both val.Unifield and ptr.UnifieldPtr via the [unifolderv2.Unifielder] interface.
// Thread-safe: acquires exclusive lock.
func (s *UnifieldStack) Push(fld unifolderv2.Unifielder) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pushNoLock(fld)
}

// pushNoLock adds fld assuming caller holds s.mu.Lock(). Not safe externally.
func (s *UnifieldStack) pushNoLock(fld unifolderv2.Unifielder) {
	s.items = append(s.items, fld.Clone())
}

// PushFields adds multiple Unifielders onto the top of the stack.
// Each element is cloned before storage to maintain immutability — modifying
// a caller-supplied Unifielder after this call does not affect the stack's contents.
// If fields is nil or empty, this is a no-op.
// Thread-safe: acquires exclusive lock.
func (s *UnifieldStack) PushFields(fields ...unifolderv2.Unifielder) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pushFieldsNoLock(fields)
}

// pushFieldsNoLock adds fields assuming caller holds s.mu.Lock(). Not safe externally.
func (s *UnifieldStack) pushFieldsNoLock(fields []unifolderv2.Unifielder) {
	if len(fields) == 0 {
		return
	}

	for i := range fields {
		s.items = append(s.items, fields[i].Clone())
	}
}

// Pop removes and returns the top element from the stack.
// Returns nil when the stack is empty. Follows the [UnitfieldList] convention:
// single return value, nil on empty, never panics.
// Thread-safe: acquires exclusive lock.
func (s *UnifieldStack) Pop() unifolderv2.Unifielder { //nolint:ireturn
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.pop()
}

// pop removes and returns the top element from the stack.
// Caller must hold s.mu.Lock(). Not safe externally.
func (s *UnifieldStack) pop() unifolderv2.Unifielder { //nolint:ireturn
	n := len(s.items)
	if n == 0 {
		return nil
	}

	idx := n - 1
	result := s.items[idx].Clone()          // ① clone BEFORE anything else
	returnValToPool(s.items[idx])           // ② put original into pool
	s.items = s.items[:idx]                 // ③ truncate
	return result                           // ④ caller gets independent copy
}

// PopField is an alias for [Pop]. Both names perform the same operation.
func (s *UnifieldStack) PopField() unifolderv2.Unifielder { //nolint:ireturn
	return s.Pop()
}

// PopN removes and returns up to count elements from the top of the stack.
// Returns nil when the stack is empty or count <= 0.
// When count exceeds the stack length, all elements are popped.
// Elements are returned in bottom-to-top order (slice order), not reversed LIFO.
// Thread-safe: acquires exclusive lock.
func (s *UnifieldStack) PopN(count int) []unifolderv2.Unifielder {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.popNNolockInt(count)
}

// popNNolockInt pops count elements assuming caller holds s.mu.Lock(). Not safe externally.
func (s *UnifieldStack) popNNolockInt(count int) []unifolderv2.Unifielder {
	if count <= 0 || len(s.items) == 0 {
		return nil
	}

	if count >= len(s.items) {
		// Branch A: pop all elements
		result := make([]unifolderv2.Unifielder, len(s.items))
		for i := 0; i < len(s.items); i++ {
			result[i] = s.items[i].Clone()    // clone each element
			returnValToPool(s.items[i])        // pool original
		}
		s.items = s.items[:0]
		return result
	}

	// Branch B: pop partial from top
	startIdx := len(s.items) - count
	result := make([]unifolderv2.Unifielder, count)
	for i := startIdx; i < len(s.items); i++ {
		result[i-startIdx] = s.items[i].Clone()  // clone each element
		returnValToPool(s.items[i])                 // pool original
	}
	s.items = s.items[:startIdx]
	return result
}

// Clear removes all elements from the stack, resetting the internal slice to empty.
// Subsequent calls to Len() return 0. No effect when the stack is already empty.
// Items are returned to the pool.
// Thread-safe: acquires exclusive lock.
func (s *UnifieldStack) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clearNoLock()
}

// clearNoLock clears all items from the stack and returns them to the pool.
// Caller must hold s.mu.Lock(). Not safe externally.
func (s *UnifieldStack) clearNoLock() {
	for _, item := range s.items {
		returnValToPool(item)
	}
	s.items = s.items[:0]
}

// Len returns the number of elements currently stored in the stack.
// Thread-safe: acquires shared read lock, allowing concurrent readers.
func (s *UnifieldStack) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.lenNoLock()
}

// lenNoLock returns the length assuming caller holds s.mu.RLock(). Not safe externally.
func (s *UnifieldStack) lenNoLock() int {
	return len(s.items)
}

// GetTop returns the top element of the stack without removing it.
// Returns nil when the stack is empty. Follows the [UnitfieldList] convention:
// single return value, nil on empty, never panics.
// Thread-safe: acquires shared read lock, allowing concurrent readers.
func (s *UnifieldStack) GetTop() unifolderv2.Unifielder { //nolint:ireturn
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.getTopNoLock()
}

// getTopNoLock returns the top element assuming caller holds s.mu.RLock(). Not safe externally.
func (s *UnifieldStack) getTopNoLock() unifolderv2.Unifielder { //nolint:ireturn
	n := len(s.items)
	if n == 0 {
		return nil
	}

	return s.items[n-1]
}

// Peek is an alias for [GetTop]. Both names peek at the top element without popping.
// Thread-safe: acquires shared read lock.
func (s *UnifieldStack) Peek() unifolderv2.Unifielder { //nolint:ireturn
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.getTopNoLock()
}

// Flush removes and returns all elements from the stack as a slice.
// The stack becomes empty afterward (equivalent to Pop + Clear combined).
// Items returned are clones; originals are returned to the pool.
// Thread-safe: acquires exclusive lock.
func (s *UnifieldStack) Flush() []unifolderv2.Unifielder {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.flush()
}

// flush removes and returns all items, returning originals to the pool.
// Caller must hold s.mu.Lock(). Not safe externally.
func (s *UnifieldStack) flush() []unifolderv2.Unifielder {
	n := len(s.items)
	if n == 0 {
		return nil
	}
	result := make([]unifolderv2.Unifielder, n)
	for i := 0; i < n; i++ {
		result[i] = s.items[i].Clone()
		returnValToPool(s.items[i])
	}
	s.items = s.items[:0]
	return result
}
