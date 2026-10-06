# Development Guidelines — bc-wallet-common-lib-unifields

> **Repository:** `github.com/crypto-bundle/bc-wallet-common-lib-unifields`
> **Maintainer:** [@gudron (Alex V Kotelnikov)](https://github.com/gudron)
> **License:** [MIT NON-AI](./LICENSE)

## Package overview

This repo contains **three** packages under `pkg/unifields/`:

| Package | Variant | Storage pattern | Per-field alloc | Struct size | Interface |
|---------|---------|----------------|-----------------|-------------|-----------|
| `unifields` | Collection + wrappers | Slice of `Unifielder` | N/A | ~16 bytes (slice header) | Holds `Unifielder` items |
| `val` | Value-based | Flat value fields | Zero | ~56 bytes | Implements `Unifielder` |
| `ptr` | Pointer-based | Pointer fields | One alloc per factory | ~80 bytes | Implements `Unifielder` |

**Quick access:** Consumers can use factory functions directly from `pkg/unifields/` without importing sub-packages: `unifields.String(...)`, `unifields.Int(...)`, etc. These wrap `val.` internally for zero-allocation usage. Both `val.Unifield` and `ptr.UnifieldPtr` share the same 15 supported types and identical API conventions. Use `val` for zero-allocation hot paths; use `ptr` when pointer identity or deferred binding is needed.

For collections, prefer [`UnitfieldList`](./pkg/unifields/) — the recommended collection with list-manipulation methods (`Merge`, `GetAfter`, `RemoveAfter`, `RemoveBefore`, `Clear`). The legacy `Unifolds` type remains for backward compatibility but is deprecated.

### File layout

```text
pkg/unifields/
├── unified_field.go                # Type alias (Unifield = val.Unifield) + 15 factory wrappers + UnitfieldList bridge
├── unified_fields.go             # Deprecated Unifolds collection + Unifielder interface alias
├── unified_fields_test.go        # Collection tests (polymorphic val/ptr)
└── unifielder/
    └── doc.go                    # Unifielder interface definition
└─️ val/
    ├── unified_field.go          # Unifield struct, factories, Clone(), MarshalTo*
    ├── unified_field_test.go     # Val unit tests
    ├─️ unified_field_benchmark_test.go # Benchmarks (factory, clone, marshal, collection)
    ├── unified_field_list.go     # UnitfieldList type + list manipulation methods
    ├── unified_field_list_test.go         # UnitfieldList unit tests
    └─️ unified_field_list_benchmark_test.go # UnitfieldList benchmarks
└─️ ptr/
    ├── unified_field_ptr.go      # UnifieldPtr struct, factories, Clone(), MarshalTo*
    ├── unified_field_ptr_test.go # Ptr unit tests
    ├─️ unified_field_ptr_benchmark_test.go # Ptr benchmarks
    └── doc.go                    # Ptr package godoc (detailed API reference)
```

### Supported types

| Type | Factory | Read method | Collection adder | Backing field (val) | Backing field (ptr) |
|------|---------|-------------|------------------|--------------------|--------------------|
| string | `String(key, val string)` | `MarshalToStr(dst *string)` | `AddStr` | `str string` | `strP *string` |
| int | `Int(key, val int)` | `MarshalToInt(dst *int)` | `AddInt` | `i64 int64` | `i64P *int64` |
| int8 | `Int8(key, val int8)` | `MarshalToInt8(dst *int8)` | `AddInt8` | `i64 int64` | `i64P *int64` |
| int16 | `Int16(key, val int16)` | `MarshalToInt16(dst *int16)` | `AddInt16` | `i64 int64` | `i64P *int64` |
| int32 | `Int32(key, val int32)` | `MarshalToInt32(dst *int32)` | `AddInt32` | `i64 int64` | `i64P *int64` |
| int64 | `Int64(key, val int64)` | `MarshalToInt64(dst *int64)` | `AddInt64` | `i64 int64` | `i64P *int64` |
| uint | `Uint(key, val uint)` | `MarshalToUint(dst *uint)` | `AddUint` | `u64 uint64` | `u64P *uint64` |
| uint8 | `Uint8(key, val uint8)` | `MarshalToUint8(dst *uint8)` | `AddUint8` | `u64 uint64` | `u64P *uint64` |
| uint16 | `Uint16(key, val uint16)` | `MarshalToUint16(dst *uint16)` | `AddUint16` | `u64 uint64` | `u64P *uint64` |
| uint32 | `Uint32(key, val uint32)` | `MarshalToUint32(dst *uint32)` | `AddUint32` | `u64 uint64` | `u64P *uint64` |
| uint64 | `Uint64(key, val uint64)` | `MarshalToUint64(dst *uint64)` | `AddUint64` | `u64 uint64` | `u64P *uint64` |
| float32 | `Float32(key, val float32)` | `MarshalToFloat32(dst *float32)` | `AddFloat32` | `f64 float64` | `f64P *float64` |
| float64 | `Float64(key, val float64)` | `MarshalToFloat64(dst *float64)` | `AddFloat64` | `f64 float64` | `f64P *float64` |
| error | `Err(key, val error)` | `MarshalToError(dst *error)` | `AddErr` | `err error` | `errP error` |
| time.Time | `Time(key, val time.Time)` | `MarshalToTime(dst *time.Time)` | `AddTime` | `tm time.Time` | `tmP *time.Time` |

## Design conventions

### Flat storage pattern (val)
Each `val.Unifield` uses ONE struct with separate value fields:
```go
type Unifield struct {
    key   string      // identifier
    err   error       // error value (interface)
    tm    time.Time   // time value
    str   string      // string value
    i64   int64       // signed int storage
    u64   uint64      // unsigned int storage
    f64   float64     // float storage
    type_ valueType   // discriminator tag
}
```
Only one data field is active at a time; `type_` indicates which one. This zeroes out unused fields implicitly — no wasted heap per unused type.

### Pointer storage pattern (ptr)
Each `ptr.UnifieldPtr` stores pointers to typed values:
```go
type UnifieldPtr struct { //nolint:govet // alignment constrained by string+error+5*pointer types
    key   string     // identifier key
    errP  error      // error value stored directly (not as pointer)
    i64P  *int64     // pointer to signed int value (nil if not set)
    u64P  *uint64    // pointer to unsigned int value (nil if not set)
    f64P  *float64   // pointer to float value (nil if not set)
    tmP   *time.Time // pointer to time.Time value (nil if not set)
    strP  *string    // pointer to string value (nil if not set)
    type_ valueType   // discriminator tag
}
```
Field layout totals ~80 bytes. Due to alignment constraints with string header (16B) and error interface (16B), rearranging fields would break intentional grouping.

### Interface abstraction (`Unifielder`)
The shared [Unifielder](./unifielder/doc.go) interface enables polymorphic collection storage:
```go
type Unifielder interface {
    Clone() Unifielder
}
```
Both `val.Unifield` and `ptr.UnifieldPtr` implement this interface. The return type `Unifielder` allows heterogeneous collections (`[]Unifielder`) without type parameters. See `.agents/plans/003_unifield_val_refactor_plan.md` for architectural rationale.

### UnitfieldList design conventions

[`UnitfieldList`](./pkg/unifields/val/unified_field_list.go) is the recommended collection type with list-manipulation capabilities. It wraps `[]unifolderv2.Unifielder` and provides structured list operations:

**Methods:**
- `Merge(list UnitfieldList)` — appends each item from `list.items` via `Clone()`. Uses **pointer-receiver** with defensive copy to prevent aliasing through shared backing array. Keys may not be unique across merged result.
- `GetAfter(index uint) []Unifielder` — returns copy of items from index N (inclusive); graceful no-op when `index >= len`.
- `RemoveAfter(index uint)` — keeps element at index N, drops N+1..end; graceful no-op when `index >= len`.
- `RemoveBefore(index uint)` — keeps element at index N, drops 0..N-1; graceful no-op when `index <= 0` or `index >= len`.
- `Clear()` — resets items slice to empty (zero-capacity).

**Graceful bounds handling:** All index-based methods check bounds and return gracefully rather than panicking on out-of-range indices (`index >= len(items)`).

### Naming conventions
- **Factory functions**: bare type names → `String`, `Int`, `Float64`, `Err`, `Time`
- **Read methods**: `MarshalTo<Type>` prefix → `MarshalToStr`, `MarshalToInt`, `MarshalToFloat64`
- **Collection adders**: `Add<Type>` prefix → `AddStr`, `AddInt`, `AddAll`
- **Errors**: package-level vars → `ErrDstNil`, `ErrTypeMismatch`, `ErrEmptyUnifield`

### Error handling
- Nil destination pointer → `return ErrDstNil` (first check)
- No value stored / all pointers nil (ptr) / wrong tag (val) → returns descriptive `%w` wrapped error
- Type mismatch → `return fmt.Errorf("%w: got %s, want <target>", ErrTypeMismatch, typeName(u.type_))`
- Use `%w` wrapping, never `errors.New()` inline

### Immutability guarantees
- `Clone()` returns a copy — all fields copied by Go pass-by-value semantics (val); pointers themselves are copied (ptr)
- Collection `Add()` auto-clones input before appending
- External mutation of returned Unifield does NOT affect collection contents
- `AddAll(s []Unifielder)` pre-allocates cloned slice to avoid reallocation
- `Merge()` clones each source item individually, matching `Add()` invariant

### Float cross-type compatibility
Both `MarshalToFloat32` and `MarshalToFloat64` accept **either** `valueFloat32` OR `valueFloat64` since they share the same `f64` backing field. Narrowing/widening is implicit.

### Wrapper factory functions (parent `unifields` package)

The parent `unifields` package exposes 15 wrapper factory functions that delegate to `val`. Consumers do not need to import `val/` for normal usage:

```go
import "github.com/crypto-bundle/bc-wallet-common-lib-unifields/pkg/unifields"

// All work through the parent package — no val/ import needed
f := unifields.String("key", "value")   // delegates to val.String()
g := unifields.Int("counter", 42)       // delegates to val.Int()
h := unifields.Err("err", someErr)      // delegates to val.Err()
```

A type alias `Unifield = val.Unifield` is defined at the parent level so consumers can use the `Unifield` type name without importing `val/`.

#### Factory wrapper template (parent unifields package)

```go
func TypeName(key string, val <Type>) Unifield { // Unifield is aliased to val.Unifield
    return valpkg.TypeName(key, val)
}
```
Single-line delegation. Zero runtime overhead — compiler inlines these calls.

### Pointer-based variant (`unifields/ptr` sub-package)

`UnifieldPtr` stores pointers to typed values instead of flat fields. Only one pointer is non-nil at a time; all others remain nil. This trades per-call heap allocation (one alloc per factory call) for zero cost on unused slots in collections.

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
    if dst == nil {                      // blank line required after if-block (nlreturn rule)
        return ErrDstNil
    }
    if u.<field>P == nil {               // blank line required after if-block (nlreturn rule)
        return ErrEmptyUnifield
    }
    if !is<Predicate>(u.type_) {         // or exact match: if u.type_ != value<Type>
        return fmt.Errorf("%w: got %s, want <type>", ErrTypeMismatch, typeName(u.type_))
    }
    *dst = <conversion>(*u.<field>P)

    return nil                           // blank line before return (nlreturn rule)
}
```

#### Test patterns (ptr)

In addition to positive, type-mismatch, nil-dst, and zero-value tests, ptr tests must cover empty UnifieldPtr access errors (calling any MarshalTo on `{key: "k"}` should return `ErrEmptyUnifield`). Clone separation tests verify that cloned structs have independent `key` and `type_` fields. Note that `Clone()` returns `Unifielder` interface — type assertions in tests must handle the interface-to-concrete cast.

## Build & test commands

```bash
make lint                    # golangci-lint run with project rules
make test                    # go test -race ./...
go vet ./...                 # static analysis
go build ./...               # verify compilation
go test -bench=. -benchmem ./pkg/unifields/... # benchmarks
```

## Code style rules

### File headers
Every `.go` file must include the MIT NON-AI license block from `copyright.txt.tpl`. License text must match verbatim.

### Factory function template (val)
```go
func TypeName(key string, val <Type>) Unifield {
    return Unifield{key: key, <field>: val, type_: value<Type>}
}
```
Return plain struct literal. No nil checks — invalid inputs produce zero-value Unifield (documented behavior).

### Clone method template (both)
```go
func (u Unifield) Clone() unifolderv2.Unifielder { //nolint:ireturn // intentionally returns interface for polymorphic collections
    return u
}
```
`Clone()` returns `Unifielder` interface — this is intentional design for polymorphic collections. The `//nolint:ireturn` comment suppresses the linter warning since returning an interface is the core abstraction mechanism.

### MarshalTo method template (val)
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

### Test patterns
Group tests by target file:
- `val/unified_field_test.go`: factory + Clone + MarshalTo + error cases (same coverage as ptr)
- `ptr/unified_field_ptr_test.go`: factory + Clone + MarshalTo + empty access + error cases
- `val/unified_field_list_test.go`: all UnitfieldList methods — Add, AddAll, Merge (two-lists, self, empty-source, single-element, zero-value receiver), GetAfter (zero/middle/boundary/out-of-range/extra-out), RemoveAfter (keep-head, last-index, out-of-range, extra-out, single-element), RemoveBefore (keep-tail, index-zero, out-of-range, extra-out, single-element), Clear, Items (read-only), polymorphic mix
- `unified_fields_test.go`: NewUnifolds, Add, AddAll, polymorphic val/ptr, AddVal/AddPtr typed tests, ItemsReadOnly, UnitfolderListNew/Typed tests

Each supported type needs: positive test, type-mismatch negative, nil-dst negative, zero-value test.

For Clone assertions after `Clone()`, note that the method returns `Unifielder` interface — tests should use either value assertion `(Unifield)` or check results via `MarshalTo*` methods rather than direct field access across package boundaries. Add `//nolint:forcetypeassert` when asserting concrete types from `Clone()` return values.

### Linter exclusions (version 2 config)
See `.golangci.yml` for full rules. Key exclusions:
- `_test\.go` → cyclop, exhaustruct_v5, err113, funlen, gocognit, gocritic, goconst, nlreturn, paralleltest, govet, wsl_v5
- `unified_field\.go`, `unified_field_ptr\.go` → cyclop, exhaustruct_v5, goheader, gosec, nlreturn, ireturn, wsl_v5
- `unified_fields\.go` → cyclop, exhaustruct_v5, gosec, nlreturn, wsl_v5
- `unified_field_list\.go` → cyclop, exhaustruct_v5, gosec, nlreturn, wsl_v5
- `doc\.go` → goheader, gci
- Source lines starting with `* ` → lll (long lines in comments)

## Task tracking

Use `.agents/tasks/NNN_task_name.md` for task descriptions and `.agents/plans/NNN_plan_name.md` for execution plans. Each plan step maps to one git commit. Numbering must be sequential and match between tasks and plans. Update the plan's status table after completing each commit.
