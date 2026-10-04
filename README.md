# bc-wallet-common-lib-unifields

A Go library providing a `zapcore.Field`-like typed value container called **Unifield** — a struct that holds exactly one typed value paired with a string key identifier. Unifields use flat storage internally (a single struct with separate fields for each type and a discriminator tag), which avoids per-field heap allocation and enables cheap copy semantics via `Clone()`.

## Features

- **15 supported types**: string, integers (`int`, `int8`–`int64`, `uint`, `uint8`–`uint64`), floats (`float32`, `float64`), `error`, and `time.Time`
- **Flat storage** — one struct, no heap allocation per field
- **Typed deserialization** via `MarshalTo*` methods with descriptive errors on type mismatch
- **Shallow cloning** through `Clone()` — all stored fields are value-types except `error` (interface)
- **Immutable collection** (`Unifolds`) — all adds store clones, so external mutation cannot affect collection contents
- **Zero allocations** when used as pure value types

## Installation

```bash
go get github.com/crypto-bundle/bc-wallet-common-lib-unifields
```

## Quick Start

### Creating Unifields

Each factory function takes a key string and a typed value:

```go
import (
    "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields"
    "errors"
    "time"
)

// Integer values
f1 := unifields.Int("status_code", 200)
f2 := unifields.Int64("count", 42)
f3 := unifields.Uint64("id", 18446744073709551615)

// String and error
f4 := unifields.String("message", "ok")
f5 := unifields.Err("error", errors.New("connection refused"))

// Time
f6 := unifields.Time("timestamp", time.Now())
```

### Reading values

Each factory has corresponding `MarshalTo*` methods that copy the stored value into a caller-provided pointer:

```go
var val int
if err := f1.MarshalToInt(&val); err != nil {
    // handle type mismatch or nil dst
}
// val == 200
```

Error cases:

```go
var i int
err := f4.MarshalToStr(&i)  // error: type mismatch — expected str, got str → wait...
// Correct example:
err = f4.MarshalToUint(&i)  // ErrTypeMismatch: got str, want uint
```

Passing a nil destination always returns an error:

```go
err := f4.MarshalToStr(nil)  // ErrDstNil
```

### Cloning

`Clone()` returns a shallow copy of the Unifield. Modifying the original does not affect the clone's stored values:

```go
f := unifields.Int("key", 99)
c := f.Clone()

f.Str = "mutated"       // does not affect c.Str (int type doesn't change)
fmt.Println(c.Key)       // "key"
fmt.Println(c.IntVal())  // 99
```

Since all stored fields except `error` (an interface) are pure value-types, `Clone()` fully preserves stored values even after mutating the original.

### Immutable Collection

The `Unifolds` type provides a wrapper around slices of `*Unifield` pointers. All `Add` operations automatically clone their inputs, so modifying a returned Unifield never affects the collection's contents.

```go
cols := unifields.NewUnifolds()
cols.AddInt("user_id", 1234)
cols.AddStr("status", "active")
cols.AddTime("created_at", time.Now())

// AddAll accepts multiple Unifolds at once
cols.AddAll([]unifolds.Unifold{
    unifolds.Int("page", 2),
    unifolds.Bool("enabled", true),
})

// Typed adders are available for every supported type
// AddStr, AddInt, AddInt8..AddInt64, AddUint..AddUint64,
// AddFloat32, AddFloat64, AddErr, AddTime, AddBool
```

## Architecture

```
┌─────────────┐        ┌──────────────────┐
│ Unifold     │───────▶│ Factory Functions│
│ + Clone()   │        │ String(), Int(), │
└─────────────┘        │ Float64(), etc.  │
                       └──────────────────┘

┌─────────────┐        ┌──────────────────┐
│ Unifold     │───────▶│ MarshalTo*       │
│ (typed)     │        │ MarshToStr(),    │
└─────────────┘        │ MarshalTo*, etc. │
                       └──────────────────┘

┌─────────────┐
│ UnifoldCol  │──▶ []*Unifold (internal slice)
│ NewCol(),   │
│ Add(),       │──▶ Auto-clones inputs
│ AddAll()     │
│ AddStr, ...  │──▶ Convenience typed adders
└─────────────┘
```

**Key design decisions:**

| Decision | Rationale |
|----------|-----------|
| Flat storage (one struct per Unifold) | Zero heap allocation per field; cheap copies |
| Discriminator tag (`type_`) instead of `any` + switch | Static type checking at construction time |
| `MarshalTo<T>` naming convention | Matches zapcore.Field deserialization pattern |
| Collection auto-clones | Guarantees immutability — callers can safely mutate return values |

## Package structure

```
pkg/unifields/
├── doc.go                  # Package-level godoc (human-readable)
├── unified_field.go        # Unifold type, factory functions, Clone(), MarshalTo* methods
├── unified_field_test.go   # Tests for Unifold type and factory/MarshalTo functionality
├── unifold_collection.go   # UnifoldCol collection type and methods
└── unifold_collection_test.go  # Tests for UnifoldCol collection behavior
```

## ptr package

The `ptr` sub-package (`github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/ptr`)
provides a **pointer-based variant** of Unifield called `UnifieldPtr`. Instead of storing values
directly in flat struct fields, it stores *pointers* to typed heap-allocated values. Only one pointer
is non-nil at a time, determined by the `type_` discriminator tag.

### Key differences from `pkg/unifields`

| Aspect | `unifields.Unifield` | `ptr.UnifieldPtr` |
|--------|---------------------|--------------------|
| Storage | Flat value fields | Pointer fields |
| Per-field alloc | Zero | One alloc per factory call |
| Struct size | ~56 bytes | ~80 bytes |
| Nil-slot cost | Non-zero (field still exists) | Zero (nil pointer) |
| Clone | Full value copy | Shallow (same underlying refs) |

Use `ptr` to benchmark trade-offs between allocation patterns against the value-based baseline.

### Quick-start example

```go
import (
    "fmt"
    "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/ptr"
)

// Create — each factory allocates one heap copy of the value
p := ptr.Int("user_id", 1234)
s := ptr.String("status", "active")

// Read — MarshalTo dereferences the stored pointer into dst
var val int
if err := p.MarshalToInt(&val); err != nil {
    // handle mismatch or empty
}
fmt.Println(val) // 1234

// Each factory allocates a heap copy, so mutating the original is safe:
original := 999
p2 := ptr.Int("key", original)
original = 0
var v int
p2.MarshalToInt(&v) // v == 999 (not 0)
```

See [AGENTS.md](./pkg/unifields/ptr/AGENTS.md) for detailed API reference, design conventions, and extension points.

## License

This project is licensed under the [MIT NON-AI License](./LICENSE).
