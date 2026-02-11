package main

import (
	"fmt"

	"github.com/rouzbehsbz/spenta"
)

func main() {
	arr := []int{1, 2, 3, 4, 5, 6}

	parIter := spenta.SliceParFilter(&arr, func(e int) bool {
		return e%2 == 0
	})

	parIter.Done()

	fmt.Println(arr)
}
