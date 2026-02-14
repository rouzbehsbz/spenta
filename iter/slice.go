package iter

import (
	"sync"

	"github.com/rouzbehsbz/spenta/pool"
)

func NewSliceParIter[V any](slice *[]V, cb func(i int, v V), opts ...ParIterOptions) *ParIter {
	options := BuildParIterOptions(opts)
	length := len(*slice)

	parIter := NewParIter()

	pool.SpawnJob(0, length, int(options.MaxChunkSize), int(options.MinChunkSize), parIter.wg, parIter.errCh, func(i int) {
		cb(i, (*slice)[i])
	})

	return parIter
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
	options := BuildParIterOptions(opts)
	input := append([]V(nil), (*slice)...)
	length := len(input)
	keep := make([]bool, length)

	parIter := NewParIter()
	predicateWG := &sync.WaitGroup{}

	pool.SpawnJob(0, length, int(options.MaxChunkSize), int(options.MinChunkSize), predicateWG, parIter.errCh, func(i int) {
		keep[i] = cb(i, input[i])
	})

	parIter.wg.Add(1)
	go func() {
		defer parIter.wg.Done()
		predicateWG.Wait()

		s := *slice
		j := 0
		for i, v := range input {
			if keep[i] {
				s[j] = v
				j++
			}
		}
		*slice = s[:j]
	}()

	return parIter
}
