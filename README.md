# Spenta

[![Go Reference](https://pkg.go.dev/badge/github.com/rouzbehsbz/spenta.svg)](https://pkg.go.dev/github.com/rouzbehsbz/spenta)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/rouzbehsbz/spenta)
[![Build Status](https://github.com/rouzbehsbz/spenta/actions/workflows/build.yml/badge.svg)](https://github.com/rouzbehsbz/spenta/actions/workflows/build.yml)

> Spenta (سپنتا) — pronounced SPEN-ta — a Persian name of Avestan origin meaning “holy,” “pure,” and “life-giving.”

spenta is a lightweight and lock-free parallel iterator generator for Go. The library provides a simple and easy-to-use abstraction layer over data parallelism for iterable types such as `slices` and `map`.

Here’s how you can perform a simple slice mapping:

```go
arr := []int{1, 2, 3, 4, 5, 6}

parIter := iter.SliceParMap(&arr, func(i int, v int) int  {
	return v * 2
})

_ = parIter.Wait()
```

## Overview
spenta divides the original task into multiple subtasks, each performing computation over a portion of the original data. It returns the results exactly as a sequential execution would, but Spenta does it in parallel using multiple goroutines.

- Completely lock-free algorithms
- Automatically spawns workers based on CPU cores
- Supports common operations like `forEach`, `map`, and others
- Type-safe closures using Go’s generics
- Capture errors while allowing other tasks to continue

## Iterator Oprations
| Operation | Slice | Map |
|:--:|:--------:|:--------:|
|`forEach`|✅ |✅ |
|`map`|✅ |✅ |
|`filter`|✅ |✅ |
|`find`|✅ | ✅ |
|`reduce`|❌ | ❌ |

## Optimizations
- Adaptive binary divide-and-conquer chunking for balanced parallel execution.

## Optional Configurations
You can pass optional configurations to override the defaults as needed.

```go
parIter := iter.SliceParFilter(&arr, func(i int,  v int) bool {
	return v%2 == 0
},
	iter.WithMinChunkSize(5),  // Default min chunk size is 256
	iter.WithMaxChunkSize(20), // Default max chunk size is 4096
)
```
