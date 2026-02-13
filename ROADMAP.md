# Roadmap

## v1.0.0
- [x] Supported parallel iterator operations:
  - `forEach`
  - `map`
  - `filter`
- [x] Support parallel iterator for `map[K]V` type
- Support parallel iterator for `string` type
- Add tests and integrate with CI
- Add benchmarks for the library itself
- Add benchmarks to compare with other libraries
- Complete documents and pkg reference

# v1.1.0
- Add new operations that collect results from all jobs and perform calculations at the end:
  - `find`
  - `reduce`
  - `min`
  - `max`
  - `sum`
- Chaining operations
