package iter

// SliceFindResult stores the result of SliceParFind.
type SliceFindResult[V any] struct {
	*ParIter

	found bool
	index int
	value V
}

// Found reports whether a matching element was found.
func (r *SliceFindResult[V]) Found() bool {
	return r.found
}

// Index returns the index of the matching element.
// If no element matched, it returns -1.
func (r *SliceFindResult[V]) Index() int {
	return r.index
}

// Value returns the matching value.
// If no element matched, it returns the zero value of V.
func (r *SliceFindResult[V]) Value() V {
	return r.value
}

// WaitResult waits for all jobs and returns the final find result.
func (r *SliceFindResult[V]) WaitResult() (index int, value V, found bool, err error) {
	err = r.Wait()
	return r.index, r.value, r.found, err
}

// MapFindResult stores the result of MapParFind.
type MapFindResult[K comparable, V any] struct {
	*ParIter

	found bool
	key   K
	value V
}

// Found reports whether a matching key/value was found.
func (r *MapFindResult[K, V]) Found() bool {
	return r.found
}

// Key returns the matching key.
// If no element matched, it returns the zero value of K.
func (r *MapFindResult[K, V]) Key() K {
	return r.key
}

// Value returns the matching value.
// If no element matched, it returns the zero value of V.
func (r *MapFindResult[K, V]) Value() V {
	return r.value
}

// WaitResult waits for all jobs and returns the final find result.
func (r *MapFindResult[K, V]) WaitResult() (key K, value V, found bool, err error) {
	err = r.Wait()
	return r.key, r.value, r.found, err
}
