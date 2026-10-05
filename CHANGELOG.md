# Change Log

## [v0.0.2] — 2026-10-05

### Added

- Pointer-based value storage variant — trades per-item heap allocation for zero-cost empty slots in collections
- Expanded error handling — nil destination, type mismatch, and empty storage access all return clear errors
- Benchmarks for all operations measuring speed and memory usage
- Full godoc coverage on every exported method, AI-agent reference guide, updated README with both variants compared
- Zero-allocation value-storage variant extracted into its own sub-package for hot-path efficiency
- Convenience factory functions at the parent package level — no need to import sub-packages for everyday use
- Shared polymorphic interface enabling heterogeneous collections mixing both storage variants

## [v0.0.1] — 2026-10-04

### What's new

**unifields** is a library for storing different types of data in a single wrapper with named labels: strings, integers of any size (signed and unsigned), floating-point numbers, dates, and errors. Each object holds exactly one value and supports safe reading — if you try to read a string as a number, the library returns an error instead of crashing the program. Supports copying and collections without risk that external changes will affect stored data.
