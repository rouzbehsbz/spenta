package iter

import (
	"sync"

	"github.com/rouzbehsbz/spenta/pool"
)

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
	parIter.startJobsDoneWatcher()

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

// Finds a matching key/value in the map in parallel.
// Because Go map iteration order is non-deterministic, the selected
// match is not guaranteed across runs when multiple keys match.
func MapParFind[K comparable, V any](Map *map[K]V, cb func(k K, v V) bool, opts ...ParIterOptions) *MapFindResult[K, V] {
	type localMatch struct {
		pos   int
		key   K
		value V
	}

	matches := make([]localMatch, 0)
	mu := &sync.Mutex{}

	p := NewMapParIter[K, V](Map, func(start, end int, keys []K) {
		for i := start; i < end; i++ {
			key := keys[i]
			value := (*Map)[key]

			if cb(key, value) {
				mu.Lock()
				matches = append(matches, localMatch{
					pos:   i,
					key:   key,
					value: value,
				})
				mu.Unlock()
				return
			}
		}
	}, opts...)

	result := &MapFindResult[K, V]{
		ParIter: p,
	}

	go func() {
		<-p.jobsDoneCh

		best := 0
		for i := range matches {
			if !result.found || matches[i].pos < matches[best].pos {
				best = i
				result.found = true
				result.key = matches[i].key
				result.value = matches[i].value
			}
		}

		p.postJobsDone()
	}()

	return result
}
