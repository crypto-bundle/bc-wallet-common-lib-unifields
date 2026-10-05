# bc-wallet-common-lib-unifields

A Go library providing a `zapcore.Field`-like typed value container called **Unifield** — a struct that holds exactly one typed value paired with a string key identifier. Unifields use flat storage internally (a single struct with separate fields for each type and a discriminator tag), which avoids per-field heap allocation and enables cheap copy semantics via `Clone()`.

## Features

- **Two implementations**: value-based (`val.Unifield`) and pointer-based (`ptr.UnifieldPtr`)
- **15 supported types**: string, integers (`int`, `int8`–`int64`, `uint`, `uint8`–`uint64`), floats (`float32`, `float64`), `error`, and `time.Time`
- **Flat storage** — one struct, no heap allocation per field (val variant)
- **Typed deserialization** via `MarshalTo*` methods with descriptive errors on type mismatch
- **Immutable collection** (`Unifields`) — accepts both `val.Unifield` and `ptr.UnifieldPtr` via [Unifielder](./pkg/unifields/unifielder/doc.go) interface
- **Zero allocations** when used as pure value types

## Installation

```bash
go get github.com/crypto-bundle/bc-wallet-common-lib-unifields
```

## Quick Start

### Creating Values

Each factory function takes a key string and a typed value. Choose between two variants:

```go
import (
    "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields"
    "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/val"
    "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/ptr"
    "errors"
    "time"
)

// Value-based (zero alloc)
f1 := val.String("name", "alice")
f2 := val.Int("status_code", 200)

// Pointer-based (one heap alloc per call)
f3 := ptr.Int("user_id", 1234)
f4 := ptr.Err("error", errors.New("connection refused"))
```

### Reading values

Each factory has corresponding `MarshalTo*` methods that copy the stored value into a caller-provided pointer:

```go
var val int
if err := f2.MarshalToInt(&val); err != nil {
    // handle type mismatch or nil dst
}
// val == 200
```

Error cases:

```go
err = f2.MarshalToUint(&val)  // ErrTypeMismatch: got int, want uint
err := f4.MarshalToStr(nil)   // ErrDstNil
```

### Cloning

`Clone()` returns a shallow copy implementing [Unifielder]. Modifying the original does not affect the clone's stored values (for val types; ptr types share the same underlying pointers since Clone only copies the struct, not heap data).

### Collection

The `Unifields` type provides an immutable collection wrapper that accepts any [Unifielder] interface, allowing mixed-val/ptr collections:

```go
cols := unifields.NewUnifields()
cols.Add(val.String("name", "alice"))    // val.Unifield via Unifielder
cols.Add(ptr.Int("count", 42))           // ptr.UnifieldPtr via Unifielder
cols.AddStr("message", "hello")          // backward-compatible typed adder (defaults to val.)

fmt.Println(cols.Len())                  // 3
items := cols.Items()                    // []Unifielder
```

Available methods:
- `NewUnifields()` — creates an empty collection
- `Add(fld Unifielder)` — adds a cloned Unifielder (supports both val and ptr)
- `AddAll(flds []Unifielder)` — bulk adds cloned Unifielders
- `Len()` — returns item count
- `Items()` — returns read-only copy of stored items
- `AddStr, AddInt, AddInt8..AddInt64, AddUint..AddUint64, AddFloat32, AddFloat64, AddErr, AddTime` — typed adders

## Choosing Between `val` and `ptr`

| Criteria         | `val.Unifield`              | `ptr.UnifieldPtr`          |
|------------------|-----------------------------|----------------------------|
| Allocs per field | 0                           | 1 (heap copy)              |
| Struct size      | ~56 bytes                   | ~80 bytes + heap pointers  |
| Best for         | Zero-cost hot paths         | Pointer identity needed    |

Benchmarks show `val` factories execute at ~0.3ns/ops with zero allocations vs `ptr` at ~0.8ns/ops with one allocation per factory call. For high-throughput scenarios creating thousands of Unifields per second, `val` reduces both CPU cycles and GC pressure.

See `.agents/reports/benchmarks_val_unifield.md` for detailed benchmark results.

## Package Layout

```
pkg/unifields/
├── doc.go                        # Package-level godoc (human-readable API reference)
├── unified_field.go              # Package doc only — types re-exported from val/ & ptr/
├── unified_fields.go             # Unifields collection + Unifielder interface type alias
├── unified_fields_test.go        # Collection tests (polymorphic val/ptr)
├── unifielder/
│   └── doc.go                    # Unifielder interface definition
├─️ val/
│   ├── doc.go                    # Val package godoc
│   ├── unified_field.go          # Unifield struct, factories, Clone(), MarshalTo*
│   ├── unified_field_test.go     # Val unit tests
│   └─️ unified_field_benchmark_test.go # Benchmarks (factory, clone, marshal, collection)
└─️ ptr/
    ├── unified_field_ptr.go      # UnifieldPtr struct, factories, Clone(), MarshalTo*
    ├── unified_field_ptr_test.go # Ptr unit tests
    ├─️ unified_field_ptr_benchmark_test.go # Ptr benchmarks
    └── doc.go                    # Ptr package godoc (detailed API reference)
```

### Import paths

| Package | Path |
|---------|------|
| Parent collection | `github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields` |
| Interface | `github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/unifielder` |
| Value-based | `github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/val` |
| Pointer-based | `github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/ptr` |

## Supported types

Both `val.Unifield` and `ptr.UnifieldPtr` support identical APIs for the same 15 types:

| Type | Factory | Read method | Collection adder |
|------|---------|-------------|------------------|
| string | String(key, val) | MarshalToStr(*string) | AddStr |
| int | Int(key, val) | MarshalToInt(*int) | AddInt |
| int8 | Int8(key, val) | MarshalToInt8(*int8) | AddInt8 |
| int16 | Int16(key, val) | MarshalToInt16(*int16) | AddInt16 |
| int32 | Int32(key, val) | MarshalToInt32(*int32) | AddInt32 |
| int64 | Int64(key, val) | MarshalToInt64(*int64) | AddInt64 |
| uint | Uint(key, val) | MarshalToUint(*uint) | AddUint |
| uint8 | Uint8(key, val) | MarshalToUint8(*uint8) | AddUint8 |
| uint16 | Uint16(key, val) | MarshalToUint16(*uint16) | AddUint16 |
| uint32 | Uint32(key, val) | MarshalToUint32(*uint32) | AddUint32 |
| uint64 | Uint64(key, val) | MarshalToUint64(*uint64) | AddUint64 |
| float32 | Float32(key, val) | MarshalToFloat32(*float32) | AddFloat32 |
| float64 | Float64(key, val) | MarshalToFloat64(*float64) | AddFloat64 |
| error | Err(key, val) | MarshalToError(*error) | AddErr |
| time.Time | Time(key, val) | MarshalToTime(*time.Time) | AddTime |

## Error handling

Every `MarshalTo*` method follows a three-check pattern:

1. Nil destination pointer → returns `ErrDstNil`
2. No value stored (empty/uninitialized) → returns `ErrEmptyUnifield`
3. Type mismatch → returns formatted `%w` error wrapping `ErrTypeMismatch`

Cross-type-compatible methods accept ranges of types sharing the same backing field:
- `MarshalToInt` / `MarshalToInt64` accept any signed integer type
- `MarshalToUint` / `MarshalToUint64` accept any unsigned integer type
- `MarshalToFloat32` / `MarshalToFloat64` accept either (they share `f64` backing field)

## License

This project is licensed under the [MIT NON-AI License](./LICENSE).
