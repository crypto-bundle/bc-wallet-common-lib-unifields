# bc-wallet-common-lib-unifields

A Go library providing a `zapcore.Field`-like typed value container called **Unifield** — a struct that holds exactly one typed value paired with a string key identifier. Unifields use flat storage internally (a single struct with separate fields for each type and a discriminator tag), which avoids per-field heap allocation and enables cheap copy semantics via `Clone()`.

## Features

- **Two implementations**: value-based (`val.Unifield`) and pointer-based (`ptr.UnifieldPtr`)
- **Recommended collection**: [`UnitfieldList`](./pkg/unifields/) with list manipulation methods (`Merge`, `GetAfter`, `RemoveAfter`, `RemoveBefore`, `Clear`)
- **Legacy collection**: `Unifolds` (deprecated in favor of `UnitfieldList`)
- **15 supported types**: string, integers (`int`, `int8`–`int64`, `uint`, `uint8`–`uint64`), floats (`float32`, `float64`), `error`, and `time.Time`
- **Flat storage** — one struct, no heap allocation per field (val variant)
- **Typed deserialization** via `MarshalTo*` methods with descriptive errors on type mismatch
- **Immutable collection** — accepts both `val.Unifield` and `ptr.UnifieldPtr` via [Unifielder](./pkg/unifields/unifielder/doc.go) interface
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
    "errors"
    "time"
)

// Zero-allocation path — use unifields package directly, no sub-package import needed
f1 := unifields.String("name", "alice")
f2 := unifields.Int("status_code", 200)
f3 := unifields.Err("error", errors.New("connection refused"))
f4 := unifields.Time("timestamp", time.Now())

// Pointer-based (one heap alloc per call) — import ptr explicitly
// import "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/ptr"
// p := ptr.Int("user_id", 1234)
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
err := f3.MarshalToStr(nil)   // ErrDstNil
```

### Cloning

`Clone()` returns a shallow copy implementing [Unifielder]. Modifying the original does not affect the clone's stored values (for val types; ptr types share the same underlying pointers since Clone only copies the struct, not heap data).

### Collection — Unifolds (legacy)

The `Unifolds` type provides an immutable collection wrapper that accepts any [Unifielder] interface, allowing mixed-val/ptr collections:

```go
cols := unifields.NewUnifolds()
cols.Add(unifields.String("name", "alice"))   // val.Unifield via Unifielder
cols.Add(unifields.Int("count", 42))          // val.Unifield via Unifielder

fmt.Println(cols.Len())                       // 2
items := cols.Items()                         // []Unifielder
```

Available methods:
- `NewUnifolds()` — creates an empty collection
- `Add(fld Unifielder)` — adds a cloned Unifielder (supports both val and ptr)
- `AddAll(flds []Unifielder)` — bulk adds cloned Unifielders
- `Len()` — returns item count
- `Items()` — returns read-only copy of stored items

### Collection — UnitfieldList (recommended)

[`UnitfieldList`](./pkg/unifields/) is the recommended collection with list-manipulation capabilities:

```go
import "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields"

list := unifields.NewUnitfieldList()
list.AddStr("name", "alice")           // shorthand typed adders
list.AddInt("count", 42)

// Merge another list into this one
other := unifields.NewUnitfieldList()
other.AddStr("key", "value")
list.Merge(other)

// Split a list by index
tail := list.GetAfter(0)               // elements from index 0 onward
list.RemoveAfter(0)                    // keep only element at index 0
list.RemoveBefore(2)                   // keeps elements from index 2 onward
list.Clear()                           // drop all elements
```

Available methods:
- `NewUnitfieldList()` — creates an empty collection
- `Add(fld Unifielder)` — adds a cloned Unifielder
- `AddAll(flds []Unifielder)` — bulk adds cloned Unifielders
- `Merge(list UnitfieldList)` — appends all items from another list (each element is cloned)
- `GetAfter(index uint) []Unifielder` — returns copy of items starting from index N (inclusive)
- `RemoveAfter(index uint)` — keeps element at index, drops everything after it
- `RemoveBefore(index uint)` — keeps element at index, drops everything before it
- `Clear()` — resets the list to empty
- `Len()` — returns item count
- `Items()` — returns read-only copy of stored items
- `AddStr`, `AddInt`, `AddInt8`…`AddInt64`, `AddUint`…`AddUint64`, `AddFloat32`, `AddFloat64`, `AddErr`, `AddTime` — typed adders

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
├── unified_field.go              # Type alias + factory wrappers + UnitfieldList bridge
├── unified_fields.go             # Deprecated Unifolds collection + Unifielder interface alias
├── unified_fields_test.go        # Collection tests (polymorphic val/ptr)
├── unifielder/
│   └── doc.go                    # Unifielder interface definition
├─️ val/
│   ├── unified_field.go          # Unifield struct, factories, Clone(), MarshalTo*
│   ├── unified_field_test.go     # Val unit tests
│   ├── unified_field_benchmark_test.go       # Benchmarks (factory, clone, marshal, collection)
│   ├── unified_field_list.go                 # UnitfieldList type + list manipulation methods
│   ├── unified_field_list_test.go            # UnitfieldList unit tests
│   └─️ unified_field_list_benchmark_test.go  # UnitfieldList benchmarks
└─️ ptr/
    ├── unified_field_ptr.go          # UnifieldPtr struct, factories, Clone(), MarshalTo*
    ├── unified_field_ptr_test.go     # Ptr unit tests
    ├─️ unified_field_ptr_benchmark_test.go # Ptr benchmarks
    └── doc.go                        # Ptr package godoc (detailed API reference)
```

### Import paths

| Package | Path |
|---------|------|
| Parent collection + wrappers | `github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields` |
| Interface | `github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/unifielder` |
| Value-based | `github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/val` |
| Pointer-based | `github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/ptr` |

### Wrapper functions (no `val/` import needed)

For everyday usage, all 15 factory functions are available from the parent `unifields` package:

```go
import "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields"

f := unifields.String("key", "value")   // delegates to val.String()
g := unifields.Int("counter", 42)       // delegates to val.Int()
h := unifields.Err("err", someErr)      // delegates to val.Err()
```

A type alias `Unifield = val.Unifield` is defined at the parent level so consumers can use the `Unifield` type name without importing `val/`. This eliminates boilerplate: instead of `import ".../val"` just to call `val.String()`, you write `unifields.String()`.

## Supported types

Both `val.Unifield` and `ptr.UnifieldPtr` support identical APIs for the same 15 types. Additionally, `UnitfieldList` exposes typed adders matching each type.

| Type | Factory (val/ptr) | Read method | Collection adder |
|------|--------------------|-------------|------------------|
| string | String(key, val) | MarshalToStr(*string) | AddStr |
| int | Int(key, val) | MarshalToInt(*int) | AddInt |
| int8 | Int8(key, val) | MarshalToInt8(*int8) | AddInt8 |
| int16 | Int16(key, val) | MarshalToInt16(*int16) | AddInt16 |
| int32 | Int32(key, val) | MarshalToUint16(*uint16) | AddInt16 |
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

## UnitfieldList methods

[`UnitfieldList`](./pkg/unifields/) extends the core collection with list manipulation operations:

| Method | Description |
|--------|-------------|
| `Merge(list UnitfieldList)` | Appends all items from another list (each cloned) |
| `GetAfter(index uint) []Unifielder` | Returns copy of items starting from index N (inclusive) |
| `RemoveAfter(index uint)` | Keeps element at index, drops everything after |
| `RemoveBefore(index uint)` | Keeps element at index, drops everything before |
| `Clear()` | Resets the list to empty |

## Error handling

Every `MarshalTo*` method follows a three-check pattern:

1. Nil destination pointer → returns `ErrDstNil`
2. No value stored (empty/uninitialized) → returns `ErrEmptyUnifield`
3. Type mismatch → returns formatted `%w` error wrapping `ErrTypeMismatch`

Cross-type-compatible methods accept ranges of types sharing the same backing field:
- `MarshalToInt` / `MarshalToInt64` accept any signed integer type
- `MarshalToUint` / `MarshalToUint64` accept any unsigned integer type
- `MarshalToFloat32` / `MarshalToFloat64` accept either (they share `f64` backing field)

## Deprecated — Unifolds

The legacy `Unifolds` type remains in the codebase for backward compatibility but is deprecated. All typed adders have been migrated exclusively to `UnitfieldList`. Migrate your code at convenience:

```go
// Old (still works):
cols := unifields.NewUnifolds()
cols.Add(unifields.String("key", "value"))

// New (recommended):
list := unifields.NewUnitfieldList()
list.AddStr("key", "value")  // typed adder
list.Merge(otherList)        // list manipulation
```

## License

This project is licensed under the [MIT NON-AI License](./LICENSE).
