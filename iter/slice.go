package iter

import (
	"sync"

	"github.com/rouzbehsbz/spenta/pool"
)

func NewSliceParIter[V any](slice *[]V, cb func(i int, v V), opts ...ParIterOptions) *ParIter {
	options := BuildParIterOptions(opts)

	sLen, chunkSize, chunkCount := SliceChunk(slice, options.MaxChunkSize)

	wg := &sync.WaitGroup{}
	wg.Add(chunkCount)

	jobs, errCh := pool.NewSliceJobs(sLen, chunkCount, chunkSize, wg, func(i int) {
		cb(i, (*slice)[i])
	})

	go pool.SpentaPool().SendJobs(jobs)

	return NewParIter(wg, errCh)
}

func SliceParForEach[V any](slice *[]V, cb func(i int, v V), opts ...ParIterOptions) *ParIter {
	return NewSliceParIter[V](slice, func(i int, v V) {
		cb(i, v)
	}, opts...)
}

func SliceParMap[V any](slice *[]V, cb func(i int, v V) V, opts ...ParIterOptions) *ParIter {
	return NewSliceParIter[V](slice, func(i int, v V) {
		(*slice)[i] = cb(i, v)
	}, opts...)
}

func SliceParFilter[V any](slice *[]V, cb func(i int, v V) bool, opts ...ParIterOptions) *ParIter {
	return NewSliceParIter[V](slice, func(i int, v V) {
		s := *slice
		j := 0
		for _, v := range s {
			if cb(i, v) {
				s[j] = v
				j++
			}
		}
		*slice = s[:j]
	}, opts...)
}
