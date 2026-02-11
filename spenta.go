package spenta

import "github.com/rouzbehsbz/spenta/internal/types"

type SliceParIter = types.SliceParIter

func SliceParForEach[T comparable](slice *[]T, cb func(e T)) *SliceParIter {
	return types.SliceParForEach(slice, cb)
}

func SliceParMap[T comparable](slice *[]T, cb func(e T) T) *SliceParIter {
	return types.SliceParMap(slice, cb)
}

func SliceParFilter[T comparable](slice *[]T, cb func(e T) bool) *SliceParIter {
	return types.SliceParFilter(slice, cb)
}
