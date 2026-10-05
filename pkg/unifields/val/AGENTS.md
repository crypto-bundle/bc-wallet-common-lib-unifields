# Development Guidelines — pkg/unifields/val

> **Package:** `github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields/val`
> **Purpose:** Zero-allocation Unifield implementation using flat value fields
> **Interface:** [Unifielder](../unifielder/doc.go) — shared polymorphic contract

## Overview

This package provides `Unifield` — a zero-allocation typed value container using flat storage, modeled after `zapcore.Field`. Each `Unifield` holds exactly one typed value paired with a string key identifier.

Use this package when allocation-free hot paths are critical. For pointer-based semantics (one alloc per factory call, ~80-byte struct), see `pkg/unifields/ptr` instead.

## File layout

```text
pkg/unifields/val/
├── doc.go                        # Package-level godoc
├── unified_field.go              # Unifield type, factories, Clone(), MarshalTo*
├── unified_field_test.go         # Unit tests
└─️ unified_field_benchmark_test.go # Benchmarks (factory, clone, marshal, collection)
```

## Key types

| Type | Purpose | Field |
|------|---------|-------|
| `Unifield` | Single typed value + key | `key`, `err`, `tm`, `str`, `i64`, `u64`, `f64`, `type_` |
| `valueType` | Discriminator enum (iota) | internal, unexported |

### Unifield struct layout

```go
type Unifield struct {
    key   string      // identifier (first for compactness)
    err   error       // error value (interface)
    tm    time.Time   // time value storage
    str   string      // string value storage
    i64   int64       // signed int storage (int, int8..int64)
    u64   uint64      // unsigned int storage (uint, uint8..uint64)
    f64   float64     // float storage (float32, float64)
    type_ valueType   // discriminator tag — only one data field active
}
```

Struct size: ~56 bytes. All fields are zeroed on default/zero-value except `key` which defaults to `""`. Only one of `str`, `i64`, `u64`, `f64`, `err`, `tm` is active based on `type_`.

### Supported types & API table

| Type | Factory | Read method | Backing field |
|------|---------|-------------|---------------|
| string | `String(key, val string)` | `MarshalToStr(dst *string)` | `str` |
| int | `Int(key, val int)` | `MarshalToInt(dst *int)` | `i64` |
| int8 | `Int8(key, val int8)` | `MarshalToInt8(dst *int8)` | `i64` |
| int16 | `Int16(key, val int16)` | `MarshalToInt16(dst *int16)` | `i64` |
| int32 | `Int32(key, val int32)` | `MarshalToInt32(dst *int32)` | `i64` |
| int64 | `Int64(key, val int64)` | `MarshalToInt64(dst *int64)` | `i64` |
| uint | `Uint(key, val uint)` | `MarshalToUint(dst *uint)` | `u64` |
| uint8 | `Uint8(key, val uint8)` | `MarshalToUint8(dst *uint8)` | `u64` |
| uint16 | `Uint16(key, val uint16)` | `MarshalToUint16(dst *uint16)` | `u64` |
| uint32 | `Uint32(key, val uint32)` | `MarshalToUint32(dst *uint32)` | `u64` |
| uint64 | `Uint64(key, val uint64)` | `MarshalToUint64(dst *uint64)` | `u64` |
| float32 | `Float32(key, val float32)` | `MarshalToFloat32(dst *float32)` | `f64` |
| float64 | `Float64(key, val float64)` | `MarshalToFloat64(dst *float64)` | `f64` |
| error | `Err(key, val error)` | `MarshalToError(dst *error)` | `err` |
| time.Time | `Time(key, val time.Time)` | `MarshalToTime(dst *time.Time)` | `tm` |

## Design conventions

### Flat storage pattern

Each `Unifield` uses ONE struct with separate value fields. Only one data field is active at a time; `type_` indicates which one. This zeroes out unused fields implicitly — no wasted heap per unused type.

### Factory function template

```go
func TypeName(key string, val <Type>) Unifield {
    return Unifield{key: key, <field>: val, type_: value<Type>}
}
```

Return plain struct literal. No nil checks — invalid inputs produce zero-value Unifield (documented behavior). Cross-type conversions happen inline (e.g., `int8` → `int64`).

### Clone method template

```go
func (u Unifield) Clone() unifolderv2.Unifielder { //nolint:ireturn
    return u
}
```

`Clone()` returns `Unifielder` interface — intentional design for polymorphic collections in the parent `Unifields` type. The `//nolint:ireturn` comment suppresses golangci-lint since returning an interface is the core abstraction mechanism. Since `Unifield` is a value type with all-copy semantics, cloning produces a fully independent copy.

### MarshalTo method template

```go
func (u Unifield) MarshalTo<Type>(dst *<Type>) error {
    if dst == nil {
        return ErrDstNil          // blank line required after if-block (nlreturn rule)
    }
    if !matches<u.Type>() {
        return fmt.Errorf("%w: got %s, want <type>", ErrTypeMismatch, typeName(u.type_))
    }
    *dst = u.<field>

    return nil                   // blank line before return (nlreturn rule)
}
```

### Error handling

Every `MarshalTo*` method follows this three-check pattern:
1. Nil destination pointer → `return ErrDstNil`
2. Wrong `type_` tag → descriptive `%w` wrapped `ErrTypeMismatch` error
3. Success → writes value into `dst`, returns `nil`

Use `%w` wrapping, never `errors.New()` inline. All three error variables (`ErrDstNil`, `ErrTypeMismatch`, `ErrEmptyUnifield`) are package-level vars defined in `unified_field.go`.

### Float cross-type compatibility

Both `MarshalToFloat32` and `MarshalToFloat64` accept **either** `valueFloat32` OR `valueFloat64` since they share the same `f64` backing field. Narrowing/widening is implicit — callers must be aware that precision loss occurs when reading a `Float64` into a `float32`.

## Testing patterns

Group tests by target file:
- `unified_field_test.go`: factory + Clone + MarshalTo + error cases
- Each supported type needs: positive test, type-mismatch negative, nil-dst negative, zero-value test

For Clone assertions after `Clone()`, note that the method returns `Unifielder` interface — tests within the same package can use value assertion `(Unifield)` to access private fields. Tests in other packages should verify via `MarshalTo*` methods. Add `//nolint:forcetypeassert` when asserting concrete types from `Clone()` return values.

Benchmark tests cover:
- Factory functions (all 15 types)
- Clone operations (representative types)
- MarshalTo methods (all 15 types)
- Collection add overhead (val vs ptr comparison in shared benchmarks)

## Linter rules

Source-file-level exclusions apply:
- cyclop, exhaustruct_v5, goheader, gosec, nlreturn, wsl_v5 → relevant to code complexity
- ireturn → suppressed via `//nolint:ireturn` on Clone() (intentional interface return)

Test-file-level exclusions:
- cyclop, exhaustruct_v5, err113, funlen, gocognit, gocritic, goconst, paralleltest, govet, wsl_v5, forcetypeassert

See `.golangci.yml` for full rule configuration.
