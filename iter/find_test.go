package iter

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestSliceParFind_FindsFirstMatchingIndexAcrossChunks(t *testing.T) {
	t.Parallel()

	arr := make([]int, 1000)
	arr[120] = 7
	arr[845] = 7
	arr[999] = 7

	res := SliceParFind(&arr, func(i int, v int) bool {
		return v == 7
	}, WithMaxChunkSize(64), WithMinChunkSize(16))

	index, value, found, err := res.WaitResult()
	if err != nil {
		t.Fatalf("WaitResult() unexpected error: %v", err)
	}
	if !found {
		t.Fatalf("expected found=true, got false")
	}
	if index != 120 {
		t.Fatalf("expected first matching index=120, got %d", index)
	}
	if value != 7 {
		t.Fatalf("expected value=7, got %d", value)
	}
}

func TestSliceParFind_NoMatch(t *testing.T) {
	t.Parallel()

	arr := []int{1, 3, 5, 7, 9}

	res := SliceParFind(&arr, func(i int, v int) bool {
		return v%2 == 0
	}, WithMaxChunkSize(2), WithMinChunkSize(1))

	index, value, found, err := res.WaitResult()
	if err != nil {
		t.Fatalf("WaitResult() unexpected error: %v", err)
	}
	if found {
		t.Fatalf("expected found=false, got true")
	}
	if index != -1 {
		t.Fatalf("expected index=-1 when no match exists, got %d", index)
	}
	if value != 0 {
		t.Fatalf("expected zero value when no match exists, got %d", value)
	}
}

func TestSliceParFind_EmptySlice(t *testing.T) {
	t.Parallel()

	arr := []int{}

	res := SliceParFind(&arr, func(i int, v int) bool {
		return true
	}, WithMaxChunkSize(2), WithMinChunkSize(1))

	index, value, found, err := res.WaitResult()
	if err != nil {
		t.Fatalf("WaitResult() unexpected error: %v", err)
	}
	if found {
		t.Fatalf("expected found=false for empty slice")
	}
	if index != -1 {
		t.Fatalf("expected index=-1 for empty slice, got %d", index)
	}
	if value != 0 {
		t.Fatalf("expected zero value for empty slice, got %d", value)
	}
}

func TestSliceParFind_DoesNotMutateInput(t *testing.T) {
	t.Parallel()

	arr := []int{2, 4, 6, 8, 10, 12}
	original := slices.Clone(arr)

	res := SliceParFind(&arr, func(i int, v int) bool {
		return v == 8
	}, WithMaxChunkSize(3), WithMinChunkSize(1))

	_, _, _, err := res.WaitResult()
	if err != nil {
		t.Fatalf("WaitResult() unexpected error: %v", err)
	}

	if !reflect.DeepEqual(arr, original) {
		t.Fatalf("slice was mutated: got=%v want=%v", arr, original)
	}
}

func TestSliceParFind_WaitResultMatchesAccessors(t *testing.T) {
	t.Parallel()

	arr := []string{"a", "bb", "ccc", "dddd"}
	res := SliceParFind(&arr, func(i int, v string) bool {
		return len(v) > 2
	}, WithMaxChunkSize(2), WithMinChunkSize(1))

	index, value, found, err := res.WaitResult()
	if err != nil {
		t.Fatalf("WaitResult() unexpected error: %v", err)
	}

	if found != res.Found() || index != res.Index() || value != res.Value() {
		t.Fatalf(
			"inconsistent result state: wait=(%v,%d,%q) accessors=(%v,%d,%q)",
			found, index, value,
			res.Found(), res.Index(), res.Value(),
		)
	}
}

func TestSliceParFind_CapturesPanicsAsError(t *testing.T) {
	t.Parallel()

	arr := []int{1, 2, 3, 4}

	res := SliceParFind(&arr, func(i int, v int) bool {
		if i == 2 {
			panic("slice find panic")
		}
		return false
	}, WithMaxChunkSize(32), WithMinChunkSize(1))

	_, _, _, err := res.WaitResult()
	if err == nil {
		t.Fatalf("expected panic to be captured as error, got nil")
	}
	if !strings.Contains(err.Error(), "slice find panic") {
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
			panic(fmt.Sprintf("map find panic key=%s", k))
		}
		return false
	}, WithMaxChunkSize(32), WithMinChunkSize(1))

	_, _, _, err := res.WaitResult()
	if err == nil {
		t.Fatalf("expected panic to be captured as error, got nil")
	}
	if !strings.Contains(err.Error(), "map find panic key=") {
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
