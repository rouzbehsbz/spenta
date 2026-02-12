package main

import (
	"fmt"

	"github.com/rouzbehsbz/spenta"
)

func main() {
	arr := []int{1, 2, 3, 4, 5, 6}

	parIter := spenta.SliceParMap(&arr, func(a int) int {
		return a * 2
	})

	parIter.Done()

	fmt.Println(arr)
}
