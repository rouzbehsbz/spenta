package main

import (
	"fmt"

	"github.com/rouzbehsbz/spenta/iter"
)

func main() {
	arr := []int{1, 2, 3, 4, 5}

	parIter := iter.SliceParMap(&arr, func(i int, v int) int {
		return v * 2
	}, iter.WithMaxChunkSize(3))

	_ = parIter.Wait()

	fmt.Println(arr)
}
