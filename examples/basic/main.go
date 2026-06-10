package main

import (
	"fmt"

	"github.com/rouzbehsbz/spenta"
)

func main() {
	arr := []int{1, 2, 3, 4, 5}

	parIter := spenta.SliceParMap(arr, func(i int, v int) int {
		return v * 2
	},
		spenta.WithMaxChunkSize(4),
		spenta.WithMinChunkSize(2),
	)

	_ = parIter.Wait()

	fmt.Println(arr)
}
