package iter

import (
	"reflect"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
)

func TestSliceParForEach_VisitsAllElementsExactlyOnce(t *testing.T) {
	t.Parallel()

	arr := make([]int, 1024)
	for i := range arr {
		arr[i] = i
	}

	seen := make([]int32, len(arr))
	p := SliceParForEach(&arr, func(i int, v int) {
		if v != i {
			t.Fatalf("unexpected value at index %d: got=%d want=%d", i, v, i)
		}
		atomic.AddInt32(&seen[i], 1)
	}, WithMaxChunkSize(32), WithMinChunkSize(8))

	if err := p.Wait(); err != nil {
		t.Fatalf("Wait() unexpected error: %v", err)
	}

	for i := range seen {
		if seen[i] != 1 {
			t.Fatalf("index %d processed %d times, expected 1", i, seen[i])
		}
	}
}

func TestSliceParForEach_CapturesPanicsAsError(t *testing.T) {
	t.Parallel()

	arr := []int{1, 2, 3, 4}
	p := SliceParForEach(&arr, func(i int, v int) {
		if i == 2 {
			panic("slice foreach panic")
		}
	}, WithMaxChunkSize(2), WithMinChunkSize(1))

	err := p.Wait()
	if err == nil {
		t.Fatalf("expected panic to be captured as error, got nil")
	}
	if !strings.Contains(err.Error(), "slice foreach panic") {
		t.Fatalf("expected panic message in error, got %v", err)
	}
}

func TestSliceParMap_TransformsInPlace(t *testing.T) {
	t.Parallel()

	arr := []int{1, 2, 3, 4, 5}
	want := []int{2, 5, 8, 11, 14}

	p := SliceParMap(&arr, func(i int, v int) int {
		return v*2 + i
	}, WithMaxChunkSize(2), WithMinChunkSize(1))

	if err := p.Wait(); err != nil {
		t.Fatalf("Wait() unexpected error: %v", err)
	}

	if !reflect.DeepEqual(arr, want) {
		t.Fatalf("unexpected mapped result: got=%v want=%v", arr, want)
	}
}

func TestSliceParMap_CapturesPanicsAsError(t *testing.T) {
	t.Parallel()

	arr := []int{10, 20, 30}
	p := SliceParMap(&arr, func(i int, v int) int {
		if i == 1 {
			panic("slice map panic")
		}
		return v
	}, WithMaxChunkSize(2), WithMinChunkSize(1))

	err := p.Wait()
	if err == nil {
		t.Fatalf("expected panic to be captured as error, got nil")
	}
	if !strings.Contains(err.Error(), "slice map panic") {
		t.Fatalf("expected panic message in error, got %v", err)
	}
}

func TestSliceParFilter_KeepsMatchingElementsUnordered(t *testing.T) {
	t.Parallel()

	arr := make([]int, 200)
	for i := range arr {
		arr[i] = i
	}

	p := SliceParFilter(&arr, func(i int, v int) bool {
		return v%3 == 0
	}, WithMaxChunkSize(16), WithMinChunkSize(4))

	if err := p.Wait(); err != nil {
		t.Fatalf("Wait() unexpected error: %v", err)
	}

	want := make([]int, 0, 67)
	for i := 0; i < 200; i++ {
		if i%3 == 0 {
			want = append(want, i)
		}
	}

	got := slices.Clone(arr)
	slices.Sort(got)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected filtered values: got=%v want=%v", got, want)
	}
}

func TestSliceParFilter_NoMatchesResultsInEmptySlice(t *testing.T) {
	t.Parallel()

	arr := []int{1, 3, 5, 7}
	p := SliceParFilter(&arr, func(i int, v int) bool {
		return v%2 == 0
	}, WithMaxChunkSize(2), WithMinChunkSize(1))

	if err := p.Wait(); err != nil {
		t.Fatalf("Wait() unexpected error: %v", err)
	}
	if len(arr) != 0 {
		t.Fatalf("expected empty slice after filter, got len=%d values=%v", len(arr), arr)
	}
}

func TestSliceParFilter_CapturesPanicsAsError(t *testing.T) {
	t.Parallel()

	arr := []int{1, 2, 3, 4}
	p := SliceParFilter(&arr, func(i int, v int) bool {
		if i == 1 {
			panic("slice filter panic")
		}
		return true
	}, WithMaxChunkSize(2), WithMinChunkSize(1))

	err := p.Wait()
	if err == nil {
		t.Fatalf("expected panic to be captured as error, got nil")
	}
	if !strings.Contains(err.Error(), "slice filter panic") {
		t.Fatalf("expected panic message in error, got %v", err)
	}
}

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
