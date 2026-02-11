package spenta

import "github.com/rouzbehsbz/spenta/internal"

type ParSliceIter = internal.ParSliceIter

func SliceParForEach[T comparable](slice *[]T, cb func(e T)) *ParSliceIter {
	return internal.SliceParForEach(slice, cb)
}

func SliceParMap[T comparable](slice *[]T, cb func(e T) T) *ParSliceIter {
	return internal.SliceParMap(slice, cb)
}
