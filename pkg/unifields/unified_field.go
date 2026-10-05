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

// Package unifields provides typed value containers and an immutable collection.
//
// Factory functions create Unifield instances. Each function takes a key string and a typed value:
//
//	import "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields"
//
//	// Integer values
//	f1 := unifields.Int("status_code", 200)
//	f2 := unifields.Int64("count", 42)
//	f3 := unifields.Uint64("id", 18446744073709551615)
//
//	// String and error
//	f4 := unifields.String("message", "ok")
//	f5 := unifields.Err("error", someError)
//
//	// Time
//	f6 := unifields.Time("timestamp", time.Now())
//
// # Reading values
//
// Each factory function has corresponding typed methods prefixed with MarshalTo. These methods
// copy the stored value into a caller-provided pointer:
//
//	var val int
//	if err := f1.MarshalToInt(&val); err != nil {
//	    // handle type mismatch
//	}
//
// If the Unifield holds a different type than expected, MarshalTo returns a descriptive error.
// Passing a nil destination pointer also returns an error.
//
// # Cloning
//
// Clone returns a shallow copy of the Unifield. Since all stored fields are value-types
// (except error which is an interface), modifying the original after cloning does not affect
// the clone's stored values.
//
// # Unifields collection
//
// The Unifields type provides an immutable collection wrapper around []Unifielder. All Add operations
// store clones of input values, ensuring external mutation cannot affect the collection's contents.
//
//	import "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields"
//
//	cols := unifields.NewUnifields()
//	cols.AddStr("user_id", "1234")
//	cols.AddInt("page", 2)
//	cols.AddTime("created_at", time.Now())
//
// Available methods:
//   - NewUnifields() — creates an empty collection
//   - Add(fld Unifield) — adds a cloned Unifield
//   - AddAll(flds []Unifield) — bulk adds cloned Unifields
//   - AddStr, AddInt, AddInt8..AddInt64, AddUint..AddUint64, AddFloat32, AddFloat64, AddErr, AddTime — typed adders
package unifields
