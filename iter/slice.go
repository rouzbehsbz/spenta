package iter

import (
	"sync"

	"github.com/rouzbehsbz/spenta/pool"
)

// Creates a ParIter for processing a slice in parallel.
func newSliceParIter[V any](slice *[]V, cb func(start, end int), opts ...ParIterOptions) *ParIter {
	options := BuildParIterOptions(opts)
	length := len(*slice)

	parIter := NewParIter()

	pool.SpawnJob(
		0,
		length,
		int(options.MaxChunkSize),
		int(options.MinChunkSize),
		parIter.jobsWg,
		parIter.errCh,
		func(start, end int) {
			cb(start, end)
		},
	)
	parIter.startJobsDoneWatcher()

	return parIter
}

// Applies the given callback function to each element
// of the slice in parallel.
func SliceParForEach[V any](slice *[]V, cb func(i int, v V), opts ...ParIterOptions) *ParIter {
	p := newSliceParIter[V](slice, func(start, end int) {
		for i := start; i < end; i++ {
			cb(i, (*slice)[i])
		}
	}, opts...)

	p.postJobsDone()

	return p
}

// Applies the given transformation function to each element
// of the slice in parallel and replaces each element with the returned value.
func SliceParMap[V any](slice *[]V, cb func(i int, v V) V, opts ...ParIterOptions) *ParIter {
	p := newSliceParIter[V](slice, func(start, end int) {
		for i := start; i < end; i++ {
			(*slice)[i] = cb(i, (*slice)[i])
		}
	}, opts...)

	p.postJobsDone()

	return p
}

// Filters the slice in parallel according to the
// provided predicate function. (it is unordered)
func SliceParFilter[V any](slice *[]V, cb func(i int, v V) bool, opts ...ParIterOptions) *ParIter {
	merge := []V{}
	mu := &sync.Mutex{}

	p := newSliceParIter(slice, func(start, end int) {
		local := make([]V, 0, end-start)

		copy(local, (*slice)[start:end])

		for i := start; i < end; i++ {
			if cb(i, (*slice)[i]) {
				local = append(local, (*slice)[i])
			}
		}

		mu.Lock()
		merge = append(merge, local...)
		mu.Unlock()
	}, opts...)

	go func() {
		<-p.jobsDoneCh
		*slice = merge
		p.postJobsDone()
	}()

	return p
}

// Finds the first matching element in the slice in parallel.
// It returns the same result as sequential execution would,
// based on the smallest matching index.
func SliceParFind[V any](slice *[]V, cb func(i int, v V) bool, opts ...ParIterOptions) *SliceFindResult[V] {
	type localMatch struct {
		index int
		value V
	}

	matches := make([]localMatch, 0)
	mu := &sync.Mutex{}

	p := newSliceParIter[V](slice, func(start, end int) {
		for i := start; i < end; i++ {
			if cb(i, (*slice)[i]) {
				mu.Lock()
				matches = append(matches, localMatch{
					index: i,
					value: (*slice)[i],
				})
				mu.Unlock()
				return
			}
		}
	}, opts...)

	result := &SliceFindResult[V]{
		ParIter: p,
		index:   -1,
	}

	go func() {
		<-p.jobsDoneCh

		for _, match := range matches {
			if !result.found || match.index < result.index {
				result.found = true
				result.index = match.index
				result.value = match.value
			}
		}

		p.postJobsDone()
	}()

	return result
}
