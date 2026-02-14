package iter

import "github.com/rouzbehsbz/spenta/pool"

// Creates a ParIter for processing a map in parallel.
func NewMapParIter[K comparable, V any](Map *map[K]V, cb func(start, end int, keys []K), opts ...ParIterOptions) *ParIter {
	options := BuildParIterOptions(opts)
	length := len(*Map)

	parIter := NewParIter()

	keys := make([]K, 0, length)
	for key := range *Map {
		keys = append(keys, key)
	}

	pool.SpawnJob(0, length, int(options.MaxChunkSize), int(options.MinChunkSize), parIter.jobsWg, parIter.errCh, func(start, end int) {
		cb(start, end, keys)
	})

	return parIter
}

// Applies the given callback function to each
// key, value of the map in parallel.
func MapParForEach[K comparable, V any](Map *map[K]V, cb func(k K, v V), opts ...ParIterOptions) *ParIter {
	p := NewMapParIter[K, V](Map, func(start, end int, keys []K) {
		for i := start; i < end; i++ {
			key := keys[i]

			cb(key, (*Map)[key])
		}
	}, opts...)

	p.postJobsDone()

	return p
}

// Applies the given transformation function to each key
// of the map in parallel and replaces each key with the returned value.
func MapParMap[K comparable, V any](Map *map[K]V, cb func(k K, v V) V, opts ...ParIterOptions) *ParIter {
	p := NewMapParIter[K, V](Map, func(start, end int, keys []K) {
		for i := start; i < end; i++ {
			key := keys[i]

			(*Map)[key] = cb(key, (*Map)[key])
		}
	}, opts...)

	p.postJobsDone()

	return p
}

// Filters the map in parallel according to the
// provided predicate function.
func MapParFilter[K comparable, V any](Map *map[K]V, cb func(k K, v V) bool, opts ...ParIterOptions) *ParIter {
	p := NewMapParIter[K, V](Map, func(start, end int, keys []K) {
		for i := start; i < end; i++ {
			key := keys[i]

			if !cb(key, (*Map)[key]) {
				delete(*Map, key)
			}
		}
	}, opts...)

	p.postJobsDone()

	return p
}
