# Change Log

## [v0.0.9] — 2026-10-09

### Changed 

**Breaking Changes**:

Removed the unused `key` parameter from all factory functions. This simplifies usage and reduces memory per object by ~8–16 bytes. If you create values like `unifields.String("id", value)`, just drop the first argument: `unifields.String(value)`.

## [v0.0.8] — 2026-10-08

### Added
- Automatic memory pre-allocation for `val.Unifield` via `sync.Pool` — factory functions now reuse pooled instances for reduced GC pressure
- New `Flush()` method on `UnifieldStack` — removes and returns all elements as a slice of clones; originals are returned to the pool
- Automatic pool recycling in list removal operations (`Clear`, `RemoveAfter`, `RemoveBefore`) — recycled `val.Unifield` instances are returned to the internal pool
- Automatic pool recycling in stack operations (`Pop`, `PopN`, `Clear`, `Flush`) — popped items are cloned before return; originals go to pool
- `returnValToPool()` helper — type-safe pool return that silently skips non-`val.Unifield` items (e.g., `*ptr.UnifieldPtr`)

### Changed
- Factory functions now use `sync.Pool` for reduced allocation churn (per-call memory footprint unchanged: 112 B/op, 1 alloc/op)
- `Pop()` and `PopN()` now return **clones** of popped elements instead of direct references — preserves immutability invariant
- `UnitfieldList.Clear()`, `RemoveAfter()`, `RemoveBefore()` now delegate to `*NoLock` variants with pool recycling
- `UnifieldStack.pop()`, `popNNolockInt()`, `clearNoLock()` updated to clone-before-pool pattern

## [v0.0.7] — 2026-10-08

### Added
- `UnitfieldList` is now thread-safe via `sync.RWMutex`, enabling concurrent access from multiple goroutines
- Three new concurrent tests validating race-free behavior under contention

### Changed
- Migrated all 15 typed adder methods from parent package to `val.UnitfieldList`
- Replaced parent-package struct wrapper with type alias (`type UnitfieldList = val.UnitfieldList`)
- `Merge()` now accepts a pointer to source list (`Merge(source *UnitfieldList)`) for proper mutex coordination
- Benchmarks recorded before and after sync.RWMutex integration

### Fixed
- Data race condition in `UnitfieldList` concurrent read/write operations — eliminated by adding mutex protection to all exported methods
- Removed `TestUnitfieldList_Merge_Self` which caused deadlock when calling `list.Merge(list)` due to Go's non-reentrant mutex

### Added
- `UnifieldStack` now supports safe concurrent access from multiple goroutines using `sync.RWMutex`. Write methods acquire exclusive locks; read-only methods use shared read locks for better throughput under contention.

### Changed
- Renamed `UnitfieldStack` → `UnifieldStack` throughout the codebase for consistency with the `val.Unifield` type name. All constructors, method receivers, tests, benchmarks, and documentation updated accordingly.
- Public stack methods delegate to private `*NoLock` variants that assume the caller holds the lock, keeping critical sections tight.

### Fixed
- Data race condition in `UnifieldStack` when accessed concurrently from multiple goroutines — eliminated by adding `sync.RWMutex` protection to all exported methods.


## [v0.0.5] — 2026-10-07

### Changed

- Stack operations now return the Unifielder interface instead of concrete types, allowing both value and pointer variants to be stored and retrieved from the same stack
- Empty-stack reads return nil consistently, matching how collection methods handle missing data

### Added

- Polymorphic val/ptr mixing supported in stack push operations

## [v0.0.3] — 2026-10-07

### Added

- New collection type with list manipulation methods (merge, split by index, clear) for organizing typed values

### Changed

- Recommended collection type updated; the older collection type remains available for backward compatibility but is no longer recommended for new code

### Fixed

- Immutable collection guarantee extended to all new list operations — every item addition clones before storing so external mutations never affect saved data


## [v0.0.2] — 2026-10-05

### Added

- Pointer-based value storage variant — trades per-item heap allocation for zero-cost empty slots in collections
- Expanded error handling — nil destination, type mismatch, and empty storage access all return clear errors
- Benchmarks for all operations measuring speed and memory usage
- Full godoc coverage on every exported method, AI-agent reference guide, updated README with both variants compared
- Zero-allocation value-storage variant extracted into its own sub-package for hot-path efficiency
- Convenience factory functions at the parent package level — no need to import sub-packages for everyday use
- Shared polymorphic interface enabling heterogeneous collections mixing both storage variants

### Changed

- Minimum Go version raised from 1.23 to 1.27

## [v0.0.1] — 2026-10-04

### What's new

**unifields** is a library for storing different types of data in a single wrapper with named labels: strings, integers of any size (signed and unsigned), floating-point numbers, dates, and errors. Each object holds exactly one value and supports safe reading — if you try to read a string as a number, the library returns an error instead of crashing the program. Supports copying and collections without risk that external changes will affect stored data.
