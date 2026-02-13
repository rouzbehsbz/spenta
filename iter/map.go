package iter

import "github.com/rouzbehsbz/spenta/pool"

func NewMapParIter[K comparable, V any](Map *map[K]V, cb func(k K, v V), opts ...ParIterOptions) *ParIter {
	options := BuildParIterOptions(opts)
	length := len(*Map)

	parIter := NewParIter()

	keys := make([]K, 0, length)
	for key := range *Map {
		keys = append(keys, key)
	}

	pool.SpawnJob(0, length, int(options.MaxChunkSize), int(options.MinChunkSize), parIter.wg, parIter.errCh, func(i int) {
		key := keys[i]
		cb(key, (*Map)[key])
	})

	return parIter
}

func MapParForEach[K comparable, V any](Map *map[K]V, cb func(k K, v V), opts ...ParIterOptions) *ParIter {
	return NewMapParIter[K, V](Map, func(k K, v V) {
		cb(k, v)
	}, opts...)
}

func MapParMap[K comparable, V any](Map *map[K]V, cb func(k K, v V) V, opts ...ParIterOptions) *ParIter {
	return NewMapParIter[K, V](Map, func(k K, v V) {
		(*Map)[k] = cb(k, v)
	}, opts...)
}

func MapParFilter[K comparable, V any](Map *map[K]V, cb func(k K, v V) bool, opts ...ParIterOptions) *ParIter {
	return NewMapParIter[K, V](Map, func(k K, v V) {
		if !cb(k, v) {
			delete(*Map, k)
		}
	}, opts...)
}
