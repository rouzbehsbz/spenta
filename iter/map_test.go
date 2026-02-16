package iter

import (
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestMapParForEach_VisitsAllPairsExactlyOnce(t *testing.T) {
	t.Parallel()

	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
		"f": 6,
	}

	seen := map[string]int{}
	var mu sync.Mutex

	p := MapParForEach(&m, func(k string, v int) {
		mu.Lock()
		seen[k]++
		mu.Unlock()

		if source, ok := m[k]; !ok || source != v {
			t.Fatalf("mismatch for key %q: callback value=%d map value=%d", k, v, source)
		}
	}, WithMaxChunkSize(2), WithMinChunkSize(1))

	if err := p.Wait(); err != nil {
		t.Fatalf("Wait() unexpected error: %v", err)
	}

	if len(seen) != len(m) {
		t.Fatalf("expected %d visited keys, got %d", len(m), len(seen))
	}
	for k := range m {
		if seen[k] != 1 {
			t.Fatalf("key %q processed %d times, expected 1", k, seen[k])
		}
	}
}

func TestMapParForEach_CapturesPanicsAsError(t *testing.T) {
	t.Parallel()

	m := map[string]int{"a": 1, "b": 2}
	p := MapParForEach(&m, func(k string, v int) {
		if k == "b" {
			panic("map foreach panic")
		}
	}, WithMaxChunkSize(2), WithMinChunkSize(1))

	err := p.Wait()
	if err == nil {
		t.Fatalf("expected panic to be captured as error, got nil")
	}
	if !strings.Contains(err.Error(), "map foreach panic") {
		t.Fatalf("expected panic message in error, got %v", err)
	}
}

func TestMapParMap_TransformsValues(t *testing.T) {
	t.Parallel()

	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}

	p := MapParMap(&m, func(k string, v int) int {
		return v * 10
	}, WithMaxChunkSize(uint(len(m)+1)), WithMinChunkSize(1))

	if err := p.Wait(); err != nil {
		t.Fatalf("Wait() unexpected error: %v", err)
	}

	want := map[string]int{
		"a": 10,
		"b": 20,
		"c": 30,
		"d": 40,
	}
	if !reflect.DeepEqual(m, want) {
		t.Fatalf("unexpected mapped map: got=%v want=%v", m, want)
	}
}

func TestMapParMap_CapturesPanicsAsError(t *testing.T) {
	t.Parallel()

	m := map[string]int{"a": 1, "b": 2, "c": 3}
	p := MapParMap(&m, func(k string, v int) int {
		if k == "b" {
			panic("map map panic")
		}
		return v
	}, WithMaxChunkSize(uint(len(m)+1)), WithMinChunkSize(1))

	err := p.Wait()
	if err == nil {
		t.Fatalf("expected panic to be captured as error, got nil")
	}
	if !strings.Contains(err.Error(), "map map panic") {
		t.Fatalf("expected panic message in error, got %v", err)
	}
}

func TestMapParFilter_KeepsOnlyMatchingPairs(t *testing.T) {
	t.Parallel()

	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
	}

	p := MapParFilter(&m, func(k string, v int) bool {
		return v%2 == 0
	}, WithMaxChunkSize(uint(len(m)+1)), WithMinChunkSize(1))

	if err := p.Wait(); err != nil {
		t.Fatalf("Wait() unexpected error: %v", err)
	}

	want := map[string]int{"b": 2, "d": 4}
	if !reflect.DeepEqual(m, want) {
		t.Fatalf("unexpected filtered map: got=%v want=%v", m, want)
	}
}

func TestMapParFilter_NoMatchesResultsInEmptyMap(t *testing.T) {
	t.Parallel()

	m := map[string]int{"a": 1, "b": 3, "c": 5}

	p := MapParFilter(&m, func(k string, v int) bool {
		return false
	}, WithMaxChunkSize(uint(len(m)+1)), WithMinChunkSize(1))

	if err := p.Wait(); err != nil {
		t.Fatalf("Wait() unexpected error: %v", err)
	}
	if len(m) != 0 {
		t.Fatalf("expected empty map after filter, got %v", m)
	}
}

func TestMapParFilter_CapturesPanicsAsError(t *testing.T) {
	t.Parallel()

	m := map[string]int{"a": 1, "b": 2, "c": 3}
	p := MapParFilter(&m, func(k string, v int) bool {
		if k == "b" {
			panic("map filter panic")
		}
		return true
	}, WithMaxChunkSize(uint(len(m)+1)), WithMinChunkSize(1))

	err := p.Wait()
	if err == nil {
		t.Fatalf("expected panic to be captured as error, got nil")
	}
	if !strings.Contains(err.Error(), "map filter panic") {
		t.Fatalf("expected panic message in error, got %v", err)
	}
}

func TestMapParFind_FindsMatchingPair(t *testing.T) {
	t.Parallel()

	m := map[string]int{
		"a": 1,
		"b": 4,
		"c": 9,
		"d": 16,
	}

	res := MapParFind(&m, func(k string, v int) bool {
		return v > 8
	}, WithMaxChunkSize(2), WithMinChunkSize(1))

	key, value, found, err := res.WaitResult()
	if err != nil {
		t.Fatalf("WaitResult() unexpected error: %v", err)
	}
	if !found {
		t.Fatalf("expected found=true, got false")
	}
	if valueFromMap, ok := m[key]; !ok || valueFromMap != value {
		t.Fatalf("result key/value mismatch against map: key=%q value=%d", key, value)
	}
	if value <= 8 {
		t.Fatalf("returned pair does not satisfy predicate: key=%q value=%d", key, value)
	}
}

func TestMapParFind_NoMatch(t *testing.T) {
	t.Parallel()

	m := map[string]int{
		"a": 1,
		"b": 2,
	}

	res := MapParFind(&m, func(k string, v int) bool {
		return v > 100
	}, WithMaxChunkSize(2), WithMinChunkSize(1))

	key, value, found, err := res.WaitResult()
	if err != nil {
		t.Fatalf("WaitResult() unexpected error: %v", err)
	}
	if found {
		t.Fatalf("expected found=false, got true")
	}
	if key != "" {
		t.Fatalf("expected zero key when no match exists, got %q", key)
	}
	if value != 0 {
		t.Fatalf("expected zero value when no match exists, got %d", value)
	}
}

func TestMapParFind_EmptyMap(t *testing.T) {
	t.Parallel()

	m := map[string]int{}

	res := MapParFind(&m, func(k string, v int) bool {
		return true
	}, WithMaxChunkSize(2), WithMinChunkSize(1))

	key, value, found, err := res.WaitResult()
	if err != nil {
		t.Fatalf("WaitResult() unexpected error: %v", err)
	}
	if found {
		t.Fatalf("expected found=false for empty map")
	}
	if key != "" {
		t.Fatalf("expected zero key for empty map, got %q", key)
	}
	if value != 0 {
		t.Fatalf("expected zero value for empty map, got %d", value)
	}
}

func TestMapParFind_DoesNotMutateInput(t *testing.T) {
	t.Parallel()

	m := map[string]int{
		"x": 10,
		"y": 20,
		"z": 30,
	}
	original := cloneMap(m)

	res := MapParFind(&m, func(k string, v int) bool {
		return v == 20
	}, WithMaxChunkSize(2), WithMinChunkSize(1))

	_, _, _, err := res.WaitResult()
	if err != nil {
		t.Fatalf("WaitResult() unexpected error: %v", err)
	}

	if !reflect.DeepEqual(m, original) {
		t.Fatalf("map was mutated: got=%v want=%v", m, original)
	}
}

func TestMapParFind_WaitResultMatchesAccessors(t *testing.T) {
	t.Parallel()

	m := map[string]int{
		"x": 10,
		"y": 20,
	}
	res := MapParFind(&m, func(k string, v int) bool {
		return v%2 == 0
	}, WithMaxChunkSize(2), WithMinChunkSize(1))

	key, value, found, err := res.WaitResult()
	if err != nil {
		t.Fatalf("WaitResult() unexpected error: %v", err)
	}

	if found != res.Found() || key != res.Key() || value != res.Value() {
		t.Fatalf(
			"inconsistent result state: wait=(%v,%q,%d) accessors=(%v,%q,%d)",
			found, key, value,
			res.Found(), res.Key(), res.Value(),
		)
	}
}

func TestMapParFind_CapturesPanicsAsError(t *testing.T) {
	t.Parallel()

	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	res := MapParFind(&m, func(k string, v int) bool {
		if k == "b" {
			panic("map find panic")
		}
		return false
	}, WithMaxChunkSize(32), WithMinChunkSize(1))

	_, _, _, err := res.WaitResult()
	if err == nil {
		t.Fatalf("expected panic to be captured as error, got nil")
	}
	if !strings.Contains(err.Error(), "map find panic") {
		t.Fatalf("expected panic message in error, got %v", err)
	}
}

func cloneMap[K comparable, V any](m map[K]V) map[K]V {
	cloned := make(map[K]V, len(m))
	for k, v := range m {
		cloned[k] = v
	}
	return cloned
}
