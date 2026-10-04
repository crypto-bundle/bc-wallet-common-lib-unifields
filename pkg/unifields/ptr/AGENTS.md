# ptr package — UnifieldPtr

## Project overview

Pointer-based variant of Unifield — `UnifieldPtr`. Stores pointers to typed values instead of flat values. Each struct has separate pointer slots (`strP`, `i64P`, `u64P`, `f64P`, `tmP`) and a direct error field (`errP`). Only one pointer is non-nil at a time, determined by the `type_` discriminator tag.

This variant trades increased per-allocation heap overhead (one alloc per factory call) for potential memory savings when stored in large collections (nil pointers are smaller than fat value slots on some platforms). Use it to benchmark trade-offs against the value-based `Unifield` in the parent `unifields` package.

### Supported types

| Type | Factory | Read method | Backing field |
|------|---------|-------------|---------------|
| string | `String(key, val string)` | `MarshalToStr(dst *string)` | `strP` |
| int | `Int(key, val int)` | `MarshalToInt(dst *int)` | `i64P` |
| int8 | `Int8(key, val int8)` | `MarshalToInt8(dst *int8)` | `i64P` |
| int16 | `Int16(key, val int16)` | `MarshalToInt16(dst *int16)` | `i64P` |
| int32 | `Int32(key, val int32)` | `MarshalToInt32(dst *int32)` | `i64P` |
| int64 | `Int64(key, val int64)` | `MarshalToInt64(dst *int64)` | `i64P` |
| uint | `Uint(key, val uint)` | `MarshalToUint(dst *uint)` | `u64P` |
| uint8 | `Uint8(key, val uint8)` | `MarshalToUint8(dst *uint8)` | `u64P` |
| uint16 | `Uint16(key, val uint16)` | `MarshalToUint16(dst *uint16)` | `u64P` |
| uint32 | `Uint32(key, val uint32)` | `MarshalToUint32(dst *uint32)` | `u64P` |
| uint64 | `Uint64(key, val uint64)` | `MarshalToUint64(dst *uint64)` | `u64P` |
| float32 | `Float32(key, val float32)` | `MarshalToFloat32(dst *float32)` | `f64P` |
| float64 | `Float64(key, val float64)` | `MarshalToFloat64(dst *float64)` | `f64P` |
| error | `Err(key, val error)` | `MarshalToError(dst *error)` | `errP` |
| time.Time | `Time(key, val time.Time)` | `MarshalToTime(dst *time.Time)` | `tmP` |

## Design conventions

### Pointer storage pattern
Each `UnifieldPtr` uses pointer fields with ONE non-nil at a time:
```go
type UnifieldPtr struct {
    key   string      // identifier
    errP  error       // error value (direct, not pointer — error is already an interface)
    i64P  *int64      // signed int pointer
    u64P  *uint64     // unsigned int pointer
    f64P  *float64    // float pointer
    tmP   *time.Time  // time pointer
    strP  *string     // string pointer
    type_ valueType   // discriminator tag
}
```
The `errP` field stores error directly because `error` is already an interface — storing `*error` would require double-indirection (`*interface{}`) which adds no value. All other pointer fields start nil; only one is non-nil after construction.

### Naming conventions
- **Factory functions**: bare type names → `String`, `Int`, `Float64`, `Err`, `Time`
- **Read methods**: `MarshalTo<Type>` prefix → `MarshalToStr`, `MarshalToInt`, `MarshalToFloat64`
- **Errors**: package-level vars → `ErrDstNil`, `ErrTypeMismatch`, `ErrEmptyUnifield`

### Error handling
- Nil destination pointer → `return ErrDstNil` (first check)
- No value stored (all pointers nil) → `return ErrEmptyUnifield` (second check)
- Type mismatch → `return fmt.Errorf("%w: got %s, want <target>", ErrTypeMismatch, typeName(u.type_))`
- Use `%w` wrapping, never `errors.New()` inline
- Cross-type-compatible methods (MarshalToInt/MarshalToInt64 accept any signed int; MarshalToUint/MarshalToUint64 accept any unsigned int; MarshalToFloat32/MarshalToFloat64 accept either) use range predicates rather than exact equality

### Immutability guarantees
- Factory functions allocate a heap copy of the input value before storing the pointer
- `Clone()` returns a shallow copy — both original and clone share references to the same underlying heap values
- Since no mutable APIs exist, shared pointers are safe in practice
- External mutation of original values does NOT affect stored UnifieldPtr values

### Float cross-type compatibility
Both `MarshalToFloat32` and `MarshalToFloat64` accept **either** `valueFloat32` OR `valueFloat64` since they share the same `f64P` backing field. Narrowing/widening is implicit.

### Struct layout
Due to alignment constraints with string header (16B) and error interface (16B), the struct totals 80 bytes. The `//nolint:govet` comment suppresses FieldAlignment lint because rearranging fields to optimize alignment would break the intentional grouping.

## Build & test commands

```bash
make lint                    # golangci-lint run with project rules (covers ptr/)
make test                    # go test -race ./...
go vet ./...                 # static analysis
go build ./...               # verify compilation
```

## Code style rules

### File headers
Every `.go` file must include the MIT NON-AI license block from `copyright.txt.tpl`. License text must match verbatim.

### Factory function template
```go
// Type creates a new UnifieldPtr with key and typed value.
//
// It allocates a heap copy of val so that mutating the original value after this
// call does not affect the stored value. Only one pointer field will be non-nil
// in the returned UnifieldPtr, determined by the type_.
func Type(key string, val Type) UnifieldPtr {
    v := val            // or v := Type(val) if conversion needed
    return UnifieldPtr{key: key, <field>P: &v, type_: value<Type>}
}
```
Return plain struct literal. Heap copy via local variable ensures immutability. For `Err(key, val error)`, pass-through nil produces an empty UnifieldPtr.

### MarshalTo method template
```go
// MarshalTo<Type> copies the stored <type> value into dst.
//
// It returns ErrDstNil if dst is nil, ErrEmptyUnifield if no value is stored,
// or ErrTypeMismatch if the stored type is not <type>.
func (u UnifieldPtr) MarshalTo<DstType>(dst *DstType) error {
    if dst == nil {                          // blank line required after if-block (nlreturn rule)
        return ErrDstNil
    }
    if u.<field>P == nil {                   // blank line required after if-block (nlreturn rule)
        return ErrEmptyUnifield
    }
    if !is<Predicate>(u.type_) {             // or exact match: if u.type_ != value<Type>
        return fmt.Errorf("%w: got %s, want <type>", ErrTypeMismatch, typeName(u.type_))
    }
    *dst = <conversion>(*u.<field>P)

    return nil                               // blank line before return (nlreturn rule)
}
```

### Test patterns
Tests should cover: factory creation + verification, MarshalTo positive cases, nil-dst negative, empty/unset negative, type-mismatch negative, zero-value tests, Clone separation, and cross-type compatibility for signed/unsigned/float ranges.

### Linter exclusions
Key rules:
- `//nolint:govet` on `UnifieldPtr struct` — alignment constrained by string+error+5*pointer types
- `_test\.go` paths → cyclop, exhaustruct_v5, funlen, gocognit, govet, wsl_v5 (matching parent config)

## Agent task tracking

Use `.agents/tasks/NNN_task_name.md` for task descriptions and `.agents/plans/NNN_plan_name.md` for execution plans. Numbering matches repo-level task tracking system. See `../../../AGENTS.md` at repo root for expanded style rules, CSG references, and full task-tracking conventions.
