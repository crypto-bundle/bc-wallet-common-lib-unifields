# Change Log

## [v0.0.1] — 2026-10-04

### What's new

**unifields** is a library for storing different types of data in a single wrapper with named labels: strings, integers of any size (signed and unsigned), floating-point numbers, dates, and errors. Each object holds exactly one value and supports safe reading — if you try to read a string as a number, the library returns an error instead of crashing the program. Supports copying and collections without risk that external changes will affect stored data.
