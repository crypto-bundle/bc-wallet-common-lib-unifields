# Development Guidelines — bc-wallet-common-lib-unifields

> **Repository:** `github.com/crypto-bundle/bc-wallet-common-lib-unifields`
> **Maintainer:** [@gudron (Alex V Kotelnikov)](https://github.com/gudron)
> **License:** [MIT NON-AI](./LICENSE)

## Package overview

This repo contains **two** packages under `pkg/unifields/`:

| Package | Variant | Storage pattern | Per-field alloc | Struct size |
|---------|---------|----------------|-----------------|-------------|
| `unifields` | Value-based | Flat value fields | Zero | ~56 bytes |
| `unifields/ptr` | Pointer-based | Pointer fields | One alloc per factory | ~80 bytes |

Both share the same 15 supported types and identical API conventions (factory names, MarshalTo methods, error variables). Use `unifields` for zero-allocation hot paths; use `unifields/ptr` when benchmarking trade-offs against heap allocation vs. struct size.

### File layout

```text
pkg/unifields/
├── doc.go                      # Package-level godoc + AI reference
├── unified_field.go            # Unifold type, factories, Clone(), MarshalTo*
├── unified_field_test.go       # Unifold unit tests
├── unified_fields.go           # Unifolds collection (immutable)
├── unified_fields_test.go      # Collection unit tests
└── ptr/
    ├── unified_field_ptr.go    # UnifieldPtr struct, factories, Clone(), MarshalTo*
    ├── unified_field_ptr_test.go   # UnifieldPtr unit tests
    ├── unified_field_ptr_benchmark_test.go # Benchmarks
    └── AGENTS.md               # Expanded agent reference for ptr package
```

### Supported types

| Type | Factory | Read method | Collection adder | Backing field |
|------|---------|-------------|------------------|---------------|
| string | `String(key, val string)` | `MarshalToStr(dst *string)` | `AddStr` | `str` |
| int | `Int(key, val int)` | `MarshalToInt(dst *int)` | `AddInt` | `i64` |
| int8 | `Int8(key, val int8)` | `MarshalToInt8(dst *int8)` | `AddInt8` | `i64` |
| int16 | `Int16(key, val int16)` | `MarshalToInt16(dst *int16)` | `AddInt16` | `i64` |
| int32 | `Int32(key, val int32)` | `MarshalToInt32(dst *int32)` | `AddInt32` | `i64` |
| int64 | `Int64(key, val int64)` | `MarshalToInt64(dst *int64)` | `AddInt64` | `i64` |
| uint | `Uint(key, val uint)` | `MarshalToUint(dst *uint)` | `AddUint` | `u64` |
| uint8 | `Uint8(key, val uint8)` | `MarshalToUint8(dst *uint8)` | `AddUint8` | `u64` |
| uint16 | `Uint16(key, val uint16)` | `MarshalToUint16(dst *uint16)` | `AddUint16` | `u64` |
| uint32 | `Uint32(key, val uint32)` | `MarshalToUint32(dst *uint32)` | `AddUint32` | `u64` |
| uint64 | `Uint64(key, val uint64)` | `MarshalToUint64(dst *uint64)` | `AddUint64` | `u64` |
| float32 | `Float32(key, val float32)` | `MarshalToFloat32(dst *float32)` | `AddFloat32` | `f64` |
| float64 | `Float64(key, val float64)` | `MarshalToFloat64(dst *float64)` | `AddFloat64` | `f64` |
| error | `Err(key, val error)` | `MarshalToError(dst *error)` | `AddErr` | `err` |
| time.Time | `Time(key, val time.Time)` | `MarshalToTime(dst *time.Time)` | `AddTime` | `tm` |

## Design conventions

### Flat storage pattern
Each `Unifold` uses ONE struct with separate fields:
```go
type Unifold struct {
    key   string      // identifier
    err   error       // error value
    tm    time.Time   // time value
    str   string      // string value
    i64   int64       // signed int storage
    u64   uint64      // unsigned int storage
    f64   float64     // float storage
    type_ valueType   // discriminator tag
}
```
Only one data field is active at a time; `type_` indicates which one. This zeroes out unused fields implicitly — no wasted heap per unused type.

### Naming conventions
- **Factory functions**: bare type names → `String`, `Int`, `Float64`, `Err`, `Time`
- **Read methods**: `MarshalTo<Type>` prefix → `MarshalToStr`, `MarshalToInt`, `MarshalToFloat64`
- **Collection adders**: `Add<Type>` prefix → `AddStr`, `AddInt`, `AddAll`
- **Errors**: package-level vars → `ErrDstNil`, `ErrTypeMismatch`, `ErrEmptyUnifold`

### Error handling
- Nil destination pointer → `return ErrDstNil` (first check)
- Type mismatch → `return fmt.Errorf("%w: got %s, want <target>", ErrTypeMismatch, typeName(u.type_))`
- Empty access → `return ErrEmptyUnifold`
- Use `%w` wrapping, never `errors.New()` inline

### Immutability guarantees
- `Clone()` returns value-copy of struct — all fields copied by Go pass-by-value semantics
- Collection `Add()` auto-clones input before appending
- External mutation of returned Unifold does NOT affect collection contents
- `AddAll(s []Unifold)` pre-allocates cloned slice to avoid reallocation

### Float cross-type compatibility
Both `MarshalToFloat32` and `MarshalToFloat64` accept **either** `valueFloat32` OR `valueFloat64` since they share the same `f64` backing field. Narrowing/widening is implicit.

### Pointer-based variant (`unifields/ptr` sub-package)

`UnifieldPtr` stores pointers to typed values instead of flat fields. Only one pointer is non-nil at a time; all others remain nil. This trades per-call heap allocation (one alloc per factory call) for zero cost on unused slots in collections.

#### UnifieldPtr struct layout

```go
type UnifieldPtr struct {
    key   string     // identifier key (first, like Unifold)
    errP  error      // error value (direct, not pointer — error is already an interface)
    i64P  *int64     // signed int pointer
    u64P  *uint64    // unsigned int pointer
    f64P  *float64   // float pointer
    tmP   *time.Time // time pointer
    strP  *string    // string pointer
    type_ valueType // discriminator tag
}
```

Field layout totals ~80 bytes. Due to alignment constraints with string header (16B) and error interface (16B), rearranging fields to optimize memory usage would break intentional grouping. A `//nolint:govet` comment suppresses FieldAlignment lint on this struct.

#### Factory function template (ptr)

```go
func TypeName(key string, val Type) UnifieldPtr {
    v := val             // or v := Type(val) if conversion needed
    return UnifieldPtr{key: key, <field>P: &v, type_: value<Type>}
}
```
Return plain struct literal. Heap copy via local variable ensures immutability. For `Err(key, val error)`, passing nil produces an empty UnifieldPtr.

#### MarshalTo method template (ptr)

```go
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

#### Test patterns (ptr)

In addition to positive, type-mismatch, nil-dst, and zero-value tests, ptr tests must cover empty UnifieldPtr access errors (calling any MarshalTo on `{key: "k"}` should return `ErrEmptyUnifield`). Clone separation tests verify that cloned structs have independent `key` and `type_` fields.

## Build & test commands

```bash
make lint                    # golangci-lint run with project rules
make test                    # go test -race ./...
go vet ./...                 # static analysis
go build ./...               # verify compilation
```

## Code style rules

### File headers
Every `.go` file must include the MIT NON-AI license block from `copyright.txt.tpl`. License text must match verbatim.

### Factory function template
```go
func TypeName(key string, val <Type>) Unifold {
    return Unifold{key: key, <field>: val, type_: value<Type>}
}
```
Return plain struct literal. No nil checks — invalid inputs produce zero-value Unifold (documented behavior).

### MarshalTo method template
```go
func (u Unifold) MarshalTo<Type>(dst *<Type>) error {
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

### Test patterns
Group tests by target file:
- `unified_field_test.go`: factory + Clone + MarshalTo + error cases
- `unifolds_collection_test.go`: New, Add, AddAll, immutability, typed adders, clone separation

Each supported type needs: positive test, type-mismatch negative, nil-dst negative, zero-value test.

### Linter exclusions (version 2 config)
See `.golangci.yml` for full rules. Key exclusions:
- `_test\.go` → cyclop, err113, funlen, gocognit, gocritic, goconst, nlreturn, paralleltest, govet, wsl_v5
- `unified_field\.go` → cyclop, exhaustruct_v5, goheader, gosec, nlreturn, wsl_v5
- `unified_field_ptr\.go` → cyclop, exhaustruct_v5, goheader, gosec, nlreturn, wsl_v5
- `doc\.go` → goheader, gci
- Source lines starting with `* ` → lll (long lines in comments)

## Task tracking

Use `.agents/tasks/NNN_task_name.md` for task descriptions and `.agents/plans/NNN_plan_name.md` for execution plans. Each plan step maps to one git commit. Numbering must be sequential and match between tasks and plans. Update the plan's status table after completing each commit.
