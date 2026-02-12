# Spenta

[![Go Reference](https://pkg.go.dev/badge/github.com/rouzbehsbz/spenta.svg)](https://pkg.go.dev/github.com/rouzbehsbz/spenta)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/rouzbehsbz/spenta)

> Spenta (سپنتا) — pronounced SPEN-ta — a Persian name of Avestan origin meaning “holy,” “pure,” and “life-giving.”

spenta is a lightweight and lock-free parallel iterator generator for Go. The library provides a simple and easy-to-use abstraction layer over data parallelism for iterable types such as `slices`.

Here’s how you can perform a simple slice mapping:

```go
	arr := []int{1, 2, 3, 4, 5, 6}

	parIter := spenta.SliceParMap(&arr, func(a int) int  {
		return a * 2
	})

	parIter.Done()
```

## Overview
spenta divides the original task into multiple subtasks, each performing computation over a portion of the original data. It returns the results exactly as a sequential execution would, but Spenta does it in parallel using multiple goroutines.

- Completely lock-free algorithms allow maximum parallel computation without using locking mechanisms
- An internal lightweight thread pool automatically spawns worker goroutines based on the number of available CPU cores.
- Supports multiple iterator operations, including `forEach`, `map` and more.
- Type-safe closures using Go’s recently introduced generics.

## Iterator Functions
| Operation | Slice | Map |
|:--:|:--------:|:--------:|
|`forEach`|✅ | ❌ |
|`map`|✅ | ❌ |
|`filter`|✅ | ❌ |
|`reduce`|❌ | ❌ |
|`find`|❌ | ❌ |
