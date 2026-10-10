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
	"time"

	valpkg "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/val"
)

// Unifield re-exported from val for convenience. Consumers can use unifields.String() without importing val/.
type Unifield = valpkg.Unifield

// String creates a new Unifield with string value. Zero-allocation path via val.
func String(val string) Unifield { return valpkg.String(val) }

// Int creates a new Unifield with int value. Zero-allocation path via val.
func Int(val int) Unifield { return valpkg.Int(val) }

// Int8 creates a new Unifield with int8 value. Zero-allocation path via val.
func Int8(val int8) Unifield { return valpkg.Int8(val) }

// Int16 creates a new Unifield with int16 value. Zero-allocation path via val.
func Int16(val int16) Unifield { return valpkg.Int16(val) }

// Int32 creates a new Unifield with int32 value. Zero-allocation path via val.
func Int32(val int32) Unifield { return valpkg.Int32(val) }

// Int64 creates a new Unifield with int64 value. Zero-allocation path via val.
func Int64(val int64) Unifield { return valpkg.Int64(val) }

// Uint creates a new Unifield with uint value. Zero-allocation path via val.
func Uint(val uint) Unifield { return valpkg.Uint(val) }

// Uint8 creates a new Unifield with uint8 value. Zero-allocation path via val.
func Uint8(val uint8) Unifield { return valpkg.Uint8(val) }

// Uint16 creates a new Unifield with uint16 value. Zero-allocation path via val.
func Uint16(val uint16) Unifield { return valpkg.Uint16(val) }

// Uint32 creates a new Unifield with uint32 value. Zero-allocation path via val.
func Uint32(val uint32) Unifield { return valpkg.Uint32(val) }

// Uint64 creates a new Unifield with uint64 value. Zero-allocation path via val.
func Uint64(val uint64) Unifield { return valpkg.Uint64(val) }

// Float32 creates a new Unifield with float32 value. Zero-allocation path via val.
func Float32(val float32) Unifield { return valpkg.Float32(val) }

// Float64 creates a new Unifield with float64 value. Zero-allocation path via val.
func Float64(val float64) Unifield { return valpkg.Float64(val) }

// Err creates a new Unifield with error value. Zero-allocation path via val.
func Err(val error) Unifield { return valpkg.Err(val) }

// Time creates a new Unifield with time.Time value. Zero-allocation path via val.
func Time(val time.Time) Unifield { return valpkg.Time(val) }

// UnifieldStack re-exported from val for convenience. Consumers can use unifields.UnifieldStack
// without importing val/.
type UnifieldStack = valpkg.UnifieldStack

// NewUnifieldStack creates a new empty UnifieldStack LIFO stack.
// Wraps [val.NewUnifieldStack] so consumers never need to import the val/ sub-package directly.
func NewUnifieldStack() *UnifieldStack {
	return valpkg.NewUnifieldStack()
}
